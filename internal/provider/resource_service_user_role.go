package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"

	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func resourceServiceUserRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceUserRoleCreate,
		ReadContext:   resourceServiceUserRoleRead,
		DeleteContext: resourceServiceUserRoleDelete,
		Schema: map[string]*schema.Schema{
			"id": {Type: schema.TypeString, Computed: true},
			"service_user_ids": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"role_id": {Type: schema.TypeString, Required: true, ForceNew: true},
			"bindings_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SHA-256 of sorted, canonical service_user_ids.",
			},
		},
	}
}

func resourceServiceUserRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	roleStr := d.Get("role_id").(string)
	roleID, err := uuid.FromString(roleStr)
	if err != nil {
		return diag.Errorf("invalid role_id format: %s", err)
	}

	input := fromSchemaSetToStrings(d.Get("service_user_ids").(*schema.Set))
	converted, err := convertStringsToUUIDArray(input)
	if err != nil {
		return diag.Errorf("invalid service_user_ids: %s", err)
	}
	ids := uniqueSorted(converted)
	_ = d.Set("bindings_hash", hashOfIDs(ids))

	suUUIDs := make([]uuid.UUID, 0, len(ids))
	for _, s := range ids {
		u, _ := uuid.FromString(s)
		suUUIDs = append(suUUIDs, u)
	}

	if len(suUUIDs) > 0 {
		if err := c.IAMClient.BulkAddServiceUsersToRole(*c.WorkspaceUUID, roleID, suUUIDs); err != nil {
			low := strings.ToLower(err.Error())
			if strings.Contains(low, "already") || strings.Contains(low, "exists") || strings.Contains(low, "conflict") {
				tflog.Warn(ctx, "Service users may already have this role",
					map[string]interface{}{"roleID": roleID, "serviceUserIDs": uuidsToStringSlice(suUUIDs)})
			} else {
				return diag.Errorf("failed to add service users to role %s: %s", roleID, err)
			}
		}
		tflog.Debug(ctx, "Added service users to role",
			map[string]interface{}{"roleID": roleID, "serviceUserIDs": uuidsToStringSlice(suUUIDs)})
	}

	d.SetId(roleID.String() + ":" + hashOfIDs(ids)[:16])
	return resourceServiceUserRoleRead(ctx, d, meta)
}

func resourceServiceUserRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := fromSchemaSetToStrings(d.Get("service_user_ids").(*schema.Set))
	converted, err := convertStringsToUUIDArray(input)
	if err == nil {
		_ = d.Set("bindings_hash", hashOfIDs(uniqueSorted(converted)))
	}
	tflog.Info(ctx, "Reading service user role", map[string]interface{}{"id": d.Id()})
	return nil
}

func resourceServiceUserRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	roleStr := d.Get("role_id").(string)
	roleID, err := uuid.FromString(roleStr)
	if err != nil {
		return diag.Errorf("invalid role_id %q: %s", roleStr, err)
	}

	raw := d.Get("service_user_ids").(*schema.Set).List()
	for _, v := range raw {
		u, err := uuid.FromString(v.(string))
		if err != nil {
			return diag.Errorf("invalid service_user_id in list: %s", err)
		}
		if err := c.UnbindRoleFromServiceUser(c.WorkspaceUUID, &roleID, &u, map[string]string{}); err != nil {
			return diag.Errorf("unbind service user %s from role %s failed: %s", u.String(), roleID.String(), err)
		}
	}
	d.SetId("")
	return nil
}
