//go:build acceptance
// +build acceptance

package iamdatasources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testutil "github.com/sotoon/terraform-provider-sotoon/internal/utils"
)

func TestAccDataSourceServiceUserTokens_basic(t *testing.T) {
	name := testutil.RandName("tf-acc-su")
	cfg1 := fmt.Sprintf(`
%s

resource "sotoon_iam_service_user" "su" {
  name        = "%s"
  description = "acceptance service user"
}

resource "sotoon_iam_service_user_token" "tok" {
  service_user_id = sotoon_iam_service_user.su.id
}
`, testutil.BaseProviderBlock(), name)
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
`, testutil.BaseProviderBlock(), name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		ProviderFactories: testutil.ProviderFactories,
		Steps: []resource.TestStep{
			{Config: cfg1},
			{Config: cfg2, Check: resource.TestCheckOutput("tokens_non_empty", "true")},
		},
	})
}
