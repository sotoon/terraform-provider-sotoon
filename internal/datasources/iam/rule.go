package iam

import (
	"context"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/client"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/models"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/utils"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ruleDataSource{}
	_ datasource.DataSourceWithConfigure = &ruleDataSource{}
)

func NewRuleDataSource() datasource.DataSource {
	return &ruleDataSource{}
}

type ruleDataSource struct {
	iamClient client.Client
}

func (d *ruleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_rule"
}

func (d *ruleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.iamClient = req.ProviderData.(*utils.SotoonConfigDataHolder).IAMClient
}

func (d *ruleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a IAM rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the rule.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the rule.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "ID of the Workspace which this rule is defined in that. (Default: global rules)",
				Optional:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"actions": schema.ListAttribute{
				Description: "List of the actions which this rule applied on that.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"service": schema.StringAttribute{
				Description: "Service which this rule defined on that.",
				Computed:    true,
			},
			"path": schema.StringAttribute{
				Description: "The path which the rule applied on that",
				Computed:    true,
			},
			"is_denial": schema.BoolAttribute{
				Description: "Defines is this rule denial or not.",
				Computed:    true,
			},
		},
	}
}

func (d *ruleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.Rule
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceName, err := utils.GetWokspaceNameFromUUID(state.WorkspaceUUID, &d.iamClient, resp.Diagnostics)
	if err != nil {
		return
	}

	rule, err := d.iamClient.GetRuleByName(state.Name.ValueString(), workspaceName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Rule",
			err.Error(),
		)
		return
	}

	state = *models.NewRuleModelFromRuleObject(rule)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
