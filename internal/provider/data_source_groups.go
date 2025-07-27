package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
	uuid "github.com/satori/go.uuid"
)

func dataSourceGroups() *schema.Resource {
	return &schema.Resource{
		Description: "Fetches a list of IAM groups within a specific Sotoon workspace.",
		ReadContext: dataSourceGroupsRead,
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the workspace to fetch groups from.",
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of groups found in the workspace.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for the group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the group.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the group.",
						},
					},
				},
			},
		},
	}
}

func dataSourceGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	workspaceID := d.Get("workspace_id").(string)

	tflog.Debug(ctx, "Reading groups for workspace", map[string]interface{}{"workspace_id": workspaceID})

	workspaceUUID, err := uuid.FromString(workspaceID)
	if err != nil {
		return diag.Errorf("Invalid workspace_id format: not a valid UUID")
	}

	groups, err := c.GetWorkspaceGroups(ctx, &workspaceUUID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return diag.Errorf(
				"Forbidden: The API token does not have permissions to list groups in this workspace. Please check the token's IAM roles. Original error: %s",
				err,
			)
		}
		return diag.FromErr(err)
	}

	groupList := make([]map[string]interface{}, 0, len(groups))
	for _, group := range groups {
		groupData := map[string]interface{}{
			"id":          group.UUID.String(),
			"name":        group.Name,
		}
		groupList = append(groupList, groupData)
	}

	if err := d.Set("groups", groupList); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set groups list: %w", err))
	}

	d.SetId(workspaceID)

	return nil
}
