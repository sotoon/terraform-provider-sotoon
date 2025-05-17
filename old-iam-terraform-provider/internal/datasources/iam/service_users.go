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
	_ datasource.DataSource              = &serviceUsersDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceUsersDataSource{}
)

func NewServiceUsersDataSource() datasource.DataSource {
	return &serviceUsersDataSource{}
}

type serviceUsersDataSource struct {
	iamClient client.Client
}

func (d *serviceUsersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_service_users"
}

func (d *serviceUsersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.iamClient = req.ProviderData.(*utils.SotoonConfigDataHolder).IAMClient
}

// Schema defines the schema for the data source.
func (d *serviceUsersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of service-users from a workspace of group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of service-user set.",
				Computed:    true,
			},
			"workspace_id": schema.StringAttribute{
				Description: "ID of the workspace which this user is retrieved from.",
				Required:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"group_id": schema.StringAttribute{
				Description: "ID of the workspace which this user is retrieved from. Set this attribute if you want to get list of service-users of a group.",
				Optional:    true,
				Validators: []validator.String{
					validators.UUID(),
				},
			},
			"service_users": schema.ListNestedAttribute{
				Description: "List of retrieved service users.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the service user.",
							Computed:    true,
						},
						"workspace_id": schema.StringAttribute{
							Description: "Workspace which the service-user retrived from that.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the service-user.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *serviceUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state models.ServiceUsers
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceUUID, _ := uuid.FromString(state.WorkspaceUUID.ValueString())

	var serviceUsers []*bepatypes.ServiceUser

	var err error
	if !state.GroupUUID.IsNull() {
		GroupUUID, _ := uuid.FromString(state.GroupUUID.ValueString())
		state.UUID = state.GroupUUID
		serviceUsers, err = d.iamClient.GetAllGroupServiceUsers(&workspaceUUID, &GroupUUID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Group Service Users",
				err.Error(),
			)
			return
		}
	} else {
		state.UUID = state.WorkspaceUUID
		serviceUsers, err = d.iamClient.GetServiceUsers(&workspaceUUID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Workspace Service Users",
				err.Error(),
			)
			return
		}
	}

	// Map response body to model
	for _, serviceUser := range serviceUsers {
		newServiceUserModel := models.NewServiceUserModelFromServiceUserObject(serviceUser, workspaceUUID.String())
		state.ServiceUsers = append(state.ServiceUsers, *newServiceUserModel)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
