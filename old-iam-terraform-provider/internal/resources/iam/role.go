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

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &roleResource{}
	_ resource.ResourceWithConfigure   = &roleResource{}
	_ resource.ResourceWithImportState = &roleResource{}
)

func NewRoleResource() resource.Resource {
	return &roleResource{}
}

type roleResource struct {
	iamClient client.Client
}

func (r *roleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_role"
}

func (r *roleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.iamClient = req.ProviderData.(*utils.SotoonConfigDataHolder).IAMClient
}

func (r *roleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a IAM role instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the role.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the role.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "ID of the Workspace which role is going to be defined in that.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"rules": schema.SetNestedAttribute{
				Description: "List of the rules which this role contains.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the rule.",
							Required:    true,
							Validators: []validator.String{
								validators.UUID(),
							},
						},
					},
				},
			},
		},
	}
}

func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.RoleWithRulesID
	diags := req.Plan.Get(ctx, &plan)
	savedPlan := plan.ToRole()

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	roleName := savedPlan.Name.ValueString()
	workspaceUUID, _ := uuid.FromString(savedPlan.WorkspaceUUID.ValueString())

	createdRole, err := r.iamClient.CreateRole(roleName, &workspaceUUID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_CREATE, errorMessage)
		return
	}

	rules := []models.Rule{}
	for _, rule := range savedPlan.Rules {
		ruleUUID, _ := uuid.FromString(rule.UUID.ValueString())
		err = r.iamClient.BindRuleToRole(createdRole.UUID, &ruleUUID, &workspaceUUID)
		if err != nil {
			resp.Diagnostics.AddWarning(
				utils.ERROR_ROLE_RULE_BINDING_CREATE,
				utils.GetIAMErrorMessage(err),
			)
		} else {
			rules = append(rules, rule)
		}
	}
	savedPlan.Rules = rules

	savedPlan.UUID = types.StringValue((*createdRole.UUID).String())
	diags = resp.State.Set(ctx, savedPlan.ToRoleWithRulesID())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.RoleWithRulesID

	diags := req.State.Get(ctx, &state)
	savedState := state.ToRole()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleUUID, _ := uuid.FromString(savedState.UUID.ValueString())
	workspaceUUID, _ := uuid.FromString(savedState.WorkspaceUUID.ValueString())

	role, err := r.iamClient.GetRole(&roleUUID, &workspaceUUID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_READ, errorMessage)
		return
	}

	rules, err := r.iamClient.GetRoleRules(&roleUUID, &workspaceUUID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_RULE_BINDING_READ, errorMessage)
		return
	}

	savedState = models.NewRoleModelFromRoleResAndRulesObject(role, rules)

	diags = resp.State.Set(ctx, savedState.ToRoleWithRulesID())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.RoleWithRulesID
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var currentState models.RoleWithRulesID
	diags = req.State.Get(ctx, &currentState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	savedPlan := plan.ToRole()

	roleUUID, _ := uuid.FromString(savedPlan.UUID.ValueString())
	workspaceUUID, _ := uuid.FromString(savedPlan.WorkspaceUUID.ValueString())

	// Update role itself if changed
	if !plan.Equal(&currentState) {
		roleName := savedPlan.Name.ValueString()
		_, err := r.iamClient.UpdateRole(&roleUUID, roleName, &workspaceUUID)
		if err != nil {
			errorMessage := utils.GetIAMErrorMessage(err)
			resp.Diagnostics.AddError(utils.ERROR_ROLE_UPDATE, errorMessage)
			return
		}
	}

	// Update role bindings
	rulesMustBeBinded := models.GetAdditionsOfUUIDAttributes(plan.Rules, currentState.Rules)
	rulesMustBeUnbinded := models.GetAdditionsOfUUIDAttributes(currentState.Rules, plan.Rules)

	for _, rule := range rulesMustBeBinded {
		ruleUUID, _ := uuid.FromString(rule.UUID.ValueString())
		err := r.iamClient.BindRuleToRole(&roleUUID, &ruleUUID, &workspaceUUID)
		if err != nil {
			errorMessage := utils.GetIAMErrorMessage(err)
			resp.Diagnostics.AddError(utils.ERROR_ROLE_RULE_BINDING_CREATE, errorMessage)
			return
		}
	}

	for _, rule := range rulesMustBeUnbinded {
		ruleUUID, _ := uuid.FromString(rule.UUID.ValueString())
		err := r.iamClient.UnbindRuleFromRole(&roleUUID, &ruleUUID, &workspaceUUID)
		if err != nil {
			errorMessage := utils.GetIAMErrorMessage(err)
			resp.Diagnostics.AddError(utils.ERROR_ROLE_RULE_BINDING_DELETE, errorMessage)
			return
		}
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.RoleWithRulesID
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleUUID, _ := uuid.FromString(state.UUID.ValueString())
	workspaceUUID, _ := uuid.FromString(state.WorkspaceUUID.ValueString())

	err := r.iamClient.DeleteRole(&roleUUID, &workspaceUUID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_ROLE_DELETE, errorMessage)
		return
	}
}

func (r *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := strings.Split(req.ID, ":")
	if len(importID) < 2 {
		resp.Diagnostics.AddError(
			"Error Importing Role",
			"The given ID must be in the form of {role id}:{workspace id}",
		)
		return
	}
	roleUUID := importID[0]
	workspaceUUID := importID[1]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), roleUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_id"), workspaceUUID)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
