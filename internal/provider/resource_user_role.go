package provider

import (
	"context"

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
		UpdateContext: resourceUserRoleUpdate,
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
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Set of user UUIDs to bind to the role."},

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

	roleID := d.Get("role_id").(string)
	roleUUID, err := uuid.FromString(roleID)
	if err != nil {
		return diag.Errorf("invalid role_id: %s", err)
	}

	sortedUserIds := uniqueSorted(fromSchemaSetToStrings(d.Get("user_ids").(*schema.Set)))

	usersList, err := c.IAMClient.GetRoleUsers(&roleUUID, c.WorkspaceUUID)
	if err != nil {
		return diag.Errorf("read users of role: %s", err)
	}
	remoteUsersID := make([]string, 0, len(usersList))
	for _, u := range usersList {
		remoteUsersID = append(remoteUsersID, u.UUID.String())
	}
	remoteUsersID = uniqueSorted(remoteUsersID)

	toAddList := diff(toSet(sortedUserIds), toSet(remoteUsersID))
	if len(toAddList) > 0 {
		uuids := make([]uuid.UUID, 0, len(toAddList))
		for _, id := range toAddList {
			uuid, _ := uuid.FromString(id)
			uuids = append(uuids, uuid)
		}
		if err := c.IAMClient.BulkAddUsersToRole(*c.WorkspaceUUID, roleUUID, uuids); err != nil {
			return diag.Errorf("add users to role %s: %s", roleUUID, err)
		}
	}

	bindHash := hashOfIDs(sortedUserIds)
	d.Set("bindings_hash", bindHash)
	d.SetId(roleUUID.String() + ":" + bindHash)

	return resourceUserRoleRead(ctx, d, meta)
}

func resourceUserRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	roleID := d.Get("role_id").(string)
	roleUUID, err := uuid.FromString(roleID)
	if err != nil {
		d.SetId("")
		return nil
	}

	sortedUserIds := uniqueSorted(fromSchemaSetToStrings(d.Get("user_ids").(*schema.Set)))

	usersList, err := c.IAMClient.GetRoleUsers(&roleUUID, c.WorkspaceUUID)
	if err != nil {
		return diag.Errorf("read users of role %s: %s", roleUUID, err)
	}

	remoteUsersID := make([]string, 0, len(usersList))
	for _, u := range usersList {
		remoteUsersID = append(remoteUsersID, u.UUID.String())
	}
	remoteUsersID = uniqueSorted(remoteUsersID)

	eff := intersect(toSet(sortedUserIds), toSet(remoteUsersID))
	effective := uniqueSorted(setKeys(eff))

	d.Set("user_ids", effective)

	bindHash := ""
	if v, ok := d.GetOk("bindings_hash"); ok && v.(string) != "" {
		bindHash = v.(string)
	} else {
		bindHash = hashOfIDs(effective)
		d.Set("bindings_hash", bindHash)
	}

	d.SetId(roleUUID.String() + ":" + bindHash)

	return nil
}

func resourceUserRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange("user_ids") {
		return nil
	}
	return resourceUserRoleCreate(ctx, d, meta)
}

func resourceUserRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	roleID := d.Get("role_id").(string)
	roleUUID, err := uuid.FromString(roleID)
	if err != nil {
		return diag.Errorf("invalid role_id %q: %s", roleID, err)
	}

	raw := d.Get("user_ids").(*schema.Set).List()
	for _, v := range raw {
		u, err := uuid.FromString(v.(string))
		if err != nil {
			return diag.Errorf("invalid user_id in list: %s", err)
		}
		if err := c.IAMClient.UnbindRoleFromUser(c.WorkspaceUUID, &roleUUID, &u, map[string]string{}); err != nil {
			return diag.Errorf("unbind user %s from role %s failed: %s", u.String(), roleUUID.String(), err)
		}
	}
	d.SetId("")
	return nil
}
