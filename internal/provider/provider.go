package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"

	// datasources
	iamdatasources "github.com/sotoon/terraform-provider-sotoon/internal/datasources/iam"

	// resources
	iamresources "github.com/sotoon/terraform-provider-sotoon/internal/resources/iam"
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
			"sotoon_iam_user":                    iamresources.ResourceUser(),
			"sotoon_iam_group":                   iamresources.ResourceGroup(),
			"sotoon_iam_user_token":              iamresources.ResourceUserToken(),
			"sotoon_iam_user_public_key":         iamresources.ResourceUserPublicKey(),
			"sotoon_iam_user_group_membership":   iamresources.ResourceUserGroupMembership(),
			"sotoon_iam_group_role":              iamresources.ResourceGroupRole(),
			"sotoon_iam_service_user_group":      iamresources.ResourceGroupServiceUser(),
			"sotoon_iam_service_user":            iamresources.ResourceServiceUser(),
			"sotoon_iam_service_user_token":      iamresources.ResourceServiceUserToken(),
			"sotoon_iam_service_user_public_key": iamresources.ResourceServiceUserPublicKey(),
			"sotoon_iam_service_user_role":       iamresources.ResourceServiceUserRole(),
			"sotoon_iam_user_role":               iamresources.ResourceUserRole(),
			"sotoon_iam_role":                    iamresources.ResourceRole(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sotoon_iam_users":                    iamdatasources.DataSourceUsers(),
			"sotoon_iam_groups":                   iamdatasources.DataSourceGroups(),
			"sotoon_iam_user_tokens":              iamdatasources.DataSourceUserTokens(),
			"sotoon_iam_user_public_keys":         iamdatasources.DataSourceUserPublicKeys(),
			"sotoon_iam_group_users":              iamdatasources.DataSourceGroupUsers(),
			"sotoon_iam_group_details":            iamdatasources.DataSourceGroupDetails(),
			"sotoon_iam_group_roles":              iamdatasources.DataSourceGroupRoles(),
			"sotoon_iam_service_user_groups":      iamdatasources.DataSourceGroupUserServices(),
			"sotoon_iam_service_user_details":     iamdatasources.DataSourceServiceUserDetails(),
			"sotoon_iam_service_user_public_keys": iamdatasources.DataSourceServiceUserPublicKeys(),
			"sotoon_iam_service_user_tokens":      iamdatasources.DataSourceServiceUserTokens(),
			"sotoon_iam_service_users":            iamdatasources.DataSourceServiceUsers(),
			"sotoon_iam_roles":                    iamdatasources.DataSourceRoles(),
			"sotoon_iam_rules":                    iamdatasources.DataSourceRules(),
			"sotoon_iam_user_roles":               iamdatasources.DataSourceUserRoles(),
			"sotoon_iam_service_user_roles":       iamdatasources.DataSourceServiceUserRoles(),
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
