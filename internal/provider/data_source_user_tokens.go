package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sotoon/iam-client/pkg/types"
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

	var list []types.UserToken

	ptr, err := c.GetAllMyUserTokenList()
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not exists") {
			list = []types.UserToken{} // no tokens → empty
		} else if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return diag.Errorf("forbidden: cannot list tokens: %s", err)
		} else {
			return diag.Errorf("error listing tokens: %s", err)
		}
	} else {
		list = *ptr
	}

	out := make([]map[string]interface{}, 0, len(list))
	for _, tok := range list {
		out = append(out, map[string]interface{}{
			"id":    tok.UUID,
			"value": "",
		})
	}

	if err := d.Set("tokens", out); err != nil {
		return diag.Errorf("failed to set tokens: %s", err)
	}
	d.SetId("current")
	return nil
}
