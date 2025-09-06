//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceUserGroup_basic(t *testing.T) {
	groupName := randName("tf-acc-group")
	suName := randName("tf-acc-su")

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_group" "g" {
  name        = "%s"
  description = "acceptance group"
}

resource "sotoon_iam_service_user" "su" {
  name        = "%s"
  description = "acceptance service user"
}

resource "sotoon_iam_service_user_group" "bind" {
  group_id         = sotoon_iam_group.g.id
  service_user_ids = [sotoon_iam_service_user.su.id]
}
`, baseProviderBlock(), groupName, suName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps:             []resource.TestStep{{Config: cfg}},
	})
}
