package models

import (
	iamtypes "git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Workspace struct {
	UUID             types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	OrganizationUUID types.String `tfsdk:"organization_id"`
}

func NewWorkspaceModelFromWorkspaceObject(workspace *iamtypes.Workspace) *Workspace {
	return &Workspace{
		UUID:             types.StringValue((*workspace.UUID).String()),
		Name:             types.StringValue(workspace.Name),
		OrganizationUUID: types.StringValue(workspace.Organization.String()),
	}
}

type Service struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
