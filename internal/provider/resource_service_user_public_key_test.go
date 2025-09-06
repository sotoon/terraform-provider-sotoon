//go:build acceptance
// +build acceptance

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceUserPublicKey_basic(t *testing.T) {
	suName := randName("tf-acc-su")
	pubKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJ4rAgexMXRfB+ZzQpzKydKBTKSE5PTxvfVBa4J8uFLC alireza@Alirezas-MacBook-Air.local"

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_service_user" "su" {
  name        = "%s"
  description = "acceptance service user"
}

resource "sotoon_iam_service_user_public_key" "pk" {
  service_user_id = sotoon_iam_service_user.su.id
  title           = "tf-acc-builder-key"
  public_key      = "%s"

  lifecycle {
    ignore_changes = [public_key]
  }
}

output "su_pk_title" { value = sotoon_iam_service_user_public_key.pk.title }
`, baseProviderBlock(), suName, pubKey)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps:             []resource.TestStep{{Config: cfg, Check: resource.TestCheckOutput("su_pk_title", "tf-acc-builder-key")}},
	})
}
