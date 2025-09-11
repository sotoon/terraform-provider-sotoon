package iamdatasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func DataSourceGroupUsers() *schema.Resource {
	return &schema.Resource{
		Description: "Fetches a list of IAM users within a specific Sotoon group and workspace.",
		ReadContext: dataSourceGroupUsersListRead,
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the workspace to fetch users from.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the group to fetch users from.",
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of users belongs to a group found in the workspace .",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for the user.",
						},
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The email address of the user.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the user.",
						},
						"is_suspended": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the user account is suspended.",
						},
					},
				},
			},
		},
	}
}

func dataSourceGroupUsersListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	workspaceID := d.Get("workspace_id").(string)
	groupID := d.Get("group_id").(string)

	workspaceUUID, err := uuid.FromString(workspaceID)
	if err != nil {
		return diag.Errorf("Invalid workspace_id format: not a valid UUID")
	}

	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		return diag.Errorf("Invalid group_id format: not a valid UUID")
	}

	users, err := c.GetWorkspaceGroupUsersList(ctx, &workspaceUUID, &groupUUID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return diag.Errorf(
				"Forbidden: The API token does not have permissions to list groups in this workspace. Please check the token's IAM roles. Original error: %s",
				err,
			)
		}
		return diag.FromErr(err)
	}

	userList := make([]map[string]interface{}, 0, len(users))
	for _, user := range users {
		userData := map[string]interface{}{
			"id":    user.UUID.String(),
			"name":  user.Name,
			"email": user.Email,
		}
		userList = append(userList, userData)
	}

	if err := d.Set("users", userList); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set users list: %w", err))
	}

	d.SetId(workspaceID)

	return nil
}
