package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func dataSourceExternalIPs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExternalIPsRead,
		Description: "Fetches a list of External IPs from the configured workspace.",
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Optional filter to apply to the list of External IPs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unbound_only": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Set to true to return only IPs that are not currently bound to a resource.",
						},
					},
				},
			},
			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of external IPs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_bound": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceExternalIPsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	apiResponseItems, err := c.ListExternalIPs()
	fmt.Println(apiResponseItems)
	if err != nil {
		return diag.FromErr(err)
	}

	// --- Filtering Logic (remains the same) ---
	unboundOnly := false
	if v, ok := d.GetOk("filter"); ok {
		filterList := v.([]interface{})
		if len(filterList) > 0 {
			filterMap := filterList[0].(map[string]interface{})
			unboundOnly = filterMap["unbound_only"].(bool)
		}
	}

	flattenedIPs := make([]interface{}, 0)
	for _, ip := range apiResponseItems {
		isBound := ip.Spec.BoundTo != nil

		if unboundOnly && isBound {
			continue
		}

		ipMap := map[string]interface{}{
			"name":       ip.Metadata.Name,
			"ip_address": ip.Spec.IP,
			"uid":        ip.Metadata.UID,
			"is_bound":   isBound,
		}
		flattenedIPs = append(flattenedIPs, ipMap)
	}

	if err := d.Set("items", flattenedIPs); err != nil {
		return diag.FromErr(err)
	}

	// Set a unique ID for the data source state.
	d.SetId(time.Now().UTC().String())

	return nil
}
