//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserPublicKey_basic(t *testing.T) {
	title := "tf-acc-user-pk"
	pubKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJ4rAgexMXRfB+ZzQpzKydKBTKSE5PTxvfVBa4J8uFLC alireza@Alirezas-MacBook-Air.local"

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_user_public_key" "pk" {
  title      = "%s"
  public_key = "%s"
}

output "user_pk_id"    { value = sotoon_iam_user_public_key.pk.id }
output "user_pk_title" { value = sotoon_iam_user_public_key.pk.title }
`, baseProviderBlock(), title, pubKey)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
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
