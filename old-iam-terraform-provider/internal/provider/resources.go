package provider

import (
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/resources/iam"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func GetResources() []func() resource.Resource {
	datasources := []func() resource.Resource{
		iam.NewGroupResource,
		iam.NewRuleResource,
		iam.NewRoleResource,
		iam.NewRoleUserBindingResource,
		iam.NewRoleServiceUserBindingResource,
		iam.NewRoleGroupBindingResource,
		iam.NewUserGroupBindingResource,
		iam.NewServiceUserGroupBindingResource,
	}

	return datasources
}
