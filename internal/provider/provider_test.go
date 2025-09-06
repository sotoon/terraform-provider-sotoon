//go:build acceptance
// +build acceptance

package provider_test

import (
	"os"
	"testing"

	sotoon "github.com/sotoon/terraform-provider-sotoon/internal/provider"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

var ProviderFactories = map[string]func() (*schema.Provider, error){
	"sotoon": func() (*schema.Provider, error) { return sotoon.Provider(), nil },
}

func testAccPreCheck(t *testing.T) {
	req := []string{"SOTOON_API_TOKEN", "SOTOON_WORKSPACE_ID", "SOTOON_API_HOST", "SOTOON_USER_ID"}
	missing := []string{}
	for _, k := range req {
		if os.Getenv(k) == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		t.Fatalf("Missing required env vars for acceptance tests: %v", missing)
	}
}

func baseProviderBlock() string {
	return fmt.Sprintf(`
provider "sotoon" {
  api_token    = "%s"
  workspace_id = "%s"
  api_host     = "%s"
  user_id      = "%s"
}
`, os.Getenv("SOTOON_API_TOKEN"),
		os.Getenv("SOTOON_WORKSPACE_ID"),
		os.Getenv("SOTOON_API_HOST"),
		os.Getenv("SOTOON_USER_ID"),
	)
}

// Unique helpers
func randEmail() string {
	return "tf-acc-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) + "@example.test"
}

func randName(prefix string) string {
	return prefix + "-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
}
