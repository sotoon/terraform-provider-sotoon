package iam

import (
	"context"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/client"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/models"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	uuid "github.com/satori/go.uuid"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	iamClient client.Client
}

func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_user"
}

func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.iamClient = req.ProviderData.(*utils.SotoonConfigDataHolder).IAMClient
}

func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve a user from a workspace",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the user.",
				Computed:    true,
			},
			"workspace_id": schema.StringAttribute{
				Description: "ID of the workspace which user is retrived from that.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the retrieved user.",
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "Email address of the retrieved user.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(utils.REGEX_EMAIL, "value must be a valid email"),
				},
			},
		},
	}
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.User
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceUUID, _ := uuid.FromString(state.WorkspaceUUID.ValueString())

	userEmail := state.Email.ValueString()
	user, err := d.iamClient.GetUserByEmail(userEmail, &workspaceUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Users By Email",
			err.Error(),
		)
		return
	}
	state = *models.NewUserModelFromUserObject(user, state.WorkspaceUUID.ValueString())

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
