package resource

import (
	"context"
	"fmt"
	"terraform-provider-sotoon/internal/client"
	"terraform-provider-sotoon/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &Subnet{}
var _ resource.ResourceWithImportState = &Subnet{}

func NewSubnet() resource.Resource {
	return &Subnet{}
}

// Subnet defines the resource implementation.
type Subnet struct {
	client *client.Client
}

func (s *Subnet) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet"
}

func (s *Subnet) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"cidr": schema.StringAttribute{
				Required: true,
			},
			"gateway_ip": schema.StringAttribute{
				Computed: true,
			},
			"routes": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"to": schema.StringAttribute{
							Required: true,
						},
						"external_ip": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

func (s *Subnet) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	s.client = c
}

func (s *Subnet) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.Subnet
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	err := s.client.Subnet.Create(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("creating subnet failed", err.Error())
		return
	}

	diags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *Subnet) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.Subnet
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subnet, err := s.client.Subnet.Get(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("reading subnet failed", err.Error())
		return
	}

	state = *subnet
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *Subnet) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.Subnet
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.Subnet.Update(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("updating subnet failed", err.Error())
		return
	}

	diags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (s *Subnet) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.Subnet
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.Subnet.Delete(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("reading subnet failed", err.Error())
		return
	}
}

func (s *Subnet) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
