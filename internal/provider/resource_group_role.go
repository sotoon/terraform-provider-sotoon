package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	iam "github.com/sotoon/sotoon-sdk-go/sdk/core/iam_v1"
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
				MinItems:    1,
				ForceNew:    true,
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

	groupID := d.Get("group_id").(string)
	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		return diag.Errorf("invalid group_id %q: %s", groupID, err)
	}

	sortedRoleIds := uniqueSorted(fromSchemaSetToStrings(d.Get("role_ids").(*schema.Set)))

	rolesList, err := c.GetWorkspaceGroupRoleList(ctx, c.WorkspaceUUID, &groupUUID)
	if err != nil {
		return diag.Errorf("read group roles: %s ", err)
	}
	remoteRolesID := make([]string, 0, len(rolesList))
	for _, r := range rolesList {
		remoteRolesID = append(remoteRolesID, r.Uuid)
	}
	remoteRolesID = uniqueSorted(remoteRolesID)

	items := map[string]string{}
	if raw, ok := d.GetOk("items"); ok && raw != nil {
		for k, v := range raw.(map[string]interface{}) {
			items[k] = fmt.Sprintf("%v", v)
		}
	}

	toAddList := diff(toSet(sortedRoleIds), toSet(remoteRolesID))
	if len(toAddList) > 0 {
		rolesWithItems := make([]iam.IamRoleItem, 0, len(toAddList))
		for _, id := range toAddList {
			rolesWithItems = append(rolesWithItems, iam.IamRoleItem{
				RoleUuid:  id,
				ItemsList: &[]map[string]string{items}})
		}
		if err := c.BulkAddRolesToGroup(ctx, &groupUUID, rolesWithItems); err != nil {
			return diag.Errorf("bulk bind roles to group %s failed: %s", groupUUID.String(), err)
		}
	}

	bindHash := hashOfIDs(sortedRoleIds)
	if err := d.Set("bindings_hash", bindHash); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set bindings_hash: %w", err))
	}
	d.SetId(groupUUID.String() + ":" + bindHash)

	return resourceGroupRoleRead(ctx, d, meta)
}

func resourceGroupRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	groupID := d.Get("group_id").(string)
	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		d.SetId("")
		return nil
	}

	sortedRoleIds := uniqueSorted(fromSchemaSetToStrings(d.Get("role_ids").(*schema.Set)))

	rolesList, err := c.GetWorkspaceGroupRoleList(ctx, c.WorkspaceUUID, &groupUUID)
	if err != nil {
		return diag.Errorf("error reading group role: %s", err)
	}

	remoteRoles := make([]string, 0, len(rolesList))
	for _, r := range rolesList {
		if r.Uuid != "" {
			remoteRoles = append(remoteRoles, r.Uuid)
			continue
		}
	}
	remoteRoles = uniqueSorted(remoteRoles)

	eff := intersect(toSet(sortedRoleIds), toSet(remoteRoles))
	effective := uniqueSorted(setKeys(eff))

	if err := d.Set("role_ids", effective); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set role_ids: %w", err))
	}

	bindHash := d.Get("bindings_hash").(string)
	if bindHash == "" {
		bindHash = hashOfIDs(effective)
		if err := d.Set("bindings_hash", bindHash); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set bindings_hash: %w", err))
		}
	}

	d.SetId(groupUUID.String() + ":" + bindHash)
	tflog.Info(ctx, "Reading group role", map[string]interface{}{"id": d.Id(), "group_id": groupID, "have": len(effective)})
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

		if err := c.UnbindRoleFromGroup(ctx, &u, &groupUUID); err != nil {
			return diag.Errorf("unbind role %s from group %s failed: %s", u.String(), groupUUID.String(), err)
		}
	}

	d.SetId("")
	return nil
}
