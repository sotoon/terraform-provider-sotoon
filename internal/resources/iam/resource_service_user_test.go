//go:build acceptance
// +build acceptance

package iamresources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testutil "github.com/sotoon/terraform-provider-sotoon/internal/utils"
)

func TestAccServiceUser_basic(t *testing.T) {
	name := testutil.RandName("tf-acc-su")

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_service_user" "test" {
  name        = "%s"
  description = "created by acceptance test"
}
`, testutil.BaseProviderBlock(), name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		ProviderFactories: testutil.ProviderFactories,
		Steps:             []resource.TestStep{{Config: cfg}},
	})
}
