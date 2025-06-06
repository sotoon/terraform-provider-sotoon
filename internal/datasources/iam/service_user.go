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
	_ datasource.DataSource              = &serviceUserDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceUserDataSource{}
)

func NewServiceUserDataSource() datasource.DataSource {
	return &serviceUserDataSource{}
}

type serviceUserDataSource struct {
	iamClient client.Client
}

func (d *serviceUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_service_user"
}

func (d *serviceUserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.iamClient = req.ProviderData.(*utils.SotoonConfigDataHolder).IAMClient
}

func (d *serviceUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a service-user from the specified workspace",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the service-user.",
				Computed:    true,
			},
			"workspace_id": schema.StringAttribute{
				Description: "Workspace which service-user is retrived from that.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the service-user.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
		},
	}
}

func (d *serviceUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.ServiceUser
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceUUID, _ := uuid.FromString(state.WorkspaceUUID.ValueString())

	workspace, err := d.iamClient.GetWorkspace(&workspaceUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Workspace",
			err.Error(),
		)
		return
	}

	serviceUserName := state.Name.ValueString()
	// TODO: create iamClient get user service by name and workspace uuid and replace this
	serviceUser, err := d.iamClient.GetServiceUserByName(workspace.Name, serviceUserName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Service User",
			err.Error(),
		)
		return
	}

	state.UUID = types.StringValue((*serviceUser.UUID).String())

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
