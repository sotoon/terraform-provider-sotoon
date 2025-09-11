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

func TestAccUserGroupMembership_basic(t *testing.T) {
	userUUID := os.Getenv("SOTOON_TEST_USER_ID")
	if userUUID == "" {
		t.Skip("SOTOON_TEST_USER_ID not set; skipping to avoid invites")
	}

	groupName := testutil.RandName("tf-acc-group")

	cfg := fmt.Sprintf(`
%s

resource "sotoon_iam_group" "g" {
  name        = "%s"
  description = "acceptance group for user_group_membership"
}

resource "sotoon_iam_user_group_membership" "m" {
  group_id = sotoon_iam_group.g.id
  user_ids = ["%s"]
}
`, testutil.BaseProviderBlock(), groupName, userUUID)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		ProviderFactories: testutil.ProviderFactories,
		Steps:             []resource.TestStep{{Config: cfg}},
	})
}
