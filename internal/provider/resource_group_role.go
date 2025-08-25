package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/iam-client/pkg/types"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func resourceGroupRole() *schema.Resource {
	return &schema.Resource{
		Description:   "Binds one or more IAM roles to a group within a Sotoon workspace.",
		CreateContext: resourceGroupRoleCreate,
		ReadContext:   resourceGroupRoleRead,
		DeleteContext: resourceGroupRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Composite stable identifier. Does not affect lifecycle.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Group UUID.",
			},
			"role_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Set of Role UUIDs to bind to the group.",
			},
			"items": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Optional key/value items to pass to the bind API for each role.",
			},
			"bindings_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SHA-256 of sorted, canonical role_ids. Changes when the set of roles changes.",
			},
		},
	}
}

func resourceGroupRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	groupStr := d.Get("group_id").(string)
	groupUUID, err := uuid.FromString(groupStr)
	if err != nil {
		return diag.Errorf("invalid group_id %q: %s", groupStr, err)
	}

	roleIDsInput := fromSchemaSetToStrings(d.Get("role_ids").(*schema.Set))
	roleIDsCanon, err := convertStringsToUUIDArray(roleIDsInput)
	if err != nil {
		return diag.FromErr(fmt.Errorf("invalid role_ids: %w", err))
	}
	roleIDs := uniqueSorted(roleIDsCanon)
	bindingsHash := hashOfIDs(roleIDs)

	items := map[string]string{}
	if raw, ok := d.GetOk("items"); ok && raw != nil {
		for k, v := range raw.(map[string]interface{}) {
			items[k] = fmt.Sprintf("%v", v)
		}
	}

	rolesWithItems := make([]types.RoleWithItems, 0, len(roleIDs))
	for _, r := range roleIDs {
		rid, _ := uuid.FromString(r)
		rolesWithItems = append(rolesWithItems, types.RoleWithItems{
			RoleUUID: rid.String(),
			Items:    []map[string]string{items},
		})
	}

	if err := c.BulkAddRolesToGroup(c.WorkspaceUUID, &groupUUID, rolesWithItems); err != nil {
		return diag.Errorf("bulk bind roles to group %s failed: %s", groupUUID.String(), err)
	}

	d.SetId(fmt.Sprintf("%s:%s", groupUUID.String(), bindingsHash[:16]))
	_ = d.Set("bindings_hash", bindingsHash)
	return resourceGroupRoleRead(ctx, d, meta)
}

func resourceGroupRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	roleIDsInput := fromSchemaSetToStrings(d.Get("role_ids").(*schema.Set))
	roleIDsCanon, err := convertStringsToUUIDArray(roleIDsInput)
	if err == nil {
		roleIDs := uniqueSorted(roleIDsCanon)
		_ = d.Set("bindings_hash", hashOfIDs(roleIDs))
	}
	tflog.Info(ctx, "Reading group role", map[string]interface{}{"id": d.Id()})
	return nil
}

func resourceGroupRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	groupStr := d.Get("group_id").(string)
	groupUUID, err := uuid.FromString(groupStr)
	if err != nil {
		return diag.Errorf("invalid group_id %q: %s", groupStr, err)
	}

	roleIds := d.Get("role_ids").(*schema.Set).List()
	for _, v := range roleIds {
		s := v.(string)
		u, err := uuid.FromString(s)
		if err != nil {
			return diag.Errorf("invalid role_id in list: %s", err)
		}

		if err := c.UnbindRoleFromGroup(c.WorkspaceUUID, &u, &groupUUID, map[string]string{}); err != nil {
			return diag.Errorf("unbind role %s from group %s failed: %s", u.String(), groupUUID.String(), err)
		}
	}

	d.SetId("")
	return nil
}
