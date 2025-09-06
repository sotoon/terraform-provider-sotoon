//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceGroupRoles_basic(t *testing.T) {
	ws := os.Getenv("SOTOON_WORKSPACE_ID")
	if ws == "" {
		t.Fatalf("SOTOON_WORKSPACE_ID must be set")
	}
	groupName := randName("tf-acc-group")
	roleName := randName("tf-acc-role")

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_group" "g" {
  name        = "%s"
  description = "acceptance group for group_roles DS"
}

resource "sotoon_iam_role" "r" {
  name = "%s"
}

resource "sotoon_iam_group_role" "bind" {
  group_id = sotoon_iam_group.g.id
  role_ids = [sotoon_iam_role.r.id]
}

data "sotoon_iam_group_roles" "gr" {
  workspace_id = "%s"
  group_id     = sotoon_iam_group.g.id
  depends_on   = [sotoon_iam_group_role.bind]
}

output "has_expected_role" {
  value = contains(data.sotoon_iam_group_roles.gr.roles[*].name, "%s")
}
`, baseProviderBlock(), groupName, roleName, ws, roleName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{Config: cfg, Check: resource.TestCheckOutput("has_expected_role", "true")},
		},
	})
}
