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
	_ resource.Resource                = &serviceUserGroupBindingResource{}
	_ resource.ResourceWithConfigure   = &serviceUserGroupBindingResource{}
	_ resource.ResourceWithImportState = &serviceUserGroupBindingResource{}
)

func NewServiceUserGroupBindingResource() resource.Resource {
	return &serviceUserGroupBindingResource{}
}

type serviceUserGroupBindingResource struct {
	iamClient bepa.Client
}

func (r *serviceUserGroupBindingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.iamClient = req.ProviderData.(*utils.SotoonConfigDataHolder).IAMClient
}

func (r *serviceUserGroupBindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_service_user_group_binding"
}

func (r *serviceUserGroupBindingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Membership relation between a service-user and a group. This object defiens the membership of a service-user in a group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique id of the binding",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "Service-user ID.",
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
				Description: "The ID of the group to which the service-user will be binded.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
		},
	}
}

func (r *serviceUserGroupBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.UserGroup
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(plan.WorkspaceID.ValueString())
	groupID, _ := uuid.FromString(plan.GroupID.ValueString())
	userID, _ := uuid.FromString(plan.UserID.ValueString())

	err := r.iamClient.BindServiceUserToGroup(&workspaceID, &groupID, &userID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_SERVICE_USER_GROUP_CREATE, errorMessage)
		return
	}
	plan.UUID = types.StringValue(uuid.NewV4().String())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *serviceUserGroupBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.UserGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(state.WorkspaceID.ValueString())
	groupID, _ := uuid.FromString(state.GroupID.ValueString())
	userID, _ := uuid.FromString(state.UserID.ValueString())

	_, err := r.iamClient.GetGroupServiceUser(&workspaceID, &groupID, &userID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_SERVICE_USER_GROUP_READ, errorMessage)
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

func (r *serviceUserGroupBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.UserGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(state.WorkspaceID.ValueString())
	groupID, _ := uuid.FromString(state.GroupID.ValueString())
	userID, _ := uuid.FromString(state.UserID.ValueString())

	err := r.iamClient.UnbindServiceUserFromGroup(&workspaceID, &groupID, &userID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_SERVICE_USER_GROUP_DELETE, errorMessage)
		return
	}
}

func (r *serviceUserGroupBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state models.UserGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(state.WorkspaceID.ValueString())
	groupID, _ := uuid.FromString(state.GroupID.ValueString())
	userID, _ := uuid.FromString(state.UserID.ValueString())

	err := r.iamClient.UnbindServiceUserFromGroup(&workspaceID, &groupID, &userID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_SERVICE_USER_GROUP_DELETE, errorMessage)
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

	err = r.iamClient.BindServiceUserToGroup(&workspaceID, &groupID, &userID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_SERVICE_USER_GROUP_CREATE, errorMessage)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *serviceUserGroupBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := strings.Split(req.ID, ":")
	if len(importID) < 3 {
		resp.Diagnostics.AddError(
			"Error Importing Group-ServiceUser Attachment",
			"The given ID must be in the form of {service-user id}:{group id}:{workspace id}",
		)
		return
	}
	serviceUserUUID := importID[0]
	groupUUID := importID[1]
	workspaceUUID := importID[2]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), serviceUserUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), groupUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_id"), workspaceUUID)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
