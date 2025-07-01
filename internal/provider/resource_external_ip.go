package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func resourceExternalIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceExternalIPCreate,
		ReadContext:   resourceExternalIPRead,
		DeleteContext: resourceExternalIPDelete,
		Schema: map[string]*schema.Schema{
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

	ip, err := c.CreateExternalIP(name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ip.Metadata.Name) // Using name as ID, assuming it's unique
	return resourceExternalIPRead(ctx, d, m)
}

func resourceExternalIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	name := d.Id()
	var diags diag.Diagnostics

	ip, err := c.GetExternalIP(name)
	if err != nil {
		return diag.FromErr(err)
	}

	if ip == nil {
		// Resource not found
		d.SetId("")
		return diags
	}

	d.Set("name", ip.Metadata.Name)
	d.Set("workspace_id", ip.Metadata.Workspace)

	return diags
}

func resourceExternalIPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	name := d.Id()
	var diags diag.Diagnostics

	err := c.DeleteExternalIP(name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("") // Mark as deleted
	return diags
}
