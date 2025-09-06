//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupRole_basic(t *testing.T) {
	ws := os.Getenv("SOTOON_WORKSPACE_ID")
	if ws == "" {
		t.Fatalf("SOTOON_WORKSPACE_ID must be set")
	}

	groupName := randName("tf-acc-group")
	roleName := randName("tf-acc-role")

	cfg := fmt.Sprintf(`
%s

# seed a group
resource "sotoon_iam_group" "g" {
  name        = "%s"
  description = "acceptance group for group_role binding"
}

# seed a role
resource "sotoon_iam_role" "r" {
  name = "%s"
}

# bind the role to the group (note: role_ids is a list)
resource "sotoon_iam_group_role" "bind" {
  group_id = sotoon_iam_group.g.id
  role_ids = [
    sotoon_iam_role.r.id
  ]
}

# verify via DS that the role appears on the group
data "sotoon_iam_group_roles" "gr" {
  workspace_id = "%s"
  group_id     = sotoon_iam_group.g.id
}

output "group_role_names" {
  value = data.sotoon_iam_group_roles.gr.roles[*].name
}
`, baseProviderBlock(), groupName, roleName, ws)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sotoon_iam_group_role.bind", "id"),
					resource.TestCheckOutput("group_role_names", roleName),
				),
			},
			{
				ResourceName:      "sotoon_iam_group_role.bind",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
