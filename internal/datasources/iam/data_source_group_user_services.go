package iamdatasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func DataSourceGroupUserServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupUserServicesListRead,
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"roles": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceGroupUserServicesListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	ws := d.Get("workspace_id").(string)
	grp := d.Get("group_id").(string)
	workspaceUUID, err := uuid.FromString(ws)
	if err != nil {
		return diag.Errorf("Invalid workspace_id: %s is not a valid UUID", ws)
	}
	groupUUID, err := uuid.FromString(grp)
	if err != nil {
		return diag.Errorf("Invalid group_id: %s is not a valid UUID", grp)
	}
	users, err := c.GetAllGroupServiceUserList(ctx, &workspaceUUID, &groupUUID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(groupUUID.String())
	list := make([]map[string]interface{}, len(users))
	for i, u := range users {
		roleList := make([]map[string]interface{}, len(u.Roles))
		for j, r := range u.Roles {
			roleList[j] = map[string]interface{}{
				"id":   r.UUID.String(),
				"name": r.Name,
			}
		}
		entry := map[string]interface{}{
			"id":          u.UUID.String(),
			"name":        u.Name,
			"description": u.Description,
			"workspace":   "",
			"created_at":  u.CreatedAt,
			"updated_at":  u.UpdatedAt,
			"roles":       roleList,
		}
		if u.Workspace != nil {
			entry["workspace"] = u.Workspace.String()
		}
		list[i] = entry
	}
	if err := d.Set("service_users", list); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set service_users: %w", err))
	}
	return nil
}
