//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceUser_basic(t *testing.T) {
	name := randName("tf-acc-su")

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_service_user" "test" {
  name        = "%s"
  description = "created by acceptance test"
}
`, baseProviderBlock(), name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps:             []resource.TestStep{{Config: cfg}},
	})
}
