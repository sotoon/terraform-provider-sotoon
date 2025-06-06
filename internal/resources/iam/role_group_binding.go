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
	_ resource.Resource                = &roleGroupBindingResource{}
	_ resource.ResourceWithConfigure   = &roleGroupBindingResource{}
	_ resource.ResourceWithImportState = &roleGroupBindingResource{}
)

func NewRoleGroupBindingResource() resource.Resource {
	return &roleGroupBindingResource{}
}

type roleGroupBindingResource struct {
	iamClient bepa.Client
}

func (r *roleGroupBindingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.iamClient = req.ProviderData.(*utils.SotoonConfigDataHolder).IAMClient
}

func (r *roleGroupBindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_role_group_binding"
}

func (r *roleGroupBindingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Relation between role and group. Existance of an instance from this resource will applies an specified role to all memebers of the group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the role-group binding",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "ID of the workspace that group is defined in that.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"group_id": schema.StringAttribute{
				Description: "ID of the group which the role going to be binded to that.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"role_id": schema.StringAttribute{
				Description: "ID of the role which is going to be applied on the members of the specified group.",
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

func (r *roleGroupBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.RoleGroup
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(plan.WorkspaceID.ValueString())
	groupID, _ := uuid.FromString(plan.GroupID.ValueString())
	roleID, _ := uuid.FromString(plan.RoleID.ValueString())

	items := make(map[string]string)
	diags = plan.Items.ElementsAs(ctx, &items, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.iamClient.BindRoleToGroup(&workspaceID, &roleID, &groupID, items)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_GROUP_BINDING_CREATE, errorMessage)
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

func (r *roleGroupBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.RoleGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(state.WorkspaceID.ValueString())
	groupID, _ := uuid.FromString(state.GroupID.ValueString())
	roleID, _ := uuid.FromString(state.RoleID.ValueString())

	items, err := r.iamClient.GetBindedRoleToGroupItems(&workspaceID, &roleID, &groupID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_GROUP_BINDING_READ, errorMessage)
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

func (r *roleGroupBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.RoleGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(state.WorkspaceID.ValueString())
	groupID, _ := uuid.FromString(state.GroupID.ValueString())
	roleID, _ := uuid.FromString(state.RoleID.ValueString())

	items := make(map[string]string)

	err := r.iamClient.UnbindRoleFromGroup(&workspaceID, &roleID, &groupID, items)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_GROUP_BINDING_DELETE, errorMessage)
		return
	}
}

func (r *roleGroupBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state models.RoleGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ := uuid.FromString(state.WorkspaceID.ValueString())
	groupID, _ := uuid.FromString(state.GroupID.ValueString())
	roleID, _ := uuid.FromString(state.RoleID.ValueString())
	items := make(map[string]string)

	err := r.iamClient.UnbindRoleFromGroup(&workspaceID, &roleID, &groupID, items)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_GROUP_BINDING_DELETE_TO_UPDATE, errorMessage)
		return
	}

	var plan models.RoleGroup
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceID, _ = uuid.FromString(plan.WorkspaceID.ValueString())
	groupID, _ = uuid.FromString(plan.GroupID.ValueString())
	roleID, _ = uuid.FromString(plan.RoleID.ValueString())

	items = make(map[string]string)
	diags = plan.Items.ElementsAs(ctx, &items, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = r.iamClient.BindRoleToGroup(&workspaceID, &roleID, &groupID, items)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_GROUP_BINDING_CREATE, errorMessage)
		return
	}

	_, err = r.iamClient.GetBindedRoleToGroupItems(&workspaceID, &roleID, &groupID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_GROUP_BINDING_ITEMS_READ, errorMessage)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *roleGroupBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := strings.Split(req.ID, ":")
	if len(importID) < 3 {
		resp.Diagnostics.AddError(
			"Error Importing Role-Group Attachment",
			"The given ID must be in the form of {role id}:{group id}:{workspace id}",
		)
		return
	}
	roleUUID := importID[0]
	groupUUID := importID[1]
	workspaceUUID := importID[2]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_id"), roleUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), groupUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_id"), workspaceUUID)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
