//go:build acceptance
// +build acceptance

package iamdatasources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testutil "github.com/sotoon/terraform-provider-sotoon/internal/utils"
)

func TestAccDataSourceServiceUserDetails_basic(t *testing.T) {
	name := testutil.RandName("tf-acc-su")

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_service_user" "su" {
  name        = "%s"
  description = "acceptance service user"
}

data "sotoon_iam_service_user_details" "details" {
  service_user_id = sotoon_iam_service_user.su.id
}

output "service_user_name" {
  value = data.sotoon_iam_service_user_details.details.name
}
`, testutil.BaseProviderBlock(), name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		ProviderFactories: testutil.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check:  resource.TestCheckOutput("service_user_name", name),
			},
		},
	})
}
