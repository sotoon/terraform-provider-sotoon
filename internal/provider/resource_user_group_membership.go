package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

// resourceUserGroupMembership defines the schema and CRUD functions for the membership resource.
func resourceUserGroupMembership() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages the membership of a user in one or more IAM groups.",
		CreateContext: resourceUserGroupMembershipCreate,
		ReadContext:   resourceUserGroupMembershipRead,
		DeleteContext: resourceUserGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A stable identifier for this membership binding (group + users).",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The UUID of the group to add users to.",
			},
			"user_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				MinItems:    1,
				Description: "A list of user UUIDs to add to the group.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"bindings_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SHA-256 of sorted, canonical user_ids.",
			},
		},
	}
}

func resourceUserGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	groupID := d.Get("group_id").(string)
	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		return diag.Errorf("invalid group_id: %s", err)
	}

	sortedUserIds := uniqueSorted(fromSchemaSetToStrings(d.Get("user_ids").(*schema.Set)))

	usersList, err := c.GetAllGroupUserList(ctx, &groupUUID)
	if err != nil {
		return diag.Errorf("read group users: %s", err)
	}
	remoteUsersID := make([]string, 0, len(usersList))
	for _, u := range usersList {
		remoteUsersID = append(remoteUsersID, u.Uuid)
	}
	remoteUsersID = uniqueSorted(remoteUsersID)

	toAddList := diff(toSet(sortedUserIds), toSet(remoteUsersID))
	if len(toAddList) > 0 {
		if _, err := c.BulkAddUsersToGroup(ctx, groupUUID, toAddList); err != nil {
			return diag.Errorf("add users to group %s: %s", groupID, err)
		}
	}

	bindHash := hashOfIDs(sortedUserIds)
	if err := d.Set("bindings_hash", bindHash); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set bindings_hash: %w", err))
	}
	d.SetId(groupUUID.String() + ":" + bindHash)

	return resourceUserGroupMembershipRead(ctx, d, meta)
}

func resourceUserGroupMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	groupID := d.Get("group_id").(string)
	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		d.SetId("")
		return nil
	}

	sortedUserIds := uniqueSorted(fromSchemaSetToStrings(d.Get("user_ids").(*schema.Set)))

	usersList, err := c.GetAllGroupUserList(ctx, &groupUUID)
	if err != nil {
		return diag.Errorf("read group users: %s", err)
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

	d.SetId(groupUUID.String() + ":" + bindHash)

	tflog.Info(ctx, "Read user group membership", map[string]interface{}{"group_id": groupID, "have": len(effective)})
	return nil
}

func resourceUserGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	groupID := d.Get("group_id").(string)
	userIDs := d.Get("user_ids").(*schema.Set).List()

	for _, v := range userIDs {
		uid := v.(string)
		if err := c.RemoveUserFromGroup(ctx, groupID, uid); err != nil {
			return diag.Errorf("failed to remove user %s from group %s: %s", uid, groupID, err)
		}
	}
	d.SetId("")
	return nil
}
