//go:build acceptance
// +build acceptance

package iamresources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testutil "github.com/sotoon/terraform-provider-sotoon/internal/utils"
)

func TestAccUser_basic(t *testing.T) {
	email := os.Getenv("SOTOON_TEST_USER_EMAIL")
	if email == "" {
		t.Skip("SOTOON_TEST_USER_EMAIL not set; skipping user adopt test to avoid invites")
	}

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_user" "test" {
  email = "%s"
}

output "user_email" { value = sotoon_iam_user.test.email }
`, testutil.BaseProviderBlock(), email)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		ProviderFactories: testutil.ProviderFactories,
		Steps: []resource.TestStep{
			{Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_user.test", "email", email),
					resource.TestCheckResourceAttrSet("sotoon_iam_user.test", "id"),
					resource.TestCheckResourceAttrSet("sotoon_iam_user.test", "user_uuid"),
				),
			},
		},
	})
}
