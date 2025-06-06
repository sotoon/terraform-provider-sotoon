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

	uuid "github.com/satori/go.uuid"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &groupDataSource{}
	_ datasource.DataSourceWithConfigure = &groupDataSource{}
)

func NewGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

type groupDataSource struct {
	iamClient client.Client
}

func (d *groupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_group"
}

func (d *groupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.iamClient = req.ProviderData.(*utils.SotoonConfigDataHolder).IAMClient
}

func (d *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a group in the specified workspace",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the group.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the group.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"workspace_id": schema.StringAttribute{
				Description: "ID of the workspace which group defined in that.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
		},
	}
}

func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.Group
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceUUID, _ := uuid.FromString(state.WorkspaceUUID.ValueString())

	groupName := state.Name.ValueString()

	workspace, err := d.iamClient.GetWorkspace(&workspaceUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Retreive Group Workspace",
			err.Error(),
		)
		return
	}

	group, err := d.iamClient.GetGroupByName(workspace.Name, groupName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Retreive Group",
			err.Error(),
		)
		return
	}

	state = *models.NewGroupModelFromGroupObject(group)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
