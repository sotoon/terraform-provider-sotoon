//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceUserToken_basic(t *testing.T) {
	name := randName("tf-acc-su")

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_service_user" "su" {
  name        = "%s"
  description = "acceptance service user"
}

resource "sotoon_iam_service_user_token" "tok" {
  service_user_id = sotoon_iam_service_user.su.id
}
`, baseProviderBlock(), name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps:             []resource.TestStep{{Config: cfg}},
	})
}
