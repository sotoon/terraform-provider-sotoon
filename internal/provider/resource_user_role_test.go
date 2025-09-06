//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserRole_basic(t *testing.T) {
	email := randEmail()
	roleName := randName("tf-acc-role")

	cfg := fmt.Sprintf(`
%s

# create a user
resource "sotoon_iam_user" "u" {
  email = "%s"
}

# create a role
resource "sotoon_iam_role" "r" {
  name = "%s"
}

# bind the role to the user
resource "sotoon_iam_user_role" "bind" {
  role_id  = sotoon_iam_role.r.id
  user_ids = [sotoon_iam_user.u.user_uuid]
}

# check roles via data source
data "sotoon_iam_user_roles" "ur" {
  user_id = sotoon_iam_user.u.user_uuid
}

output "user_role_names" {
  value = data.sotoon_iam_user_roles.ur.roles[*].name
}
`, baseProviderBlock(), email, roleName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sotoon_iam_user_role.bind", "id"),
					resource.TestCheckOutput("user_role_names", fmt.Sprintf("[%s]", roleName)),
				),
			},
			{
				ResourceName:      "sotoon_iam_user_role.bind",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
