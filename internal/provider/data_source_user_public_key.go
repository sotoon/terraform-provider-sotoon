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

func dataSourceUserPublicKeys() *schema.Resource {
	return &schema.Resource{
		Description: "Lists all SSH/public-keys for the current IAM user.",
		ReadContext: dataSourceUserPublicKeysRead,
		Schema: map[string]*schema.Schema{
			"public_keys": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of SSH/public-key records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UUID of the public-key.",
						},
						"title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The key title.",
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The SSH public key material.",
						},
					},
				},
			},
		},
	}
}

func dataSourceUserPublicKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Debug(ctx, "Listing all public keys")

	var list []*types.PublicKey

	all, err := c.GetAllMyUserPublicKeyList()
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not exists") {
			list = []*types.PublicKey{} // no keys → empty
		} else if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return diag.Errorf("forbidden: cannot list public-keys: %s", err)
		} else {
			return diag.Errorf("error listing public-keys: %s", err)
		}
	} else {
		list = all
	}

	out := make([]map[string]interface{}, 0, len(list))
	for _, pk := range list {
		out = append(out, map[string]interface{}{
			"id":    pk.UUID,
			"title": pk.Title,
			"key":   pk.Key,
		})
	}

	if err := d.Set("public_keys", out); err != nil {
		return diag.Errorf("failed to set public_keys: %s", err)
	}
	d.SetId("current")
	return nil
}
