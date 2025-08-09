package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func dataSourceGroupRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRolesRead,
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
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
						"description_fa": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description_en": {
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
					},
				},
			},
		},
	}
}

func dataSourceGroupRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	roles, err := c.GetWorkspaceGroupRoleList(ctx, &workspaceUUID, &groupUUID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(groupUUID.String())
	list := make([]map[string]interface{}, len(roles))
	for i, r := range roles {
		list[i] = map[string]interface{}{
			"id":             r.UUID.String(),
			"name":           r.Name,
			"description_fa": r.DescriptionFa,
			"description_en": r.DescriptionEn,
			"created_at":     r.CreatedAt.Format(time.RFC3339),
			"updated_at":     r.UpdatedAt.Format(time.RFC3339),
		}
	}
	if err := d.Set("roles", list); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set roles: %w", err))
	}
	return nil
}
