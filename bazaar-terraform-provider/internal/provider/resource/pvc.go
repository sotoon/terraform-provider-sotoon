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
var _ resource.Resource = &PVC{}
var _ resource.ResourceWithImportState = &PVC{}

func NewPvc() resource.Resource {
	return &PVC{}
}

// PVC defines the resource implementation.
type PVC struct {
	client *client.Client
}

func (p *PVC) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pvc"
}

func (p *PVC) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"size": schema.StringAttribute{
				Required: true,
			},
			"tier": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (p *PVC) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	p.client = c
}

func (p *PVC) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.PVC
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	err := p.client.PVC.Create(ctx, &plan)
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

func (p *PVC) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.PVC
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pvc, err := p.client.PVC.Get(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("reading pvc failed", err.Error())
		return
	}

	state = *pvc
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *PVC) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.PVC
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := p.client.PVC.Update(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("updating pvc failed", err.Error())
		return
	}

	diags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (p *PVC) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.PVC
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := p.client.PVC.Delete(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("reading pvc failed", err.Error())
		return
	}
}

func (p *PVC) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
