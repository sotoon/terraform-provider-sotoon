package base_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	uuid "github.com/satori/go.uuid"
)

func Test_WorkspaceDataSource(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()

	organizationUUID := uuid.NewV4()
	organizationUUIDString := organizationUUID.String()

	mocks.IAMClient.EXPECT().GetWorkspace(&workspaceUUID).Return(&types.Workspace{
		UUID:         &workspaceUUID,
		Name:         "test_workspace",
		Namespace:    "test_workspace",
		Organization: &organizationUUID,
	}, nil).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: provider.TestProviderConfig + fmt.Sprintf(`data "sotoon_workspace" "test" {id = "%s"}`, workspaceUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sotoon_workspace.test", "id", workspaceUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_workspace.test", "name", "test_workspace"),
					resource.TestCheckResourceAttr("data.sotoon_workspace.test", "organization_id", organizationUUIDString),
				),
			},
		},
	})
}
