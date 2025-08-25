package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

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
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
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
		return diag.Errorf("invalid group_id format: %s", err)
	}

	raw := d.Get("user_ids").([]interface{})
	seen := make(map[string]struct{}, len(raw))
	userIDs := make([]string, 0, len(raw))
	userUUIDs := make([]uuid.UUID, 0, len(raw))

	for _, v := range raw {
		s := v.(string)
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		u, err := uuid.FromString(s)
		if err != nil {
			return diag.Errorf("invalid user_ids entry %q: %s", s, err)
		}
		userIDs = append(userIDs, u.String())
		userUUIDs = append(userUUIDs, u)
	}

	sort.Strings(userIDs)
	h := sha256.Sum256([]byte(strings.Join(userIDs, ",")))
	fullHash := hex.EncodeToString(h[:])
	shortHash := fullHash[:16]

	for i, u := range userUUIDs {
		_, err := c.IAMClient.BulkAddUsersToGroup(*c.WorkspaceUUID, groupUUID, []uuid.UUID{u})
		if err != nil {
			msg := strings.ToLower(err.Error())
			if strings.Contains(msg, "bad request") || strings.Contains(msg, "already") {
				tflog.Warn(ctx, "User may already be a member; continuing", map[string]interface{}{
					"groupID": groupID, "userID": userIDs[i],
				})
				continue
			}
			return diag.Errorf("failed to add user %s to group %s: %s", userIDs[i], groupID, err)
		}
		tflog.Debug(ctx, "Successfully added user to group", map[string]interface{}{"groupID": groupID, "userID": userIDs[i]})
	}

	d.SetId(fmt.Sprintf("%s:%s", groupID, shortHash))
	_ = d.Set("bindings_hash", fullHash)
	return resourceUserGroupMembershipRead(ctx, d, meta)
}

func resourceUserGroupMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	raw := d.Get("user_ids").([]interface{})
	userIDs := make([]string, 0, len(raw))
	for _, v := range raw {
		u, err := uuid.FromString(v.(string))
		if err != nil {
			continue
		}
		userIDs = append(userIDs, u.String())
	}
	sort.Strings(userIDs)
	h := sha256.Sum256([]byte(strings.Join(userIDs, ",")))
	_ = d.Set("bindings_hash", hex.EncodeToString(h[:]))

	tflog.Info(ctx, "Reading user group membership", map[string]interface{}{"id": d.Id()})
	return nil
}

func resourceUserGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	groupID := d.Get("group_id").(string)
	userIDs := d.Get("user_ids").([]interface{})

	for _, v := range userIDs {
		uid := v.(string)
		if err := c.RemoveUserFromGroup(ctx, groupID, uid); err != nil {
			msg := strings.ToLower(err.Error())
			if strings.Contains(msg, "not found") || strings.Contains(msg, "no such") {
				tflog.Warn(ctx, "User or group not found during removal; assuming already deleted", map[string]interface{}{"groupID": groupID, "userID": uid})
				continue
			}
			return diag.Errorf("failed to remove user %s from group %s: %s", uid, groupID, err)
		}
	}
	d.SetId("")
	return nil
}
