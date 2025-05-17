package ds

import (
	"context"
	"fmt"
	"terraform-provider-sotoon/internal/client"
	"terraform-provider-sotoon/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &Subnet{}

func NewSubnet() datasource.DataSource {
	return &Subnet{}
}

type Subnet struct {
	client *client.Client
}

func (s *Subnet) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet"
}

func (s *Subnet) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"subnets": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"cidr": schema.StringAttribute{
							Computed: true,
						},
						"gateway_ip": schema.StringAttribute{
							Computed: true,
						},
						"routes": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"to": schema.StringAttribute{
										Computed: true,
									},
									"via": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"external_ip": schema.StringAttribute{
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (s *Subnet) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	s.client = c
}

func (s *Subnet) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	subnets, _ := s.client.Subnet.List(ctx)
	diags := resp.State.Set(ctx, &model.Subnets{Subnets: subnets})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
