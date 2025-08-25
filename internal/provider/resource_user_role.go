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

func resourceUserRole() *schema.Resource {
	return &schema.Resource{
		Description:   "Binds one or more users to a role within a Sotoon workspace (bulk).",
		CreateContext: resourceUserRoleCreate,
		ReadContext:   resourceUserRoleRead,
		DeleteContext: resourceUserRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Stable identifier (anchor + hash).`,
			},
			"role_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Role UUID.",
			},
			"user_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Set of user UUIDs to bind to the role.",
			},
			"bindings_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SHA-256 of sorted user_ids. Changes when membership changes.",
			},
		},
	}
}

func resourceUserRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	roleStr := d.Get("role_id").(string)
	roleID, err := uuid.FromString(roleStr)
	if err != nil {
		return diag.Errorf("invalid role_id format: %s", err)
	}

	input := fromSchemaSetToStrings(d.Get("user_ids").(*schema.Set))
	converted, err := convertStringsToUUIDArray(input)
	if err != nil {
		return diag.Errorf("invalid user_ids: %s", err)
	}
	ids := uniqueSorted(converted)
	hash := hashOfIDs(ids)
	_ = d.Set("bindings_hash", hash)

	userUUIDs := make([]uuid.UUID, 0, len(ids))
	for _, s := range ids {
		u, _ := uuid.FromString(s)
		userUUIDs = append(userUUIDs, u)
	}

	if err := c.IAMClient.BulkAddUsersToRole(*c.WorkspaceUUID, roleID, userUUIDs); err != nil {
		low := strings.ToLower(err.Error())
		if strings.Contains(low, "already") || strings.Contains(low, "exists") || strings.Contains(low, "conflict") {
			tflog.Warn(ctx, "Users may already have this role",
				map[string]interface{}{"roleID": roleID, "userIDs": uuidsToStringSlice(userUUIDs)})
		} else {
			return diag.Errorf("failed to add users to role %s: %s", roleID, err)
		}
	}
	tflog.Debug(ctx, "Added users to role",
		map[string]interface{}{"roleID": roleID, "userIDs": uuidsToStringSlice(userUUIDs)})

	d.SetId(roleID.String() + ":" + hash[:16])
	return resourceUserRoleRead(ctx, d, meta)
}

func resourceUserRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := fromSchemaSetToStrings(d.Get("user_ids").(*schema.Set))
	converted, err := convertStringsToUUIDArray(input)
	if err == nil {
		hash := hashOfIDs(uniqueSorted(converted))
		_ = d.Set("bindings_hash", hash)
	}
	tflog.Info(ctx, "Reading user-role bulk binding", map[string]interface{}{"id": d.Id()})
	return nil
}

func resourceUserRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	roleStr := d.Get("role_id").(string)
	roleID, err := uuid.FromString(roleStr)
	if err != nil {
		return diag.Errorf("invalid role_id %q: %s", roleStr, err)
	}

	raw := d.Get("user_ids").(*schema.Set).List()
	for _, v := range raw {
		u, err := uuid.FromString(v.(string))
		if err != nil {
			return diag.Errorf("invalid user_id in list: %s", err)
		}
		if err := c.IAMClient.UnbindRoleFromUser(c.WorkspaceUUID, &roleID, &u, map[string]string{}); err != nil {
			return diag.Errorf("unbind user %s from role %s failed: %s", u.String(), roleID.String(), err)
		}
	}
	d.SetId("")
	return nil
}
