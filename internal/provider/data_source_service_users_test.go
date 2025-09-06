//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceServiceUsers_basic(t *testing.T) {
	name := randName("tf-acc-su")

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_service_user" "su" {
  name        = "%s"
  description = "acceptance service user"
}

data "sotoon_iam_service_users" "all" {
  depends_on = [sotoon_iam_service_user.su]
}

output "has_su" {
  value = contains(data.sotoon_iam_service_users.all.users[*].name, "%s")
}
`, baseProviderBlock(), name, name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{Config: cfg, Check: resource.TestCheckOutput("has_su", "true")},
		},
	})
}
