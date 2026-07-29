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

type groupRoleBindingClient interface {
	GetWorkspaceGroupRoleList(context.Context, *uuid.UUID, *uuid.UUID) ([]iam.IamRole, error)
	BulkAddRolesToGroup(context.Context, *uuid.UUID, []iam.IamRoleItem) error
	UnbindRoleFromGroup(context.Context, *uuid.UUID, *uuid.UUID) error
}

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
				Description: "SHA-256 of sorted role_ids and items. Changes when roles or items change.",
			},
		},
	}
}

func resourceGroupRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	roleBindingItemsMu.Lock()
	defer roleBindingItemsMu.Unlock()
	c := meta.(*client.Client)
	return resourceGroupRoleCreateWithClient(ctx, d, c, c.WorkspaceUUID)
}

func resourceGroupRoleCreateWithClient(ctx context.Context, d *schema.ResourceData, c groupRoleBindingClient, workspaceUUID *uuid.UUID) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		return diag.Errorf("invalid group_id %q: %s", groupID, err)
	}

	sortedRoleIDs := uniqueSorted(fromSchemaSetToStrings(d.Get("role_ids").(*schema.Set)))
	wantedItems := roleItemsFromResourceData(d)

	roles, err := c.GetWorkspaceGroupRoleList(ctx, workspaceUUID, &groupUUID)
	if err != nil {
		return diag.Errorf("read group roles: %s", err)
	}
	remoteItems := groupRoleItemsByID(roles)

	toAdd := make([]iam.IamRoleItem, 0, len(sortedRoleIDs))
	for _, roleID := range sortedRoleIDs {
		items, exists := remoteItems[roleID]
		if !roleBindingHasItems(exists, items, wantedItems) {
			itemSets := []map[string]string{wantedItems}
			toAdd = append(toAdd, iam.IamRoleItem{RoleUuid: roleID, ItemsList: &itemSets})
		}
	}
	if len(toAdd) > 0 {
		if err := c.BulkAddRolesToGroup(ctx, &groupUUID, toAdd); err != nil {
			return diag.Errorf("bind requested role items to group %s failed: %s", groupUUID, err)
		}
	}

	return resourceGroupRoleReadWithClient(ctx, d, c, workspaceUUID)
}

func resourceGroupRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	roleBindingItemsMu.Lock()
	defer roleBindingItemsMu.Unlock()
	c := meta.(*client.Client)
	return resourceGroupRoleReadWithClient(ctx, d, c, c.WorkspaceUUID)
}

func resourceGroupRoleReadWithClient(ctx context.Context, d *schema.ResourceData, c groupRoleBindingClient, workspaceUUID *uuid.UUID) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		d.SetId("")
		return nil
	}

	sortedRoleIDs := uniqueSorted(fromSchemaSetToStrings(d.Get("role_ids").(*schema.Set)))
	wantedItems := roleItemsFromResourceData(d)

	roles, err := c.GetWorkspaceGroupRoleList(ctx, workspaceUUID, &groupUUID)
	if err != nil {
		return diag.Errorf("read group roles: %s", err)
	}
	remoteItems := groupRoleItemsByID(roles)

	effective := make([]string, 0, len(sortedRoleIDs))
	for _, roleID := range sortedRoleIDs {
		items, exists := remoteItems[roleID]
		if roleBindingHasItems(exists, items, wantedItems) {
			effective = append(effective, roleID)
		}
	}
	if err := d.Set("role_ids", effective); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set role_ids: %w", err))
	}

	bindingHash := hashOfBinding(sortedRoleIDs, wantedItems)
	if err := d.Set("bindings_hash", bindingHash); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set bindings_hash: %w", err))
	}
	d.SetId(groupUUID.String() + ":" + bindingHash)
	return nil
}

func resourceGroupRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	roleBindingItemsMu.Lock()
	defer roleBindingItemsMu.Unlock()
	c := meta.(*client.Client)
	return resourceGroupRoleDeleteWithClient(ctx, d, c, c.WorkspaceUUID)
}

func resourceGroupRoleDeleteWithClient(ctx context.Context, d *schema.ResourceData, c groupRoleBindingClient, workspaceUUID *uuid.UUID) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		return diag.Errorf("invalid group_id %q: %s", groupID, err)
	}

	roles, err := c.GetWorkspaceGroupRoleList(ctx, workspaceUUID, &groupUUID)
	if err != nil {
		return diag.Errorf("read group roles before delete: %s", err)
	}
	remoteItems := groupRoleItemsByID(roles)
	wantedItems := roleItemsFromResourceData(d)

	for _, roleID := range uniqueSorted(fromSchemaSetToStrings(d.Get("role_ids").(*schema.Set))) {
		items, exists := remoteItems[roleID]
		remaining, found := removeRoleBindingItems(exists, items, wantedItems)
		if !found {
			continue
		}

		roleUUID, err := uuid.FromString(roleID)
		if err != nil {
			return diag.Errorf("invalid role_id %q: %s", roleID, err)
		}
		if err := c.UnbindRoleFromGroup(ctx, &roleUUID, &groupUUID); err != nil {
			return diag.Errorf("unbind role %s from group %s failed: %s", roleUUID, groupUUID, err)
		}
		if err := restoreGroupRoleItems(ctx, c, groupUUID, roleID, remaining); err != nil {
			rollbackErr := restoreGroupRoleItems(ctx, c, groupUUID, roleID, items)
			if rollbackErr != nil {
				return diag.Errorf("restore remaining items for role %s and group %s failed: %s; rollback failed: %s", roleID, groupUUID, err, rollbackErr)
			}
			return diag.Errorf("restore remaining items for role %s and group %s failed: %s; original binding restored", roleID, groupUUID, err)
		}
	}

	d.SetId("")
	return nil
}

func groupRoleItemsByID(roles []iam.IamRole) map[string][]map[string]string {
	itemsByID := make(map[string][]map[string]string, len(roles))
	for _, role := range roles {
		itemsByID[role.Uuid] = append(itemsByID[role.Uuid], role.Items...)
	}
	return itemsByID
}

func restoreGroupRoleItems(ctx context.Context, c groupRoleBindingClient, groupUUID uuid.UUID, roleID string, items []map[string]string) error {
	if len(items) == 0 {
		return nil
	}
	itemSets := append([]map[string]string(nil), items...)
	return c.BulkAddRolesToGroup(ctx, &groupUUID, []iam.IamRoleItem{{RoleUuid: roleID, ItemsList: &itemSets}})
}
