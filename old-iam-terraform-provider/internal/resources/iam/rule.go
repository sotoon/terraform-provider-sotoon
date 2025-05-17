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
	_ resource.Resource                = &ruleResource{}
	_ resource.ResourceWithConfigure   = &ruleResource{}
	_ resource.ResourceWithImportState = &ruleResource{}
)

func NewRuleResource() resource.Resource {
	return &ruleResource{}
}

type ruleResource struct {
	iamClient client.Client
}

func (r *ruleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_rule"
}

func (r *ruleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.iamClient = req.ProviderData.(*utils.SotoonConfigDataHolder).IAMClient
}

func (r *ruleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a IAM rule instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the rule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the rule.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "ID of the Workspace which the rule is  going to be defined in that.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"actions": schema.ListAttribute{
				Description: "List of the actions which this rule applied on that.",
				ElementType: types.StringType,
				Required:    true,
			},
			"service": schema.StringAttribute{
				Description: "The service which this rule applied on that. The ID and name of the services are same and there is no difference in passing ID or name of a service.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"path": schema.StringAttribute{
				Description: "The path that the rule is going to be applied on that",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"is_denial": schema.BoolAttribute{
				Description: "Defines is this rule denial or not.",
				Required:    true,
			},
		},
	}
}

func (r *ruleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.Rule
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleName := plan.Name.ValueString()

	workspaceUUID, _ := uuid.FromString(plan.WorkspaceUUID.ValueString())

	actionsList := utils.GetStringListFromStringValueList(plan.Actions)
	isDenial := plan.IsDenial.ValueBool()
	serviceName := plan.Service.ValueString()
	path := plan.Path.ValueString()
	object := utils.CreateRRI(plan.WorkspaceUUID.ValueString(), serviceName, path)

	createdRule, err := r.iamClient.CreateRule(ruleName, &workspaceUUID, actionsList, object, isDenial)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_RULE_CREATE, errorMessage)
		return
	}

	plan = *models.NewRuleModelFromRuleObject(createdRule)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ruleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.Rule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleUUID, _ := uuid.FromString(state.UUID.ValueString())

	workspaceUUID, _ := uuid.FromString(state.WorkspaceUUID.ValueString())

	rule, err := r.iamClient.GetRule(&ruleUUID, &workspaceUUID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_RULE_READ, errorMessage)
		return
	}

	state = *models.NewRuleModelFromRuleObject(rule)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ruleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.Rule
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleUUID, name, workspaceUUID, actions, path, service, deny, err := plan.GetParams()
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Data",
			err.Error(),
		)
		return
	}

	object := utils.CreateRRI(workspaceUUID.String(), service, path)
	updatedRule, err := r.iamClient.UpdateRule(ruleUUID, name, workspaceUUID, actions, object, deny)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_RULE_UPDATE, errorMessage)
		return
	}

	plan = *models.NewRuleModelFromRuleObject(updatedRule)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ruleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.Rule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleUUID, _ := uuid.FromString(state.UUID.ValueString())

	workspaceUUID, _ := uuid.FromString(state.WorkspaceUUID.ValueString())

	err := r.iamClient.DeleteRule(&ruleUUID, &workspaceUUID)
	if err != nil {
		errorMessage := utils.GetIAMErrorMessage(err)
		resp.Diagnostics.AddError(utils.ERROR_RULE_DELETE, errorMessage)
		return
	}
}

func (r *ruleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := strings.Split(req.ID, ":")
	if len(importID) < 2 {
		resp.Diagnostics.AddError(
			"Error Importing Rule",
			"The given ID must be in the form of {rule id}:{workspace id}",
		)
		return
	}
	ruleUUID := importID[0]
	workspaceUUID := importID[1]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), ruleUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_id"), workspaceUUID)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
