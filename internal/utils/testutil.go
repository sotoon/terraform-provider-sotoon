//go:build acceptance
// +build acceptance

package testutil

import (
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

// TestAccPreCheck ensures required environment variables are set before running acceptance tests
func TestAccPreCheck(t *testing.T) {
	req := []string{"SOTOON_API_TOKEN", "SOTOON_WORKSPACE_ID", "SOTOON_API_HOST", "SOTOON_USER_ID"}
	var missing []string
	for _, k := range req {
		if os.Getenv(k) == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		t.Fatalf("Missing required env vars for acceptance tests: %v", missing)
	}
}

// BaseProviderBlock returns the provider block without embedding secrets
// The provider should read its configuration from environment variables instead.
func BaseProviderBlock() string {
	return `
provider "sotoon" {}
`
}

// RandEmail generates a unique email address for testing
func RandEmail() string {
	return "tf-acc-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) + "@example.test"
}

// RandName generates a unique resource name with the given prefix
func RandName(prefix string) string {
	return prefix + "-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
}
