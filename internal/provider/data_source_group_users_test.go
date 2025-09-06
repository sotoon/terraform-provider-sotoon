//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceGroupUsers_emptyGroup(t *testing.T) {
	ws := os.Getenv("SOTOON_WORKSPACE_ID")
	if ws == "" {
		t.Fatalf("SOTOON_WORKSPACE_ID must be set")
	}
	name := randName("tf-acc-group")

	cfg := fmt.Sprintf(`
%s

# create an empty group
resource "sotoon_iam_group" "g" {
  name        = "%s"
  description = "acceptance group for group_users DS"
}

data "sotoon_iam_group_users" "gu" {
  workspace_id = "%s"
  group_id     = sotoon_iam_group.g.id
}

output "group_user_count" {
  value = length(data.sotoon_iam_group_users.gu.users)
}
`, baseProviderBlock(), name, ws)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check:  resource.TestCheckOutput("group_user_count", "0"),
			},
		},
	})
}
