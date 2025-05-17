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
	_ resource.Resource                = &roleUserBindingResource{}
	_ resource.ResourceWithConfigure   = &roleUserBindingResource{}
	_ resource.ResourceWithImportState = &roleUserBindingResource{}
)

func NewRoleUserBindingResource() resource.Resource {
	return &roleUserBindingResource{}
}

type roleUserBindingResource struct {
	iamClient bepa.Client
}

func (r *roleUserBindingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.iamClient = req.ProviderData.(*utils.SotoonConfigDataHolder).IAMClient
}

func (r *roleUserBindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_role_user_binding"
}

func (r *roleUserBindingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Appliance relation between role and user. Existance of an instance from this resource will applies an specified role to the selected user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the binding",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "ID of the user which the role is going to be binded to that.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"role_id": schema.StringAttribute{
				Description: "ID of the role which is going to be applied on the specified user.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "ID of the workspace which the role is goling to be applied on the scope of that.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"items": schema.MapAttribute{
				Description: "Items of the role-user binding.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *roleUserBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.RoleUser
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(plan.WorkspaceID.ValueString())
	userID, _ := uuid.FromString(plan.UserID.ValueString())
	roleID, _ := uuid.FromString(plan.RoleID.ValueString())

	items := make(map[string]string)
	diags = plan.Items.ElementsAs(ctx, &items, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.iamClient.BindRoleToUser(&workspaceID, &roleID, &userID, items)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_USER_BINDING_CREATE, errorMessage)
		return
	}
	if plan.UUID.ValueString() == "" {
		plan.UUID = types.StringValue(uuid.NewV4().String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *roleUserBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.RoleUser
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(state.WorkspaceID.ValueString())
	userID, _ := uuid.FromString(state.UserID.ValueString())
	roleID, _ := uuid.FromString(state.RoleID.ValueString())

	items, err := r.iamClient.GetBindedRoleToUserItems(&workspaceID, &roleID, &userID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_USER_BINDING_READ, errorMessage)
		return
	}

	state.Items, diags = types.MapValueFrom(ctx, types.StringType, items)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
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

func (r *roleUserBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.RoleUser
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(state.WorkspaceID.ValueString())
	userID, _ := uuid.FromString(state.UserID.ValueString())
	roleID, _ := uuid.FromString(state.RoleID.ValueString())

	items := make(map[string]string)

	err := r.iamClient.UnbindRoleFromUser(&workspaceID, &roleID, &userID, items)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_USER_BINDING_DELETE, errorMessage)
		return
	}
}

func (r *roleUserBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state models.RoleUser
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(state.WorkspaceID.ValueString())
	userID, _ := uuid.FromString(state.UserID.ValueString())
	roleID, _ := uuid.FromString(state.RoleID.ValueString())
	items := make(map[string]string)

	err := r.iamClient.UnbindRoleFromUser(&workspaceID, &roleID, &userID, items)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_USER_BINDING_DELETE_TO_UPDATE, errorMessage)
		return
	}

	var plan models.RoleUser
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ = uuid.FromString(plan.WorkspaceID.ValueString())
	userID, _ = uuid.FromString(plan.UserID.ValueString())
	roleID, _ = uuid.FromString(plan.RoleID.ValueString())

	items = make(map[string]string)
	diags = plan.Items.ElementsAs(ctx, &items, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = r.iamClient.BindRoleToUser(&workspaceID, &roleID, &userID, items)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_USER_BINDING_CREATE, errorMessage)
		return
	}

	_, err = r.iamClient.GetBindedRoleToUserItems(&workspaceID, &roleID, &userID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_USER_BINDING_READ, errorMessage)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *roleUserBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := strings.Split(req.ID, ":")
	if len(importID) < 3 {
		resp.Diagnostics.AddError(
			"Error Importing Role-User Attachment",
			"The given ID must be in the form of {role id}:{user id}:{workspace id}",
		)
		return
	}
	roleUUID := importID[0]
	userUUID := importID[1]
	workspaceUUID := importID[2]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_id"), roleUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), userUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_id"), workspaceUUID)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
