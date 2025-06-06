package provider

import (
	"testing"

	mockiam "git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/mocks/iam"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const TestProviderConfig = `
provider "sotoon" {
	token = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}
`

type MockKeeper struct {
	IAMClient *mockiam.MockClient
}

func NewTestProvider(t *testing.T) (*MockKeeper, *gomock.Controller, func() provider.Provider) {
	mockCtrl := gomock.NewController(t)
	mocks := &MockKeeper{
		IAMClient: mockiam.NewMockClient(mockCtrl),
	}
	return mocks, mockCtrl, func() provider.Provider {
		return &sotoonProvider{
			version: "test",
			mocks:   mocks,
		}
	}
}

func GetTestPrerequisites(t *testing.T) (*MockKeeper, *gomock.Controller, map[string]func() (tfprotov6.ProviderServer, error)) {
	mocks, mockCtrl, provider := NewTestProvider(t)
	return mocks, mockCtrl, map[string]func() (tfprotov6.ProviderServer, error){
		"sotoon": providerserver.NewProtocol6WithError(provider()),
	}
}
