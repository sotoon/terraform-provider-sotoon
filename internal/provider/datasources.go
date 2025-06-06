package provider

import (
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/datasources/base"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/datasources/iam"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func GetDataSources() []func() datasource.DataSource {
	datasources := []func() datasource.DataSource{
		iam.NewServiceUserDataSource,
		iam.NewServiceUsersDataSource,
		iam.NewUserDataSource,
		iam.NewUsersDataSource,
		iam.NewRuleDataSource,
		iam.NewRoleDataSource,
		iam.NewGroupDataSource,

		base.NewWorkspaceDataSource,
		base.NewServiceDataSource,
	}

	return datasources
}
