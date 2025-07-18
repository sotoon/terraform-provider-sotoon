package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func dataSourceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstancesRead,
		Description: "Fetches a list of all compute instances in the workspace.",
		Schema: map[string]*schema.Schema{
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of all instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the instance.",
						},
						"uid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier (UID) of the instance.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance type.",
						},
						"image": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The OS image used by the instance.",
						},
						"powered_on": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the instance is powered on.",
						},
						"iam_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether IAM is enabled for the instance.",
						},
						"interfaces": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Network interfaces attached to the instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"link": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ready": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The readiness status of the instance.",
						},
						"running": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The running status of the instance.",
						},
					},
				},
			},
		},
	}
}

func dataSourceInstancesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	apiInstances, err := c.ListInstances(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	instances := make([]map[string]interface{}, 0, len(apiInstances))

	for _, apiInstance := range apiInstances {
		interfaces := make([]map[string]interface{}, len(apiInstance.Spec.Interfaces))
		for i, iface := range apiInstance.Spec.Interfaces {
			interfaces[i] = map[string]interface{}{
				"name": iface.Name,
				"link": iface.Link,
			}
		}

		instanceMap := map[string]interface{}{
			"name":        apiInstance.Metadata.Name,
			"uid":         apiInstance.Metadata.UID,
			"type":        apiInstance.Spec.Type,
			"image":       apiInstance.Spec.ImageSource.Image,
			"powered_on":  apiInstance.Spec.PoweredOn,
			"iam_enabled": apiInstance.Spec.IAMEnabled,
			"interfaces":  interfaces,
		}
		instances = append(instances, instanceMap)
	}

	if err := d.Set("instances", instances); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().UTC().String())

	return nil
}
