package provider

import (
	"context"
	"fmt"

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
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Set of user UUIDs to bind to the role.",
			},
			"items": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "map of items related to this role.",
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

	roleID := d.Get("role_id").(string)
	roleUUID, err := uuid.FromString(roleID)
	if err != nil {
		return diag.Errorf("invalid role_id: %s", err)
	}

	sortedUserIds := uniqueSorted(fromSchemaSetToStrings(d.Get("user_ids").(*schema.Set)))

	usersList, err := c.GetRoleUsers(ctx, &roleUUID)
	if err != nil {
		return diag.Errorf("read users of role: %s", err)
	}
	remoteUsersID := make([]string, 0, len(usersList))
	for _, u := range usersList {
		remoteUsersID = append(remoteUsersID, u.Uuid)
	}
	remoteUsersID = uniqueSorted(remoteUsersID)

	toAddList := diff(toSet(sortedUserIds), toSet(remoteUsersID))
	if len(toAddList) > 0 {
		var itemsToAdd map[string]interface{}
		if items, found := d.GetOk("items"); found {
			itemsToAdd = items.(map[string]interface{})
		}
		if err := c.BulkAddUsersToRole(ctx, roleUUID, toAddList, itemsToAdd); err != nil {
			return diag.Errorf("add users to role %s: %s", roleUUID, err)
		}
	}

	bindHash := hashOfIDs(sortedUserIds)

	if err := d.Set("bindings_hash", bindHash); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set bindings_hash: %w", err))
	}
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

	usersList, err := c.GetRoleUsers(ctx, &roleUUID)
	if err != nil {
		return diag.Errorf("read users of role %s: %s", roleUUID, err)
	}

	remoteUsersID := make([]string, 0, len(usersList))
	for _, u := range usersList {
		remoteUsersID = append(remoteUsersID, u.Uuid)
	}
	remoteUsersID = uniqueSorted(remoteUsersID)

	eff := intersect(toSet(sortedUserIds), toSet(remoteUsersID))
	effective := uniqueSorted(setKeys(eff))

	if err := d.Set("user_ids", effective); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set user_ids: %w", err))
	}

	bindHash := ""
	if v, ok := d.GetOk("bindings_hash"); ok && v.(string) != "" {
		bindHash = v.(string)
	} else {
		bindHash = hashOfIDs(effective)
		if err := d.Set("bindings_hash", bindHash); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set bindings_hash: %w", err))
		}

	}

	d.SetId(roleUUID.String() + ":" + bindHash)

	return nil
}

func resourceUserRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	roleID := d.Get("role_id").(string)
	roleUUID, err := uuid.FromString(roleID)
	if err != nil {
		return diag.Errorf("invalid role_id %q: %s", roleID, err)
	}

	users := d.Get("user_ids").(*schema.Set).List()
	for _, v := range users {
		u, err := uuid.FromString(v.(string))
		if err != nil {
			return diag.Errorf("invalid user_id in list: %s", err)
		}
		if err := c.UnbindRoleFromUser(ctx, &roleUUID, &u); err != nil {
			return diag.Errorf("unbind user %s from role %s failed: %s", u.String(), roleUUID.String(), err)
		}
	}
	d.SetId("")
	return nil
}
