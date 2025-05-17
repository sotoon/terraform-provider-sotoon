package iam

import (
	"context"
	"strings"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/client"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/models"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/utils"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	uuid "github.com/satori/go.uuid"
)

var (
	_ resource.Resource                = &groupResource{}
	_ resource.ResourceWithConfigure   = &groupResource{}
	_ resource.ResourceWithImportState = &groupResource{}
)

func NewGroupResource() resource.Resource {
	return &groupResource{}
}

type groupResource struct {
	iamClient client.Client
}

func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_group"
}

func (r *groupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.iamClient = req.ProviderData.(*utils.SotoonConfigDataHolder).IAMClient
}

func (r *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a group instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the group.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "ID of the workspace which group is defined in that.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
		},
	}
}

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.Group
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupName := plan.Name.ValueString()

	workspaceUUID, _ := uuid.FromString(plan.WorkspaceUUID.ValueString())

	createdGroup, err := r.iamClient.CreateGroup(groupName, &workspaceUUID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_GROUP_CREATE, errorMessage)
		return
	}
	plan.UUID = types.StringValue(createdGroup.UUID.String())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.Group
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceUUID, _ := uuid.FromString(state.WorkspaceUUID.ValueString())

	groupName := state.Name.ValueString()

	workspace, err := r.iamClient.GetWorkspace(&workspaceUUID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_GROUP_WORKSPACE_READ, errorMessage)
		return
	}

	group, err := r.iamClient.GetGroupByName(workspace.Name, groupName)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_GROUP_READ, errorMessage)
		return
	}

	state = *models.NewGroupModelFromGroupObject(group)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Error
	resp.Diagnostics.AddError(
		"Group Can't Be Changed!",
		"You can delete the current group and create the new one",
	)
}

func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.Group
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceUUID, _ := uuid.FromString(state.WorkspaceUUID.ValueString())

	groupUUID, _ := uuid.FromString(state.UUID.ValueString())

	err := r.iamClient.DeleteGroup(&workspaceUUID, &groupUUID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError("Error Deleting Group", errorMessage)
		return
	}
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := strings.Split(req.ID, ":")
	if len(importID) < 2 {
		resp.Diagnostics.AddError(
			"Error Importing Group",
			"The given ID must be in the form of {group name}:{workspace id}",
		)
		return
	}
	groupName := importID[0]
	workspaceUUID := importID[1]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), groupName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_id"), workspaceUUID)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
