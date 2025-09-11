//go:build acceptance
// +build acceptance

package iamresources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testutil "github.com/sotoon/terraform-provider-sotoon/internal/utils"
)

func TestAccServiceUserRole_basic(t *testing.T) {
	suName := testutil.RandName("tf-acc-su")
	roleName := testutil.RandName("tf-acc-role")

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_service_user" "su" {
  name        = "%s"
  description = "acceptance service user"
}

resource "sotoon_iam_role" "r" {
  name = "%s"
}

resource "sotoon_iam_service_user_role" "bind" {
  role_id          = sotoon_iam_role.r.id
  service_user_ids = [sotoon_iam_service_user.su.id]
}

output "su_role_binding_id" {
  value = sotoon_iam_service_user_role.bind.id
}
`, testutil.BaseProviderBlock(), suName, roleName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		ProviderFactories: testutil.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sotoon_iam_service_user_role.bind", "id"),
				),
			},
			{
				ResourceName:      "sotoon_iam_service_user_role.bind",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
