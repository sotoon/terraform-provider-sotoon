//go:build acceptance
// +build acceptance

package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceRules_basic(t *testing.T) {
	ws := os.Getenv("SOTOON_WORKSPACE_ID")
	if ws == "" {
		t.Fatalf("SOTOON_WORKSPACE_ID must be set")
	}

	cfg := `
` + baseProviderBlock() + `
data "sotoon_iam_rules" "all" {
  workspace_id = "` + ws + `"
}

output "rules_attr_present" { value = can(data.sotoon_iam_rules.all.rules) }
output "rules_len_nonneg"  { value = length(data.sotoon_iam_rules.all.rules) >= 0 }
`

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("rules_attr_present", "true"),
					resource.TestCheckOutput("rules_len_nonneg", "true"),
				),
			},
		},
	})
}
