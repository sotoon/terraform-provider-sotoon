package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func dataSourceUserTokens() *schema.Resource {
	return &schema.Resource{
		Description: "Lists all tokens belonging to the current IAM user.",
		ReadContext: dataSourceUserTokensRead,
		Schema: map[string]*schema.Schema{
			"tokens": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of user-token objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UUID of the token.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The token value (empty for existing tokens).",
						},
					},
				},
			},
		},
	}
}

func dataSourceUserTokensRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Debug(ctx, "Listing all user tokens")

	list, err := c.GetAllMyUserTokenList(ctx)
	if err != nil {
		return diag.Errorf("error listing tokens: %s", err)
	}

	out := make([]map[string]interface{}, 0, len(list))
	for _, tok := range list {
		out = append(out, map[string]interface{}{
			"id":    tok.Uuid,
			"value": "",
		})
	}

	if err := d.Set("tokens", out); err != nil {
		return diag.Errorf("failed to set tokens: %s", err)
	}
	d.SetId("current")
	return nil
}
