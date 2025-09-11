//go:build acceptance
// +build acceptance

package iamdatasources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testutil "github.com/sotoon/terraform-provider-sotoon/internal/utils"
)

func TestAccDataSourceUserRoles_basic(t *testing.T) {
	userUUID := os.Getenv("SOTOON_TEST_USER_ID")
	if userUUID == "" {
		t.Skip("SOTOON_TEST_USER_ID not set; skipping user roles DS")
	}

	cfg := fmt.Sprintf(`
%s

data "sotoon_iam_user_roles" "ur" {
  user_id = "%s"
}

output "roles_ok" { value = can(data.sotoon_iam_user_roles.ur.roles) }
`, testutil.BaseProviderBlock(), userUUID)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		ProviderFactories: testutil.ProviderFactories,
		Steps:             []resource.TestStep{{Config: cfg, Check: resource.TestCheckOutput("roles_ok", "true")}},
	})
}
