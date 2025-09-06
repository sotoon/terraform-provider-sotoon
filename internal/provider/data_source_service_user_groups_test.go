//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceServiceUserGroups_basic(t *testing.T) {
	ws := os.Getenv("SOTOON_WORKSPACE_ID")
	if ws == "" {
		t.Fatalf("SOTOON_WORKSPACE_ID must be set")
	}
	groupName := randName("tf-acc-group")
	suName := randName("tf-acc-su")

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_group" "g" {
  name        = "%s"
  description = "acceptance group for service_user_groups DS"
}

resource "sotoon_iam_service_user" "su" {
  name        = "%s"
  description = "acceptance service user"
}

resource "sotoon_iam_service_user_group" "bind" {
  group_id         = sotoon_iam_group.g.id
  service_user_ids = [sotoon_iam_service_user.su.id]
}

data "sotoon_iam_service_user_groups" "dev" {
  workspace_id = "%s"
  group_id     = sotoon_iam_group.g.id
  depends_on   = [sotoon_iam_service_user_group.bind]
}

output "has_su" {
  value = contains(data.sotoon_iam_service_user_groups.dev.service_users[*].name, "%s")
}
`, baseProviderBlock(), groupName, suName, ws, suName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{Config: cfg, Check: resource.TestCheckOutput("has_su", "true")},
		},
	})
}
