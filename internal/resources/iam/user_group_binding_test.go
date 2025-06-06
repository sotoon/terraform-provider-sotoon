package iam_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	uuid "github.com/satori/go.uuid"
)

func Test_GroupUserResource(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	organizationUUID := uuid.NewV4()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()
	workspace := types.Workspace{
		UUID:         &workspaceUUID,
		Name:         "test_workspace",
		Namespace:    "test_workspace",
		Organization: &organizationUUID,
	}

	groupUUID := uuid.NewV4()
	groupUUIDString := groupUUID.String()
	group := types.Group{
		UUID:      &groupUUID,
		Name:      "test_group",
		Workspace: workspace,
	}

	userUUID := uuid.NewV4()
	userUUIDString := userUUID.String()
	user := types.User{
		UUID:  &userUUID,
		Name:  "test_user",
		Email: "test@sotoon.ir",
	}

	mocks.IAMClient.EXPECT().GetGroup(&workspaceUUID, &groupUUID).AnyTimes().Return(&group, nil)
	mocks.IAMClient.EXPECT().BindGroup(group.Name, &workspaceUUID, &groupUUID, &userUUID).AnyTimes().Return(nil)
	mocks.IAMClient.EXPECT().GetGroupUser(&workspaceUUID, &groupUUID, &userUUID).AnyTimes().Return(&user, nil)
	mocks.IAMClient.EXPECT().UnbindUserFromGroup(&workspaceUUID, &groupUUID, &userUUID).AnyTimes().Return(nil)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				resource "sotoon_iam_user_group_binding" "test" {
					workspace_id = "%s"
					group_id = "%s"
					user_id = "%s"
				}`, workspaceUUIDString, groupUUIDString, userUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_user_group_binding.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_user_group_binding.test", "group_id", groupUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_user_group_binding.test", "user_id", userUUIDString),
				),
			},
			{
				ResourceName:  "sotoon_iam_user_group_binding.test",
				ImportState:   true,
				ImportStateId: fmt.Sprintf("%s:%s:%s", userUUIDString, groupUUIDString, workspaceUUIDString),
			},
		},
	})
}
