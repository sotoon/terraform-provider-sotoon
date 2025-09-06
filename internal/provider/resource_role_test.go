//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRole_basic(t *testing.T) {
	name := randName("tf-acc-role")

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_role" "test" {
  name = "%s"
}
`, baseProviderBlock(), name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_role.test", "name", name),
					resource.TestCheckResourceAttrSet("sotoon_iam_role.test", "id"),
				),
			},
			{
				ResourceName:      "sotoon_iam_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
