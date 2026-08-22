package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	iam "github.com/sotoon/sotoon-sdk-go/sdk/core/iam_v1"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

type userRoleBindingClient interface {
	GetRoleUsers(context.Context, *uuid.UUID) ([]iam.IamUserWithRoleItems, error)
	BulkAddUsersToRole(context.Context, uuid.UUID, []string, map[string]string) error
	UnbindRoleFromUser(context.Context, *uuid.UUID, *uuid.UUID) error
}

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
				Description: "SHA-256 of sorted user_ids and items. Changes when membership or items change.",
			},
		},
	}
}

func resourceUserRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	roleBindingItemsMu.Lock()
	defer roleBindingItemsMu.Unlock()
	return resourceUserRoleCreateWithClient(ctx, d, meta.(*client.Client))
}

func resourceUserRoleCreateWithClient(ctx context.Context, d *schema.ResourceData, c userRoleBindingClient) diag.Diagnostics {
	roleID := d.Get("role_id").(string)
	roleUUID, err := uuid.FromString(roleID)
	if err != nil {
		return diag.Errorf("invalid role_id: %s", err)
	}

	sortedUserIDs := uniqueSorted(fromSchemaSetToStrings(d.Get("user_ids").(*schema.Set)))
	wantedItems := roleItemsFromResourceData(d)

	users, err := c.GetRoleUsers(ctx, &roleUUID)
	if err != nil {
		return diag.Errorf("read users of role: %s", err)
	}
	remoteItems := userRoleItemsByID(users)

	toAdd := make([]string, 0, len(sortedUserIDs))
	for _, userID := range sortedUserIDs {
		items, exists := remoteItems[userID]
		if !roleBindingHasItems(exists, items, wantedItems) {
			toAdd = append(toAdd, userID)
		}
	}
	if len(toAdd) > 0 {
		if err := c.BulkAddUsersToRole(ctx, roleUUID, toAdd, wantedItems); err != nil {
			return diag.Errorf("add users to role %s with requested items: %s", roleUUID, err)
		}
	}

	return resourceUserRoleReadWithClient(ctx, d, c)
}

func resourceUserRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	roleBindingItemsMu.Lock()
	defer roleBindingItemsMu.Unlock()
	return resourceUserRoleReadWithClient(ctx, d, meta.(*client.Client))
}

func resourceUserRoleReadWithClient(ctx context.Context, d *schema.ResourceData, c userRoleBindingClient) diag.Diagnostics {
	roleID := d.Get("role_id").(string)
	roleUUID, err := uuid.FromString(roleID)
	if err != nil {
		d.SetId("")
		return nil
	}

	sortedUserIDs := uniqueSorted(fromSchemaSetToStrings(d.Get("user_ids").(*schema.Set)))
	wantedItems := roleItemsFromResourceData(d)

	users, err := c.GetRoleUsers(ctx, &roleUUID)
	if err != nil {
		return diag.Errorf("read users of role %s: %s", roleUUID, err)
	}
	remoteItems := userRoleItemsByID(users)

	effective := make([]string, 0, len(sortedUserIDs))
	for _, userID := range sortedUserIDs {
		items, exists := remoteItems[userID]
		if roleBindingHasItems(exists, items, wantedItems) {
			effective = append(effective, userID)
		}
	}
	if err := d.Set("user_ids", effective); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set user_ids: %w", err))
	}

	bindingHash := hashOfBinding(sortedUserIDs, wantedItems)
	if err := d.Set("bindings_hash", bindingHash); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set bindings_hash: %w", err))
	}
	d.SetId(roleUUID.String() + ":" + bindingHash)
	return nil
}

func resourceUserRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	roleBindingItemsMu.Lock()
	defer roleBindingItemsMu.Unlock()
	return resourceUserRoleDeleteWithClient(ctx, d, meta.(*client.Client))
}

func resourceUserRoleDeleteWithClient(ctx context.Context, d *schema.ResourceData, c userRoleBindingClient) diag.Diagnostics {
	roleID := d.Get("role_id").(string)
	roleUUID, err := uuid.FromString(roleID)
	if err != nil {
		return diag.Errorf("invalid role_id %q: %s", roleID, err)
	}

	users, err := c.GetRoleUsers(ctx, &roleUUID)
	if err != nil {
		return diag.Errorf("read users of role %s before delete: %s", roleUUID, err)
	}
	remoteItems := userRoleItemsByID(users)
	wantedItems := roleItemsFromResourceData(d)

	for _, userID := range uniqueSorted(fromSchemaSetToStrings(d.Get("user_ids").(*schema.Set))) {
		items, exists := remoteItems[userID]
		remaining, found := removeRoleBindingItems(exists, items, wantedItems)
		if !found {
			continue
		}

		userUUID, err := uuid.FromString(userID)
		if err != nil {
			return diag.Errorf("invalid user_id %q: %s", userID, err)
		}
		if err := c.UnbindRoleFromUser(ctx, &roleUUID, &userUUID); err != nil {
			return diag.Errorf("unbind user %s from role %s failed: %s", userID, roleUUID, err)
		}
		if err := restoreUserRoleItems(ctx, c, roleUUID, userID, remaining); err != nil {
			rollbackErr := restoreUserRoleItems(ctx, c, roleUUID, userID, items)
			if rollbackErr != nil {
				return diag.Errorf("restore remaining items for user %s and role %s failed: %s; rollback failed: %s", userID, roleUUID, err, rollbackErr)
			}
			return diag.Errorf("restore remaining items for user %s and role %s failed: %s; original binding restored", userID, roleUUID, err)
		}
	}

	d.SetId("")
	return nil
}

func userRoleItemsByID(users []iam.IamUserWithRoleItems) map[string][]map[string]string {
	itemsByID := make(map[string][]map[string]string, len(users))
	for _, user := range users {
		itemsByID[user.Uuid] = append(itemsByID[user.Uuid], user.Items...)
	}
	return itemsByID
}

func restoreUserRoleItems(ctx context.Context, c userRoleBindingClient, roleUUID uuid.UUID, userID string, items []map[string]string) error {
	for _, itemSet := range items {
		if err := c.BulkAddUsersToRole(ctx, roleUUID, []string{userID}, itemSet); err != nil {
			return err
		}
	}
	return nil
}
