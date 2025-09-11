//go:build acceptance
// +build acceptance

package iamresources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testutil "github.com/sotoon/terraform-provider-sotoon/internal/utils"
)

func TestAccGroup_basic(t *testing.T) {
	name := testutil.RandName("tf-acc-group")
	desc := "created by acceptance test"

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_group" "test" {
  name        = "%s"
  description = "%s"
}

`, testutil.BaseProviderBlock(), name, desc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		ProviderFactories: testutil.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_group.test", "name", name),
					resource.TestCheckResourceAttr("sotoon_iam_group.test", "description", desc),
					resource.TestCheckResourceAttrSet("sotoon_iam_group.test", "id"),
				),
			},
			{
				ResourceName:      "sotoon_iam_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
