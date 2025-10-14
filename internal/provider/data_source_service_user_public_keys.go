package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func dataSourceServiceUserPublicKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServiceUserPublicKeysRead,
		Schema: map[string]*schema.Schema{
			"service_user_id": {Type: schema.TypeString, Required: true},
			"id":              {Type: schema.TypeString, Computed: true},
			"public_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":    {Type: schema.TypeString, Computed: true},
						"title": {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceServiceUserPublicKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	idStr := d.Get("service_user_id").(string)
	suID, err := uuid.FromString(idStr)
	if err != nil {
		return diag.FromErr(err)
	}

	list, err := c.GetWorkspaceServiceUserPublicKeyList(ctx, *c.WorkspaceUUID, suID)
	if err != nil {
		return diag.FromErr(err)
	}

	out := make([]map[string]interface{}, 0, len(list))
	for _, pk := range list {
		if pk.Uuid == "" {
			continue
		}
		out = append(out, map[string]interface{}{
			"id":    pk.Uuid,
			"title": pk.Title,
		})
	}

	d.SetId(suID.String())
	if err := d.Set("public_keys", out); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set public_keys: %w", err))
	}
	return nil
}
