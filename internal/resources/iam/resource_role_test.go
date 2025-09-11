//go:build acceptance
// +build acceptance

package iamresources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testutil "github.com/sotoon/terraform-provider-sotoon/internal/utils"
)

func TestAccRole_basic(t *testing.T) {
	name := testutil.RandName("tf-acc-role")

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_role" "test" {
  name = "%s"
}
`, testutil.BaseProviderBlock(), name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		ProviderFactories: testutil.ProviderFactories,
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
