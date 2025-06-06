package iam_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	uuid "github.com/satori/go.uuid"
)

func Test_GroupResource(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	organizationUUID := uuid.NewV4()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()

	groupUUID := uuid.NewV4()
	groupUUIDString := groupUUID.String()
	groupName := "test"

	mocks.IAMClient.EXPECT().CreateGroup(groupName, &workspaceUUID).AnyTimes().Return(&types.GroupRes{
		UUID:          &groupUUID,
		Name:          groupName,
		WorkspaceUUID: workspaceUUIDString,
		Descriotion:   "",
	}, nil)

	workspace := types.Workspace{
		UUID:         &workspaceUUID,
		Name:         "test_workspace",
		Namespace:    "test_workspace",
		Organization: &organizationUUID,
	}

	mocks.IAMClient.EXPECT().GetWorkspace(&workspaceUUID).AnyTimes().Return(&workspace, nil)

	mocks.IAMClient.EXPECT().GetGroupByName("test_workspace", groupName).AnyTimes().Return(&types.Group{
		UUID:      &groupUUID,
		Name:      groupName,
		Workspace: workspace,
	}, nil)

	mocks.IAMClient.EXPECT().DeleteGroup(&workspaceUUID, &groupUUID).AnyTimes().Return(nil)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				resource "sotoon_iam_group" "test" {
					name = "%s"
					workspace_id = "%s"
				}`, groupName, workspaceUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_group.test", "name", groupName),
					resource.TestCheckResourceAttr("sotoon_iam_group.test", "id", groupUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_group.test", "workspace_id", workspaceUUIDString),
				),
			},
			{
				ResourceName:      "sotoon_iam_group.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s:%s", groupName, workspaceUUIDString),
			},
		},
	})
}
