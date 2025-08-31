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

	rolesList, err := c.GetWorkspaceGroupRoleList(ctx, &groupUUID, c.WorkspaceUUID)
	if err != nil {
		return diag.Errorf("read group roles: %s", err)
	}
	remoteRolesID := make([]string, 0, len(rolesList))
	for _, r := range rolesList {
		if r.UUID != nil {
			remoteRolesID = append(remoteRolesID, r.UUID.String())
			continue
		}
		if rr, ok := any(r).(interface{ GetRoleUUID() string }); ok {
			if id := rr.GetRoleUUID(); id != "" {
				remoteRolesID = append(remoteRolesID, id)
			}
		}
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
		rolesWithItems := make([]types.RoleWithItems, 0, len(toAddList))
		for _, id := range toAddList {
			roleUUID, _ := uuid.FromString(id)
			rolesWithItems = append(rolesWithItems, types.RoleWithItems{RoleUUID: roleUUID.String(), Items: []map[string]string{items}})
		}
		if err := c.BulkAddRolesToGroup(c.WorkspaceUUID, &groupUUID, rolesWithItems); err != nil {
			return diag.Errorf("bulk bind roles to group %s failed: %s", groupUUID.String(), err)
		}
	}

	bindHash := hashOfIDs(sortedRoleIds)
	d.Set("bindings_hash", bindHash)
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

	rolesList, err := c.GetWorkspaceGroupRoleList(ctx, &groupUUID, c.WorkspaceUUID)
	if err != nil {
		return diag.Errorf("error reading group role %s: %s", groupID, err)
	}

	remoteRoles := make([]string, 0, len(rolesList))
	for _, r := range rolesList {
		if r.UUID != nil {
			remoteRoles = append(remoteRoles, r.UUID.String())
			continue
		}
		if rr, ok := any(r).(interface{ GetRoleUUID() string }); ok {
			if id := rr.GetRoleUUID(); id != "" {
				remoteRoles = append(remoteRoles, id)
			}
		}
	}
	remoteRoles = uniqueSorted(remoteRoles)

	eff := intersect(toSet(sortedRoleIds), toSet(remoteRoles))
	effective := uniqueSorted(setKeys(eff))

	d.Set("role_ids", effective)

	bindHash := d.Get("bindings_hash").(string)
	if bindHash == "" {
		bindHash = hashOfIDs(effective)
		d.Set("bindings_hash", bindHash)
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

		if err := c.UnbindRoleFromGroup(c.WorkspaceUUID, &u, &groupUUID, map[string]string{}); err != nil {
			return diag.Errorf("unbind role %s from group %s failed: %s", u.String(), groupUUID.String(), err)
		}
	}

	d.SetId("")
	return nil
}
