package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

// Provider defines the provider schema and configuration
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SOTOON_API_TOKEN", nil),
				Description: "The API token for Sotoon cloud.",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOTOON_WORKSPACE_ID", nil),
				Description: "The ID of the workspace.",
			},
			"api_host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOTOON_API_HOST", "https://bepa.sotoon.ir"),
				Description: "The Sotoon API host.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sotoon_instance":    resourceInstance(),
			"sotoon_external_ip": resourceExternalIP(),
			"sotoon_iam_user": resourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sotoon_instance":     dataSourceInstances(),
			"sotoon_external_ips": dataSourceExternalIPs(),
			"sotoon_iam_users": dataSourceUsers(),
			"sotoon_iam_groups": dataSourceGroups(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("api_token").(string)
	workspaceID := d.Get("workspace_id").(string)
	host := d.Get("api_host").(string)

	var diags diag.Diagnostics

	if token == "" || workspaceID == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Sotoon client",
			Detail:   "API token and Workspace ID must be configured for the provider.",
		})
		return nil, diags
	}

	c, err := client.NewClient(host, token, workspaceID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Sotoon client",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return c, diags
}
