//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceGroupDetails_basic(t *testing.T) {
	ws := os.Getenv("SOTOON_WORKSPACE_ID")
	if ws == "" {
		t.Fatalf("SOTOON_WORKSPACE_ID must be set")
	}
	name := randName("tf-acc-group")

	cfg := fmt.Sprintf(`
%s

# create a group
resource "sotoon_iam_group" "g" {
  name        = "%s"
  description = "acceptance group for details DS"
}

data "sotoon_iam_group_details" "details" {
  workspace_id = "%s"
  group_id     = sotoon_iam_group.g.id
}

output "ds_group_name" {
  value = data.sotoon_iam_group_details.details.name
}
`, baseProviderBlock(), name, ws)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check:  resource.TestCheckOutput("ds_group_name", name),
			},
		},
	})
}
