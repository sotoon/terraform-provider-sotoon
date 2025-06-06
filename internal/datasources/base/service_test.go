package base_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Test_ServiceDataSource(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	serviceName := "k8s"

	mocks.IAMClient.EXPECT().GetService(serviceName).Return(&types.Service{
		Name: serviceName,
	}, nil).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: provider.TestProviderConfig + fmt.Sprintf(`data "sotoon_service" "test" {name = "%s"}`, serviceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sotoon_service.test", "id", serviceName),
					resource.TestCheckResourceAttr("data.sotoon_service.test", "name", serviceName),
				),
			},
		},
	})
}
