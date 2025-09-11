//go:build acceptance
// +build acceptance

package iamresources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testutil "github.com/sotoon/terraform-provider-sotoon/internal/utils"
)

func TestAccServiceUserGroup_basic(t *testing.T) {
	groupName := testutil.RandName("tf-acc-group")
	suName := testutil.RandName("tf-acc-su")

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
`, testutil.BaseProviderBlock(), groupName, suName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		ProviderFactories: testutil.ProviderFactories,
		Steps:             []resource.TestStep{{Config: cfg}},
	})
}
