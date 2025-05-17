package provider

import (
	"context"
	"os"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/client"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/common"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &sotoonProvider{}
)

// TODO: We shouldn't use this user uuid to create new bepa
// client but because of using GetGroupByName function and
// it's need we use this dummy useless user uuid. sorry :(
const dummyUserUUID = "00000000-0000-0000-0000-000000000000"

func NewProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &sotoonProvider{
			version: version,
		}
	}
}

type sotoonProvider struct {
	version string
	mocks   *MockKeeper
}

type sotoonProviderModel struct {
	Token types.String `tfsdk:"token"`
}

func (p *sotoonProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sotoon"
}

func (p *sotoonProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Sotoon.",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Description: "Access token of Sotoon API. May also be provided via `SOTOON_TOKEN` environment variable. `token` is your account's access token and you can get it from your panel.",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *sotoonProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Sotoon client")

	var config sotoonProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Sotoon API Access Token",
			"The provider cannot create the Sotoon API client as there is an unknown configuration value for the Sotoon API access token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SOTOON_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	token := os.Getenv("SOTOON_TOKEN")
	var workspace_uuid string

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Sotoon API Access Token",
			"The provider cannot create the Sotoon API client as there is a missing or empty value for the Sotoon API access token. "+
				"Set the access token value in the configuration or use the SOTOON_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "sotoon_token", token)
	ctx = tflog.SetField(ctx, "sotoon_workspace_uuid", workspace_uuid)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "sotoon_token")

	configFile := common.Config(ctx)

	var configDataHolder *utils.SotoonConfigDataHolder

	if p.mocks != nil {
		configDataHolder = &utils.SotoonConfigDataHolder{
			IAMClient: p.mocks.IAMClient,
		}
	} else {
		// Please see the comment of dummyUserUUID's definition to see why we use this useless uuid!
		iamClient, err := client.NewReliableClient(token, configFile.BepaUrls, "", dummyUserUUID, configFile.BepaTimeout)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Create Sotoon IAM API Client",
				"An unexpected error occurred when creating the Sotoon API client. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"Sotoon Client Error: "+err.Error(),
			)
			return
		}

		configDataHolder = &utils.SotoonConfigDataHolder{
			IAMClient: iamClient,
		}
	}

	resp.DataSourceData = configDataHolder
	resp.ResourceData = configDataHolder

	tflog.Info(ctx, "Configured Sotoon client", map[string]any{"success": true})
}

func (p *sotoonProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return GetDataSources()
}

func (p *sotoonProvider) Resources(_ context.Context) []func() resource.Resource {
	return GetResources()
}
