package iam

import (
	"context"
	"strings"

	bepa "git.cafebazaar.ir/infrastructure/bepa-client/pkg/client"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/models"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/utils"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	uuid "github.com/satori/go.uuid"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &userGroupBindingResource{}
	_ resource.ResourceWithConfigure   = &userGroupBindingResource{}
	_ resource.ResourceWithImportState = &userGroupBindingResource{}
)

func NewUserGroupBindingResource() resource.Resource {
	return &userGroupBindingResource{}
}

type userGroupBindingResource struct {
	iamClient bepa.Client
}

func (r *userGroupBindingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.iamClient = req.ProviderData.(*utils.SotoonConfigDataHolder).IAMClient
}

func (r *userGroupBindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_user_group_binding"
}

func (r *userGroupBindingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Membership relation between a user and a group. This object defiens the membership of a user in a group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique id of group-user binding",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "ID of binded used to the group.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID of the group.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"group_id": schema.StringAttribute{
				Description: "The ID of the group to which the user will be binded.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
		},
	}
}

func (r *userGroupBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.UserGroup
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(plan.WorkspaceID.ValueString())
	groupID, _ := uuid.FromString(plan.GroupID.ValueString())
	userID, _ := uuid.FromString(plan.UserID.ValueString())

	group, err := r.iamClient.GetGroup(&workspaceID, &groupID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_GROUP_READ, errorMessage)
		return
	}

	err = r.iamClient.BindGroup(group.Name, &workspaceID, &groupID, &userID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_USER_GROUP_CREATE, errorMessage)
		return
	}
	plan.UUID = types.StringValue(uuid.NewV4().String())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userGroupBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.UserGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(state.WorkspaceID.ValueString())
	groupID, _ := uuid.FromString(state.GroupID.ValueString())
	userID, _ := uuid.FromString(state.UserID.ValueString())

	_, err := r.iamClient.GetGroupUser(&workspaceID, &groupID, &userID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_USER_GROUP_READ, errorMessage)
		return
	}

	if state.UUID.IsNull() {
		state.UUID = types.StringValue(uuid.NewV4().String())
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userGroupBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.UserGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(state.WorkspaceID.ValueString())
	groupID, _ := uuid.FromString(state.GroupID.ValueString())
	userID, _ := uuid.FromString(state.UserID.ValueString())

	err := r.iamClient.UnbindUserFromGroup(&workspaceID, &groupID, &userID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_USER_GROUP_DELETE, errorMessage)
		return
	}
}

func (r *userGroupBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state models.UserGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(state.WorkspaceID.ValueString())
	groupID, _ := uuid.FromString(state.GroupID.ValueString())
	userID, _ := uuid.FromString(state.UserID.ValueString())

	err := r.iamClient.UnbindUserFromGroup(&workspaceID, &groupID, &userID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_USER_GROUP_DELETE_TO_UPDATE, errorMessage)
		return
	}

	var plan models.UserGroup
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ = uuid.FromString(plan.WorkspaceID.ValueString())
	userID, _ = uuid.FromString(plan.UserID.ValueString())
	groupID, _ = uuid.FromString(plan.GroupID.ValueString())

	group, err := r.iamClient.GetGroup(&workspaceID, &groupID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_USER_GROUP_READ, errorMessage)
		return
	}

	err = r.iamClient.BindGroup(group.Name, &workspaceID, &groupID, &userID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_USER_GROUP_CREATE, errorMessage)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userGroupBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := strings.Split(req.ID, ":")
	if len(importID) < 3 {
		resp.Diagnostics.AddError(
			"Error Importing Group-User Attachment",
			"The given ID must be in the form of {user id}:{group id}:{workspace id}",
		)
		return
	}
	userUUID := importID[0]
	groupUUID := importID[1]
	workspaceUUID := importID[2]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), userUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), groupUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_id"), workspaceUUID)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
