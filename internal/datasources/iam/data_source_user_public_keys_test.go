//go:build acceptance
// +build acceptance

package iamdatasources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testutil "github.com/sotoon/terraform-provider-sotoon/internal/utils"
)

func TestAccDataSourceUserPublicKeys_basic(t *testing.T) {
	title := "tf-acc-user-pk"
	pubKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJ4rAgexMXRfB+ZzQpzKydKBTKSE5PTxvfVBa4J8uFLC alireza@Alirezas-MacBook-Air.local"

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_user_public_key" "k" {
  title      = "%s"
  public_key = "%s"
}

data "sotoon_iam_user_public_keys" "me" {
  depends_on = [sotoon_iam_user_public_key.k]
}

output "has_title" {
  value = contains(data.sotoon_iam_user_public_keys.me.public_keys[*].title, "%s")
}
`, testutil.BaseProviderBlock(), title, pubKey, title)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		ProviderFactories: testutil.ProviderFactories,
		Steps: []resource.TestStep{
			{Config: cfg, Check: resource.TestCheckOutput("has_title", "true")},
		},
	})
}
