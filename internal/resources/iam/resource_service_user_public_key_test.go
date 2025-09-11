//go:build acceptance
// +build acceptance

package iamresources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testutil "github.com/sotoon/terraform-provider-sotoon/internal/utils"
)

func TestAccServiceUserPublicKey_basic(t *testing.T) {
	suName := testutil.RandName("tf-acc-su")
	pubKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7kG7aZb1W9n3mX6Lw0N2Q7pPq3yZV7g5r8jL9z8lQmA9oJQqD4lX9nB5u2jC4t4h1uWq1Ff5rjM4cO2q9f3B8nT9pA7vC2hX info@mail.com"

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
`, testutil.BaseProviderBlock(), suName, pubKey)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		ProviderFactories: testutil.ProviderFactories,
		Steps:             []resource.TestStep{{Config: cfg, Check: resource.TestCheckOutput("su_pk_title", "tf-acc-builder-key")}},
	})
}
