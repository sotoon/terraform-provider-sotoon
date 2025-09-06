//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceGroups_basic(t *testing.T) {
	ws := os.Getenv("SOTOON_WORKSPACE_ID")
	if ws == "" {
		t.Fatalf("SOTOON_WORKSPACE_ID must be set")
	}
	grpName := randName("tf-acc-group")

	cfgStep1 := fmt.Sprintf(`
%s

# seed a group first
resource "sotoon_iam_group" "seed" {
  name        = "%s"
  description = "created by acceptance test"
}
`, baseProviderBlock(), grpName)

	cfgStep2 := fmt.Sprintf(`
%s

resource "sotoon_iam_group" "seed" {
  name        = "%s"
  description = "created by acceptance test"
}

# read AFTER the resource exists
data "sotoon_iam_groups" "all" {
  workspace_id = "%s"
  depends_on   = [sotoon_iam_group.seed]
}

# boolean we can assert
output "has_seeded_group" {
  value = contains(data.sotoon_iam_groups.all.groups[*].name, "%s")
}
`, baseProviderBlock(), grpName, ws, grpName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{Config: cfgStep1},
			{
				Config: cfgStep2,
				Check:  resource.TestCheckOutput("has_seeded_group", "true"),
			},
		},
	})
}
