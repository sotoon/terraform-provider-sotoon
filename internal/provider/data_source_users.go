package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		Description: "Fetches a list of IAM users within a specific Sotoon workspace.",
		ReadContext: dataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the workspace to fetch users from.",
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of users found in the workspace.",
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

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	workspaceID := d.Get("workspace_id").(string)

	tflog.Debug(ctx, "Reading users for workspace", map[string]interface{}{"workspace_id": workspaceID})

	workspaceUUID, err := uuid.FromString(workspaceID)
	if err != nil {
		return diag.Errorf("Invalid workspace_id format: not a valid UUID")
	}

	users, err := c.GetWorkspaceUsers(ctx, &workspaceUUID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return diag.Errorf(
				"Forbidden: The API token does not have permissions to list users in this workspace. Please check the token's IAM roles. Original error: %s",
				err,
			)
		}
		return diag.FromErr(err)
	}

	userList := make([]map[string]interface{}, 0, len(users))
	for _, user := range users {
		userData := map[string]interface{}{
			"id":    user.Uuid,
			"email": user.Email,
			"name":  user.Name,
			// "is_suspended": user.IsSuspended,
		}
		userList = append(userList, userData)
	}

	if err := d.Set("users", userList); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set users list: %w", err))
	}

	d.SetId(createDataSourceUsers_ID(workspaceID))

	return nil
}

func createDataSourceUsers_ID(workspaceID string) string {
	return "users_" + workspaceID
}
