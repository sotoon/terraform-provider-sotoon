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
				Description: "A unique identifier for the membership resource.",
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The UUID of the user to add to the groups.",
			},
			"group_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				ForceNew:    true,
				Description: "A list of group UUIDs to which the user will be added.",
			},
		},
	}
}

func resourceUserGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	userID := d.Get("user_id").(string)
	groupIds := d.Get("group_ids").([]interface{})

	userUUID, err := uuid.FromString(userID)
	if err != nil {
		return diag.Errorf("invalid user_id format: %s", err)
	}

	for _, groupID := range groupIds {
		groupUUID, err := uuid.FromString(groupID.(string))
		if err != nil {
			return diag.Errorf("invalid group_id format in list: %s", err)
		}

		_, err = c.IAMClient.BulkAddUsersToGroup(*c.WorkspaceUUID, groupUUID, []uuid.UUID{userUUID})
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "bad request") {
				tflog.Warn(ctx, "Received 'bad request' when adding user to group. This may mean the user is already a member.", map[string]interface{}{"userID": userID, "groupID": groupID})
				continue
			}
			return diag.Errorf("failed to add user %s to group %s: %s", userID, groupID, err)
		}
		tflog.Debug(ctx, "Successfully added user to group", map[string]interface{}{"userID": userID, "groupID": groupID})
	}

	// Create a unique ID for this membership resource for state management
	d.SetId(userID + ":" + uuid.NewV4().String())

	return resourceUserGroupMembershipRead(ctx, d, meta)
}

func resourceUserGroupMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Reading user group membership", map[string]interface{}{"id": d.Id()})
	return nil
}

func resourceUserGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    c := meta.(*client.Client)
    userID := d.Get("user_id").(string)
    groupIds := d.Get("group_ids").([]interface{})

    for _, groupID := range groupIds {
        err := c.RemoveUserFromGroup(ctx, groupID.(string), userID)
        if err != nil {
            if strings.Contains(strings.ToLower(err.Error()), "not found") {
                tflog.Warn(ctx, "User or group not found during removal, assuming membership is already deleted.", map[string]interface{}{"userID": userID, "groupID": groupID})
                continue
            }
            return diag.Errorf("failed to remove user %s from group %s: %s", userID, groupID.(string), err)
        }
    }

    d.SetId("")
    return nil
}
