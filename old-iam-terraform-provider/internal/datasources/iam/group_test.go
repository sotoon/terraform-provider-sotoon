package iam_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	uuid "github.com/satori/go.uuid"
)

func Test_GroupDataSource(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	organizationUUID := uuid.NewV4()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()
	workspaceName := "test_workspace"

	groupUUID := uuid.NewV4()
	groupUUIDString := groupUUID.String()
	groupName := "test_group"

	mocks.IAMClient.EXPECT().GetWorkspace(&workspaceUUID).Return(&types.Workspace{
		UUID:         &workspaceUUID,
		Name:         workspaceName,
		Namespace:    workspaceName,
		Organization: &organizationUUID,
	}, nil).AnyTimes()

	mocks.IAMClient.EXPECT().GetGroupByName(workspaceName, groupName).Return(&types.Group{
		UUID: &groupUUID,
		Name: groupName,
		Workspace: types.Workspace{
			UUID:         &workspaceUUID,
			Name:         workspaceName,
			Namespace:    workspaceName,
			Organization: &organizationUUID,
		},
	}, nil).AnyTimes()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				data "sotoon_iam_group" "test" {
					workspace_id = "%s"
					name = "%s"
				}`, workspaceUUIDString, groupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sotoon_iam_group.test", "name", groupName),
					resource.TestCheckResourceAttr("data.sotoon_iam_group.test", "id", groupUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_group.test", "workspace_id", workspaceUUIDString),
				),
			},
		},
	})
}
