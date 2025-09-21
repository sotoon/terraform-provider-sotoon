package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func dataSourceServiceUserTokens() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServiceUserTokensRead,
		Schema: map[string]*schema.Schema{
			"service_user_id": {Type: schema.TypeString, Required: true},
			"id":              {Type: schema.TypeString, Computed: true},
			"tokens": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceServiceUserTokensRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	idStr := d.Get("service_user_id").(string)
	suID, err := uuid.FromString(idStr)
	if err != nil {
		return diag.FromErr(err)
	}

	list, err := c.GetWorkspaceServiceUserTokenList(&suID, c.WorkspaceUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, 0)
	if list != nil {
		for _, t := range *list {
			if t.UUID == nil {
				continue
			}
			result = append(result, map[string]interface{}{"id": t.UUID.String()})
		}
	}
	d.SetId(suID.String())
	if err := d.Set("tokens", result); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set tokens: %w", err))
	}
	return nil
}
