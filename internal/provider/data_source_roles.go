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

func dataSourceRoles() *schema.Resource {
	return &schema.Resource{
		Description: "Fetches a list of IAM roles within a specific Sotoon workspace.",
		ReadContext: dataSourceRolesRead,
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the workspace to fetch roles from.",
			},
			"roles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of roles found in the workspace.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for the role.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the role.",
						},
					},
				},
			},
		},
	}
}

func dataSourceRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	workspaceID := d.Get("workspace_id").(string)

	tflog.Debug(ctx, "Reading roles for workspace", map[string]interface{}{"workspace_id": workspaceID})

	if _, err := uuid.FromString(workspaceID); err != nil {
		return diag.Errorf("Invalid workspace_id format: not a valid UUID: %s", err)
	}

	roles, err := c.GetWorkspaceRoles(ctx)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return diag.Errorf(
				"Forbidden: the API token lacks permissions to list roles in this workspace. Original error: %s",
				err,
			)
		}
		return diag.FromErr(err)
	}

	roleList := make([]map[string]interface{}, 0, len(roles))
	for _, role := range roles {
		roleList = append(roleList, map[string]interface{}{
			"id":   role.UUID.String(),
			"name": role.Name,
		})
	}

	if err := d.Set("roles", roleList); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set roles list: %w", err))
	}

	d.SetId(workspaceID)
	return nil
}
