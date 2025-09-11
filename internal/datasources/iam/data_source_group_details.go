package iamdatasources

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func DataSourceGroupDetails() *schema.Resource {
	return &schema.Resource{
		Description: "Fetches details of a specific IAM group within a Sotoon workspace.",
		ReadContext: dataSourceGroupDetailsRead,
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the workspace.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the group.",
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
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the group was created (RFC3339).",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the group was last updated (RFC3339).",
			},
			"users_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of users in the group.",
			},
			"service_users_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of service users in the group.",
			},
			"roles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of roles assigned to the group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UUID of the role.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the role.",
						},
					},
				},
			},
		},
	}
}

func dataSourceGroupDetailsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	workspaceID := d.Get("workspace_id").(string)
	groupID := d.Get("group_id").(string)

	workspaceUUID, err := uuid.FromString(workspaceID)
	if err != nil {
		return diag.Errorf("Invalid workspace_id: %s is not a valid UUID", workspaceID)
	}
	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		return diag.Errorf("Invalid group_id: %s is not a valid UUID", groupID)
	}

	group, err := c.GetWorkspaceGroupDetail(ctx, workspaceUUID, groupUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(group.UUID.String())

	d.Set("name", group.Name)
	d.Set("description", group.Description)
	d.Set("created_at", group.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", group.UpdatedAt.Format(time.RFC3339))
	d.Set("users_number", group.UsersNumber)
	d.Set("service_users_number", group.ServiceUsersNumber)

	roles := make([]map[string]interface{}, len(group.Roles))
	for i, r := range group.Roles {
		roles[i] = map[string]interface{}{
			"id":   r.UUID.String(),
			"name": r.Name,
		}
	}
	if err := d.Set("roles", roles); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set roles: %w", err))
	}

	return nil
}
