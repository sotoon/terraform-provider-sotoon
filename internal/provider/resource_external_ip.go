package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func resourceExternalIP() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages an External IP address on the Sotoon platform.",
		CreateContext: resourceExternalIPCreate,
		ReadContext:   resourceExternalIPRead,
		DeleteContext: resourceExternalIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the External IP.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the external IP.",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The workspace ID.",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The allocated public IP address.",
			},
		},
	}
}

func resourceExternalIPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	name := d.Get("name").(string)

	ip, err := c.CreateExternalIP(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ip.Metadata.Name)
	return resourceExternalIPRead(ctx, d, m)
}

func resourceExternalIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	name := d.Id()

	ip, err := c.GetExternalIP(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	if ip == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", ip.Metadata.Name)
	d.Set("workspace_id", ip.Metadata.Workspace)
	d.Set("address", ip.Spec.IP)

	return nil
}

func resourceExternalIPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	name := d.Id()

	err := c.DeleteExternalIP(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
