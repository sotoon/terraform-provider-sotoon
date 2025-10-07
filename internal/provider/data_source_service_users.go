package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func dataSourceServiceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServiceUsersRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":          {Type: schema.TypeString, Computed: true},
						"name":        {Type: schema.TypeString, Computed: true},
						"description": {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceServiceUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	list, err := c.GetServiceUsers(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	out := make([]map[string]interface{}, 0, len(list))
	for _, su := range list {
		if su.Uuid == "" {
			continue
		}
		out = append(out, map[string]interface{}{
			"id":          su.Uuid,
			"name":        su.Name,
			"description": su.Description,
		})
	}

	d.SetId(c.WorkspaceUUID.String())
	if err := d.Set("users", out); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set users: %w", err))
	}
	return nil
}
