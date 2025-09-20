package iamresources

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
	"github.com/sotoon/terraform-provider-sotoon/internal/common"
)

// ResourceUserGroupMembership defines the schema and CRUD functions for the membership resource.
func ResourceUserGroupMembership() *schema.Resource {
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

	sortedUserIds := common.UniqueSorted(common.FromSchemaSetToStrings(d.Get("user_ids").(*schema.Set)))

	usersList, err := c.GetAllGroupUserList(ctx, &groupUUID)
	if err != nil {
		return diag.Errorf("read group users: %s", err)
	}
	remoteUsersID := make([]string, 0, len(usersList))
	for _, u := range usersList {
		remoteUsersID = append(remoteUsersID, u.UUID.String())
	}
	remoteUsersID = common.UniqueSorted(remoteUsersID)

	toAddList := common.Diff(common.ToSet(sortedUserIds), common.ToSet(remoteUsersID))
	if len(toAddList) > 0 {
		uuids := make([]uuid.UUID, 0, len(toAddList))
		for _, id := range toAddList {
			uuid, _ := uuid.FromString(id)
			uuids = append(uuids, uuid)
		}
		if _, err := c.BulkAddUsersToGroup(ctx, groupUUID, uuids); err != nil {
			return diag.Errorf("add users to group %s: %s", groupID, err)
		}
	}

	bindHash := common.HashOfIDs(sortedUserIds)
	if err := d.Set("bindings_hash", bindHash); err != nil {
		return diag.Errorf("failed to set bindings_hash: %s", err.Error())
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

	sortedUserIds := common.UniqueSorted(common.FromSchemaSetToStrings(d.Get("user_ids").(*schema.Set)))

	usersList, err := c.GetAllGroupUserList(ctx, &groupUUID)
	if err != nil {
		return diag.Errorf("read group users: %s", err)
	}
	remoteUsersID := make([]string, 0, len(usersList))
	for _, u := range usersList {
		remoteUsersID = append(remoteUsersID, u.UUID.String())
	}
	remoteUsersID = common.UniqueSorted(remoteUsersID)
	eff := common.Intersect(common.ToSet(sortedUserIds), common.ToSet(remoteUsersID))
	effective := common.UniqueSorted(common.SetKeys(eff))
	if err := d.Set("user_ids", effective); err != nil {
		return diag.Errorf("failed to set user_ids: %s", err.Error())
	}

	bindHash := ""
	if v, ok := d.GetOk("bindings_hash"); ok && v.(string) != "" {
		bindHash = v.(string)
	} else {
		bindHash = common.HashOfIDs(effective)
		if err := d.Set("bindings_hash", bindHash); err != nil {
			return diag.Errorf("failed to set bindings_hash: %s", err.Error())
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
			msg := strings.ToLower(err.Error())
			if strings.Contains(msg, "not found") || strings.Contains(msg, "no such") {
				tflog.Warn(ctx, "User or group not found during removal; assuming already deleted",
					map[string]interface{}{"groupID": groupID, "userID": uid})
				continue
			}
			return diag.Errorf("failed to remove user %s from group %s: %s", uid, groupID, err)
		}
	}
	d.SetId("")
	return nil
}
