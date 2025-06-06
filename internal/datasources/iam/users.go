package iam

import (
	"context"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/client"
	bepatypes "git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/models"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/utils"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	uuid "github.com/satori/go.uuid"
)

var (
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

type usersDataSource struct {
	iamClient client.Client
}

func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_users"
}

func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.iamClient = req.ProviderData.(*utils.SotoonConfigDataHolder).IAMClient
}

// Schema defines the schema for the data source.
func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches list of users from the specified workspace",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of user set.",
				Computed:    true,
			},
			"workspace_id": schema.StringAttribute{
				Description: "ID of the workspace which users will be retrieved.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"group_id": schema.StringAttribute{
				Description: "ID of the group which users will be retrieved. Set this attribute if you want to get list of users of a group.",
				Optional:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"users": schema.ListNestedAttribute{
				Description: "List of retrieved users.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the user.",
							Required:    true,
						},
						"workspace_id": schema.StringAttribute{
							Description: "Workspace which user is retrived from that.",
							Required:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the retrived user.",
							Required:    true,
						},
						"email": schema.StringAttribute{
							Description: "Email address of the retrived user.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state models.Users
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceUUID, _ := uuid.FromString(state.WorkspaceUUID.ValueString())

	var users []*bepatypes.User

	// TODO: clean this section
	if !state.GroupUUID.IsNull() {
		groupUUID, _ := uuid.FromString(state.GroupUUID.ValueString())
		state.UUID = state.GroupUUID
		var err error
		users, err = d.iamClient.GetAllGroupUsers(&workspaceUUID, &groupUUID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Group Users",
				err.Error(),
			)
			return
		}
	} else {
		var err error
		state.UUID = state.WorkspaceUUID
		users, err = d.iamClient.GetWorkspaceUsers(&workspaceUUID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Workspace Users",
				err.Error(),
			)
			return
		}
	}

	// Map response body to model
	for _, user := range users {
		newUserModel := models.NewUserModelFromUserObject(user, workspaceUUID.String())
		state.Users = append(state.Users, *newUserModel)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
