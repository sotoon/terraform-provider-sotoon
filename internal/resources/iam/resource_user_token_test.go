//go:build acceptance
// +build acceptance

package iamresources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testutil "github.com/sotoon/terraform-provider-sotoon/internal/utils"
)

func TestAccUserToken_basic(t *testing.T) {
	name := "developer-token"

	cfg := fmt.Sprintf(`
%s

# create a user token for the current user (no explicit user_id needed per your example)
resource "sotoon_iam_user_token" "me" {
  name      = "%s"
  expire_at = "2030-01-01T00:00:00Z"
}

# You expose the token value as 'value' (sensitive). We won't assert it directly.
output "user_token_id" {
  value = sotoon_iam_user_token.me.id
}
`, testutil.BaseProviderBlock(), name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		ProviderFactories: testutil.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_user_token.me", "name", name),
					resource.TestCheckResourceAttrSet("sotoon_iam_user_token.me", "id"),
					resource.TestCheckOutput("user_token_id", "${sotoon_iam_user_token.me.id}"),
				),
			},
			{
				ResourceName:      "sotoon_iam_user_token.me",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
