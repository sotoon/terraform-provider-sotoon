package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func dataSourceImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImagesRead,
		Schema: map[string]*schema.Schema{
			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceImagesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics

	images, err := c.ListImages()
	if err != nil {
		return diag.FromErr(err)
	}

	imageList := make([]map[string]interface{}, 0, len(images))
	for _, image := range images {
		img := map[string]interface{}{
			"name":        image.Name,
			"version":     image.Version,
			"os_type":     image.OsType,
			"description": image.Description,
		}
		imageList = append(imageList, img)
	}

	if err := d.Set("images", imageList); err != nil {
		return diag.FromErr(err)
	}

	// Set a unique ID for the data source
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
