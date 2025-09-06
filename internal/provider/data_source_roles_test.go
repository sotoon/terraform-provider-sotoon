//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceRoles_basic(t *testing.T) {
	ws := os.Getenv("SOTOON_WORKSPACE_ID")
	if ws == "" {
		t.Fatalf("SOTOON_WORKSPACE_ID must be set")
	}
	roleName := randName("tf-acc-role")

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_role" "seed" {
  name = "%s"
}

data "sotoon_iam_roles" "all" {
  workspace_id = "%s"
  depends_on   = [sotoon_iam_role.seed]
}

output "has_seeded_role" {
  value = contains(data.sotoon_iam_roles.all.roles[*].name, "%s")
}
`, baseProviderBlock(), roleName, ws, roleName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{Config: cfg, Check: resource.TestCheckOutput("has_seeded_role", "true")},
		},
	})
}
