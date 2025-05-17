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
var _ datasource.DataSource = &PVC{}

func NewPVC() datasource.DataSource {
	return &PVC{}
}

type PVC struct {
	client *client.Client
}

func (p *PVC) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pvc"
}

func (p *PVC) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"pvcs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"size": schema.StringAttribute{
							Computed: true,
						},
						"tier": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (p *PVC) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	p.client = c
}

func (p *PVC) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	pvcs, _ := p.client.PVC.List(ctx)
	diags := resp.State.Set(ctx, &model.PVCS{PVCS: pvcs})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
