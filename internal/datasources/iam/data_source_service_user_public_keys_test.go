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
	pubKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7kG7aZb1W9n3mX6Lw0N2Q7pPq3yZV7g5r8jL9z8lQmA9oJQqD4lX9nB5u2jC4t4h1uWq1Ff5rjM4cO2q9f3B8nT9pA7vC2hX info@mail.com"
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
