package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	sdkResource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
	"terraform-provider-sotoon/internal/client"
	"terraform-provider-sotoon/internal/provider/ds"
	"terraform-provider-sotoon/internal/provider/resource"
)

var _ provider.Provider = &SotoonProvider{}
var _ provider.ProviderWithFunctions = &SotoonProvider{}

type SotoonProvider struct {
	version string
}

type SotoonProviderModel struct {
	Host          types.String `tfsdk:"host"`
	WorkspaceName types.String `tfsdk:"workspace_name"`
	WorkspaceID   types.String `tfsdk:"workspace_id"`
	Zone          types.String `tfsdk:"zone"` // todo: move it to sdkResource and datasource
	Token         types.String `tfsdk:"token"`
}

func (sp *SotoonProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sotoon"
	resp.Version = sp.version
}

func (sp *SotoonProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required: true,
			},
			"zone": schema.StringAttribute{
				Required: true,
			},
			"workspace_name": schema.StringAttribute{
				Required: true,
			},
			"workspace_id": schema.StringAttribute{
				Required: true,
			},
			"token": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (sp *SotoonProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config SotoonProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	host := config.Host.ValueString()
	workspaceName := config.WorkspaceName.ValueString()
	workspaceID := config.WorkspaceID.ValueString()
	zone := config.Zone.ValueString()
	token := os.Getenv("PROVIDER_TOKEN")

	ctx = tflog.SetField(ctx, "host", host)
	ctx = tflog.SetField(ctx, "token", token)
	//ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "token")
	tflog.Info(ctx, "Creating client")

	if config.Host.IsUnknown() && config.Host.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Sotoon API Host",
			"The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_HOST environment variable.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Sotoon API token",
			"The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_USERNAME environment variable.",
		)
	}

	if config.WorkspaceID.IsUnknown() && config.WorkspaceID.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("workspace_id"),
			"Unknown Sotoon workspace ID",
			"The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_USERNAME environment variable.",
		)
	}

	if config.WorkspaceName.IsUnknown() && config.WorkspaceName.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("wokspace_name"),
			"Unknown Sotoon workspace name",
			"The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_USERNAME environment variable.",
		)
	}

	if config.Zone.IsUnknown() && config.Zone.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("zone"),
			"Unknown Sotoon zone ",
			"The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_USERNAME environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	c := client.NewClient(
		host,
		workspaceName,
		workspaceID,
		zone,
		token,
	)
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (sp *SotoonProvider) Resources(ctx context.Context) []func() sdkResource.Resource {
	return []func() sdkResource.Resource{
		NewExampleResource,
		resource.NewSubnet,
		resource.NewPvc,
	}
}

func (sp *SotoonProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		ds.NewSubnet,
		ds.NewExternalIP,
		ds.NewPVC,
	}
}

func (sp *SotoonProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SotoonProvider{
			version: version,
		}
	}
}
