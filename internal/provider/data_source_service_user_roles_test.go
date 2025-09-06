//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceServiceUserRoles_basic(t *testing.T) {
	ws := os.Getenv("SOTOON_WORKSPACE_ID")
	if ws == "" {
		t.Fatalf("SOTOON_WORKSPACE_ID must be set")
	}

	suName := randName("tf-acc-su")
	roleName := randName("tf-acc-role")

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_service_user" "su" {
  name        = "%s"
  description = "acceptance service user"
}

resource "sotoon_iam_role" "r" {
  name = "%s"
}

resource "sotoon_iam_service_user_role" "bind" {
  role_id          = sotoon_iam_role.r.id
  service_user_ids = [sotoon_iam_service_user.su.id]
}

data "sotoon_iam_service_user_roles" "ds" {
  workspace_id    = "%s"
  service_user_id = sotoon_iam_service_user.su.id
  depends_on      = [sotoon_iam_service_user_role.bind]
}

output "has_expected_role" {
  value = contains(data.sotoon_iam_service_user_roles.ds.roles[*].name, "%s")
}
`, baseProviderBlock(), suName, roleName, ws, roleName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{Config: cfg, Check: resource.TestCheckOutput("has_expected_role", "true")},
		},
	})
}
