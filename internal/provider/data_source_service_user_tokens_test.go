//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceServiceUserTokens_basic(t *testing.T) {
	name := randName("tf-acc-su")
	cfg1 := fmt.Sprintf(`
%s

resource "sotoon_iam_service_user" "su" {
  name        = "%s"
  description = "acceptance service user"
}

resource "sotoon_iam_service_user_token" "tok" {
  service_user_id = sotoon_iam_service_user.su.id
}
`, baseProviderBlock(), name)
	cfg2 := fmt.Sprintf(`
%s

resource "sotoon_iam_service_user" "su" {
  name        = "%s"
  description = "acceptance service user"
}

# token exists from step 1
resource "sotoon_iam_service_user_token" "tok" {
  service_user_id = sotoon_iam_service_user.su.id
}

data "sotoon_iam_service_user_tokens" "list" {
  service_user_id = sotoon_iam_service_user.su.id
}

output "tokens_non_empty" {
  value = length(data.sotoon_iam_service_user_tokens.list.tokens) > 0
}
`, baseProviderBlock(), name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{Config: cfg1},
			{Config: cfg2, Check: resource.TestCheckOutput("tokens_non_empty", "true")},
		},
	})
}
