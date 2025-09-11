//go:build acceptance
// +build acceptance

package iamresources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testutil "github.com/sotoon/terraform-provider-sotoon/internal/utils"
)

func TestAccUserPublicKey_basic(t *testing.T) {
	title := "tf-acc-user-pk"
	pubKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7kG7aZb1W9n3mX6Lw0N2Q7pPq3yZV7g5r8jL9z8lQmA9oJQqD4lX9nB5u2jC4t4h1uWq1Ff5rjM4cO2q9f3B8nT9pA7vC2hX info@mail.com"

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_user_public_key" "pk" {
  title      = "%s"
  public_key = "%s"
}

output "user_pk_id"    { value = sotoon_iam_user_public_key.pk.id }
output "user_pk_title" { value = sotoon_iam_user_public_key.pk.title }
`, testutil.BaseProviderBlock(), title, pubKey)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		ProviderFactories: testutil.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_user_public_key.pk", "title", title),
					resource.TestCheckResourceAttrSet("sotoon_iam_user_public_key.pk", "id"),
					resource.TestCheckOutput("user_pk_title", title),
				),
			},
			{
				ResourceName:      "sotoon_iam_user_public_key.pk",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
