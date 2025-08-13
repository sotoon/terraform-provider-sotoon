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
				Description: `Composite ID in the form "<group_uuid>;<role_uuid>[;<role_uuid>...]" (new) or "<group_uuid>/<role_uuid>" (legacy).`,
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

	rawRoles := d.Get("role_ids").(*schema.Set)
	if rawRoles == nil || rawRoles.Len() == 0 {
		return diag.Errorf("role_ids must be provided and contain at least one UUID")
	}

	roleUUIDs := make([]uuid.UUID, 0, rawRoles.Len())
	for _, v := range rawRoles.List() {
		rs := v.(string)
		rid, err := uuid.FromString(rs)
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid role_ids entry %q: %w", rs, err))
		}
		roleUUIDs = append(roleUUIDs, rid)
	}

	items := map[string]string{}
	if raw, ok := d.GetOk("items"); ok && raw != nil {
		for k, v := range raw.(map[string]interface{}) {
			items[k] = fmt.Sprintf("%v", v)
		}
	}

	rolesWithItems := make([]types.RoleWithItems, 0, len(roleUUIDs))
	for _, r := range roleUUIDs {
		rolesWithItems = append(rolesWithItems, types.RoleWithItems{
			RoleUUID: r.String(),
			Items:    []map[string]string{items},
		})
	}

	if err := c.BulkAddRolesToGroup(c.WorkspaceUUID, &groupUUID, rolesWithItems); err != nil {
		return diag.Errorf("bulk bind roles to group %s failed: %s", groupUUID.String(), err)
	}

	d.SetId(groupUUID.String() + ":" + uuid.NewV4().String())
	return resourceGroupRoleRead(ctx, d, meta)
}

func resourceGroupRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

		if err := c.UnbindRoleFromGroup(c.WorkspaceUUID, &groupUUID, &u, map[string]string{}); err != nil {
			return diag.Errorf("unbind group %s from role %s failed: %s", u.String(), groupUUID.String(), err)
		}
	}

	d.SetId("")
	return nil
}
