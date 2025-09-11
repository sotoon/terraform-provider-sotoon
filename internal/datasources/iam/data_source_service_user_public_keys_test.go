//go:build acceptance
// +build acceptance

package iamdatasources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testutil "github.com/sotoon/terraform-provider-sotoon/internal/utils"
)

func TestAccDataSourceServiceUserPublicKeys_basic(t *testing.T) {
	suName := testutil.RandName("tf-acc-su")
	pubKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJ4rAgexMXRfB+ZzQpzKydKBTKSE5PTxvfVBa4J8uFLC alireza@Alirezas-MacBook-Air.local"
	title := "tf-acc-builder-key"

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_service_user" "su" {
  name        = "%s"
  description = "acceptance service user"
}

resource "sotoon_iam_service_user_public_key" "k" {
  service_user_id = sotoon_iam_service_user.su.id
  title           = "%s"
  public_key      = "%s"

  lifecycle {
    ignore_changes = [public_key]
  }
}

data "sotoon_iam_service_user_public_keys" "list" {
  service_user_id = sotoon_iam_service_user.su.id
  depends_on      = [sotoon_iam_service_user_public_key.k]
}

output "has_title" {
  value = contains(data.sotoon_iam_service_user_public_keys.list.public_keys[*].title, "%s")
}
`, testutil.BaseProviderBlock(), suName, title, pubKey, title)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		ProviderFactories: testutil.ProviderFactories,
		Steps:             []resource.TestStep{{Config: cfg, Check: resource.TestCheckOutput("has_title", "true")}},
	})
}
