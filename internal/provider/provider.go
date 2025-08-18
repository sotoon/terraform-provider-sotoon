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
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOTOON_USER_ID", ""),
				Description: "The Sotoon UserID.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sotoon_iam_user":            resourceUser(),
			"sotoon_iam_group":           resourceGroup(),
			"sotoon_iam_user_token":      resourceUserToken(),
			"sotoon_iam_user_public_key": resourceUserPublicKey(),
			"sotoon_iam_user_group_membership": resourceUserGroupMembership(),


			"sotoon_iam_group_role":              resourceGroupRole(),
			"sotoon_iam_service_user_group":      resourceGroupServiceUser(),
			"sotoon_iam_service_user":            resourceServiceUser(),
			"sotoon_iam_service_user_token":      resourceServiceUserToken(),
			"sotoon_iam_service_user_public_key": resourceServiceUserPublicKey(),
			"sotoon_iam_service_user_role":       resourceServiceUserRole(),
			"sotoon_iam_user_role":               resourceUserRole(),
			"sotoon_iam_role":                    resourceRole(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sotoon_iam_users":            dataSourceUsers(),
			"sotoon_iam_groups":           dataSourceGroups(),
			"sotoon_iam_user_tokens":      dataSourceUserTokens(),
			"sotoon_iam_user_public_keys": dataSourceUserPublicKeys(),

			"sotoon_iam_group_users":              dataSourceGroupUsers(),
			"sotoon_iam_group_details":            dataSourceGroupDetails(),
			"sotoon_iam_group_roles":              dataSourceGroupRoles(),
			"sotoon_iam_service_user_groups":      dataSourceGroupUserServices(),
			"sotoon_iam_service_user_details":     dataSourceServiceUserDetails(),
			"sotoon_iam_service_user_public_keys": dataSourceServiceUserPublicKeys(),	
			"sotoon_iam_service_user_tokens":      dataSourceServiceUserTokens(),
			"sotoon_iam_service_users":            dataSourceServiceUsers(),
			"sotoon_iam_roles":                    dataSourceRoles(),
			"sotoon_iam_rules":                    dataSourceRules(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("api_token").(string)
	workspaceID := d.Get("workspace_id").(string)
	host := d.Get("api_host").(string)
	userID := d.Get("user_id").(string)

	var diags diag.Diagnostics

	if token == "" || workspaceID == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Sotoon client",
			Detail:   "API token and Workspace ID must be configured for the provider.",
		})
		return nil, diags
	}

	c, err := client.NewClient(host, token, workspaceID, userID)
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
