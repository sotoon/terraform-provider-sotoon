//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceUsers_basic(t *testing.T) {
	ws := os.Getenv("SOTOON_WORKSPACE_ID")
	if ws == "" {
		t.Fatalf("SOTOON_WORKSPACE_ID must be set")
	}
	expectedEmail := os.Getenv("SOTOON_TEST_USER_EMAIL")

	cfg := fmt.Sprintf(`
%s

data "sotoon_iam_users" "all" {
  workspace_id = "%s"
}

output "has_expected_email" {
  value = %s
}
`, baseProviderBlock(), ws,
		func() string {
			if expectedEmail != "" {
				return fmt.Sprintf(`contains(data.sotoon_iam_users.all.users[*].email, "%s")`, expectedEmail)
			}
			return `can(data.sotoon_iam_users.all.users)`
		}(),
	)

	check := resource.TestCheckOutput("has_expected_email", "true")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps:             []resource.TestStep{{Config: cfg, Check: check}},
	})
}
