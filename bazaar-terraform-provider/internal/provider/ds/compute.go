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
var _ datasource.DataSource = &Compute{}

func NewCompute() datasource.DataSource {
	return &Compute{}
}

type Compute struct {
	client *client.Client
}

func (c *Compute) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compute"
}

func (c *Compute) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"Computes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"iam_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"image": schema.StringAttribute{
							Computed: true,
						},
						"username": schema.StringAttribute{
							Computed: true,
						},
						"size": schema.StringAttribute{
							Computed: true,
						},
						"volumes": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"local_disk": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"disk_size_gb": schema.Int64Attribute{
												Computed: true,
											},
											"name": schema.StringAttribute{
												Computed: true,
											},
										},
									},
									"pvc": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Computed: true,
											},
										},
									},
								},
							},
						},
						"powered_on": schema.BoolAttribute{
							Computed: true,
						},
						"subnet": schema.StringAttribute{
							Computed: true,
						},
						"external_ip": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (c *Compute) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.IClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	c.client = r
}

func (c *Compute) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	computes, _ := c.client.Compute.List(ctx)
	diags := resp.State.Set(ctx, &model.Computes{Computes: computes})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
