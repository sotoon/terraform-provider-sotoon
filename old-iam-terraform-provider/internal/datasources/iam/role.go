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
	uuid "github.com/satori/go.uuid"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &roleDataSource{}
	_ datasource.DataSourceWithConfigure = &roleDataSource{}
)

func NewRoleDataSource() datasource.DataSource {
	return &roleDataSource{}
}

type roleDataSource struct {
	iamClient client.Client
}

func (d *roleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_role"
}

func (d *roleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.iamClient = req.ProviderData.(*utils.SotoonConfigDataHolder).IAMClient
}

func (d *roleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a IAM role.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the role.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the role.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "ID of the Workspace which role is defined in that. (Default: global roles)",
				Optional:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"rules": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the rule.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the rule.",
							Computed:    true,
						},
						"workspace_id": schema.StringAttribute{
							Description: "ID of the Workspace which rule is defined in that. (Default: global rules)",
							Computed:    true,
						},
						"actions": schema.ListAttribute{
							Description: "List of the actions which this rule is binded to that.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"service": schema.StringAttribute{
							Description: "Service which this rule defined on that.",
							Computed:    true,
						},
						"path": schema.StringAttribute{
							Description: "The path which this rule applied on that",
							Computed:    true,
						},
						"is_denial": schema.BoolAttribute{
							Description: "Defines is this rule denial or not.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *roleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.Role
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var workspaceUUIDString string
	if state.WorkspaceUUID.IsNull() {
		workspaceUUIDString = utils.GLOBAL_WORKSPACE
	} else {
		workspaceUUIDString = state.WorkspaceUUID.ValueString()
	}

	workspaceName, err := utils.GetWokspaceNameFromUUID(state.WorkspaceUUID, &d.iamClient, resp.Diagnostics)
	if err != nil {
		return
	}

	role, err := d.iamClient.GetRoleByName(state.Name.ValueString(), workspaceName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Role",
			err.Error(),
		)
		return
	}

	workspaceUUID, _ := uuid.FromString(workspaceUUIDString)
	rules, err := d.iamClient.GetRoleRules(role.UUID, &workspaceUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Role Rules",
			err.Error(),
		)
		return
	}

	newRoleModel := models.NewRoleModelFromRoleResAndRulesObject(role, rules)

	diags = resp.State.Set(ctx, &newRoleModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
