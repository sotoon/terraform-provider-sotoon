package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"

	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func resourceServiceUserRole() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a role for a service user within a Sotoon workspace.",
		CreateContext: resourceServiceUserRoleCreate,
		ReadContext:   resourceServiceUserRoleRead,
		DeleteContext: resourceServiceUserRoleDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Composite stable identifier. Does not affect lifecycle.`,
			},
			"service_user_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				ForceNew:    true,
				Description: "List of service user UUIDs to bind to the role.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"role_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Role UUID.",
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
				Description: "SHA-256 of sorted, canonical service_user_ids.",
			},
		},
	}
}

func resourceServiceUserRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	roleID := d.Get("role_id").(string)
	roleUUID, err := uuid.FromString(roleID)
	if err != nil {
		return diag.Errorf("invalid role_id: %s", err)
	}

	sortedServiceUserIds := uniqueSorted(fromSchemaSetToStrings(d.Get("service_user_ids").(*schema.Set)))

	serviceUsersList, err := c.GetRoleServiceUsers(ctx, &roleUUID)
	if err != nil {
		return diag.Errorf("read service-users of role: %s", err)
	}

	remoteServiceUsersID := make([]string, 0, len(serviceUsersList))
	for _, u := range serviceUsersList {
		remoteServiceUsersID = append(remoteServiceUsersID, u.Uuid)
	}

	remoteServiceUsersID = uniqueSorted(remoteServiceUsersID)
	toAddList := diff(toSet(sortedServiceUserIds), toSet(remoteServiceUsersID))
	if len(toAddList) > 0 {
		var itemsToAdd map[string]interface{}
		if items, found := d.GetOk("items"); found {
			itemsToAdd = items.(map[string]interface{})
		}
		if err := c.BulkAddServiceUsersToRole(ctx, roleUUID, toAddList, itemsToAdd); err != nil {
			return diag.Errorf("add service users to role %s: %s", roleUUID, err)
		}
	}

	bindHash := hashOfIDs(sortedServiceUserIds)
	if err := d.Set("bindings_hash", bindHash); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set bindings_hash: %w", err))
	}
	d.SetId(roleUUID.String() + ":" + bindHash)

	return resourceServiceUserRoleRead(ctx, d, meta)
}

func resourceServiceUserRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	roleID := d.Get("role_id").(string)
	roleUUID, err := uuid.FromString(roleID)
	if err != nil {
		d.SetId("")
		return nil
	}

	sortedServiceUserIds := uniqueSorted(fromSchemaSetToStrings(d.Get("service_user_ids").(*schema.Set)))

	serviceUsersList, err := c.GetRoleServiceUsers(ctx, &roleUUID)
	if err != nil {
		return diag.Errorf("read service-users of role %s: %s", roleUUID, err)
	}

	remoteServiceUsersID := make([]string, 0, len(serviceUsersList))
	for _, u := range serviceUsersList {
		remoteServiceUsersID = append(remoteServiceUsersID, u.Uuid)
	}
	remoteServiceUsersID = uniqueSorted(remoteServiceUsersID)

	eff := intersect(toSet(sortedServiceUserIds), toSet(remoteServiceUsersID))
	effective := uniqueSorted(setKeys(eff))

	if err := d.Set("service_user_ids", effective); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set service_user_ids: %w", err))
	}

	bindHash, _ := d.Get("bindings_hash").(string)
	if bindHash == "" {
		bindHash = hashOfIDs(effective)
		if err := d.Set("bindings_hash", bindHash); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set bindings_hash: %w", err))
		}
	}

	d.SetId(roleUUID.String() + ":" + bindHash)
	return nil
}

func resourceServiceUserRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	roleID := d.Get("role_id").(string)
	roleUUID, err := uuid.FromString(roleID)
	if err != nil {
		return diag.Errorf("invalid role_id %q: %s", roleID, err)
	}

	raw := d.Get("service_user_ids").(*schema.Set).List()
	for _, v := range raw {
		u, err := uuid.FromString(v.(string))
		if err != nil {
			return diag.Errorf("invalid service_user_id in list: %s", err)
		}
		if err := c.UnbindRoleFromServiceUser(ctx, &roleUUID, &u); err != nil {
			return diag.Errorf("unbind service user %s from role %s failed: %s", u.String(), roleUUID.String(), err)
		}
	}
	d.SetId("")
	return nil
}
