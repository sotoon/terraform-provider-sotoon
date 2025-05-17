package ds

import (
	"context"
	"fmt"
	"terraform-provider-sotoon/internal/client"
	"terraform-provider-sotoon/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ExternalIP{}

func NewExternalIP() datasource.DataSource {
	return &ExternalIP{}
}

type ExternalIP struct {
	client *client.Client
}

func (e *ExternalIP) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_externalip"
}

func (e *ExternalIP) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
			},
			"external_ips": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"ip": schema.StringAttribute{
							Computed: true,
						},
						"gateway_ip": schema.StringAttribute{
							Computed: true,
						},
						"reserved": schema.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (e *ExternalIP) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.IClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	e.client = c
}

func (e *ExternalIP) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var externalIps []model.ExternalIP
	externalIps = []model.ExternalIP{}
	var state struct {
		Name        types.String       `tfsdk:"name"`
		ExternalIps []model.ExternalIP `tfsdk:"external_ips"`
	}
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to read data source configuration")
		return
	}

	if !state.Name.IsNull() {
		externalIp, err := e.client.ExternalIP.List(ctx)
		if err != nil {
			tflog.Error(ctx, "Error getting subnet", map[string]interface{}{"err": err})
		}
		externalIps = append(externalIps, externalIp...)
	} else {
		externalIps, _ = e.client.ExternalIP.List(ctx)
	}

	state.ExternalIps = externalIps
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
