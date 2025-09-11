//go:build acceptance
// +build acceptance

package testutil

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/sotoon/terraform-provider-sotoon/internal/provider"
)

// ProviderFactories for acceptance tests
var ProviderFactories = map[string]func() (*schema.Provider, error){
	"sotoon": func() (*schema.Provider, error) { return provider.Provider(), nil },
}

func TestAccPreCheck(t *testing.T) {
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

func BaseProviderBlock() string {
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
		os.Getenv("SOTOON_USER_ID"))
}

func RandEmail() string {
	return "tf-acc-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) + "@example.test"
}

func RandName(prefix string) string {
	return prefix + "-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
}
