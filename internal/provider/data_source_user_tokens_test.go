//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceUserTokens_basic(t *testing.T) {
	name := "developer-token"

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_user_token" "me" {
  name      = "%s"
  expire_at = "2030-01-01T00:00:00Z"
}

data "sotoon_iam_user_tokens" "me" {
  depends_on = [sotoon_iam_user_token.me]
}

output "has_created_token" {
  value = contains(data.sotoon_iam_user_tokens.me.tokens[*].id, sotoon_iam_user_token.me.id)
}
`, baseProviderBlock(), name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{Config: cfg, Check: resource.TestCheckOutput("has_created_token", "true")},
		},
	})
}
