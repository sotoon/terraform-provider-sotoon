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
var _ resource.Resource = &Compute{}
var _ resource.ResourceWithImportState = &Compute{}

func NewCompute() resource.Resource {
	return &Compute{}
}

// Compute defines the resource implementation.
type Compute struct {
	client *client.Client
}

func (c *Compute) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compute"
}

func (c *Compute) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"iam_enabled": schema.BoolAttribute{
				Required: true,
			},
			"image": schema.StringAttribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"size": schema.StringAttribute{
				Required: true,
			},
			"volumes": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"local_disk": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"disk_size_gb": schema.Int64Attribute{
									Required: true,
								},
								"name": schema.StringAttribute{
									Required: true,
								},
							},
							Optional: true,
						},
						"pvc": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
							},
							Optional: true,
						},
						"remote_disk": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"size": schema.Int64Attribute{
									Required: true,
								},
								"tier": schema.StringAttribute{
									Required: true,
								},
							},
							Optional: true,
						},
					},
				},
			},
			"powered_on": schema.BoolAttribute{
				Computed: true,
			},
			"subnet": schema.StringAttribute{
				Required: true,
			},
			"external_ip": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (c *Compute) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	c.client = r
}

func (c *Compute) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.Compute
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	err := c.client.Compute.Create(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("creating compute failed", err.Error())
		return
	}

	diags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *Compute) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.Compute
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	compute, err := c.client.Compute.Get(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("reading compute failed", err.Error())
		return
	}

	state = *compute
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *Compute) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.Compute
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := c.client.Compute.Update(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("updating compute failed", err.Error())
		return
	}

	diags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (c *Compute) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.Compute
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := c.client.Compute.Delete(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("reading compute failed", err.Error())
		return
	}
}

func (c *Compute) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
