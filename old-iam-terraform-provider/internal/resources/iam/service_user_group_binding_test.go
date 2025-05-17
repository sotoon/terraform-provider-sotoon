package iam_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	uuid "github.com/satori/go.uuid"
)

func Test_GroupServiceUserResource(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()

	groupUUID := uuid.NewV4()
	groupUUIDString := groupUUID.String()

	serviceUserUUID := uuid.NewV4()
	serviceUserUUIDString := serviceUserUUID.String()
	serviceuser := types.ServiceUser{
		UUID:      &serviceUserUUID,
		Name:      "test_user",
		Workspace: workspaceUUIDString,
	}

	mocks.IAMClient.EXPECT().BindServiceUserToGroup(&workspaceUUID, &groupUUID, &serviceUserUUID).AnyTimes().Return(nil)
	mocks.IAMClient.EXPECT().GetGroupServiceUser(&workspaceUUID, &groupUUID, &serviceUserUUID).AnyTimes().Return(&serviceuser, nil)
	mocks.IAMClient.EXPECT().UnbindServiceUserFromGroup(&workspaceUUID, &groupUUID, &serviceUserUUID).AnyTimes().Return(nil)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				resource "sotoon_iam_service_user_group_binding" "test" {
					workspace_id = "%s"
					group_id = "%s"
					user_id = "%s"
				}`, workspaceUUIDString, groupUUIDString, serviceUserUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_service_user_group_binding.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_service_user_group_binding.test", "group_id", groupUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_service_user_group_binding.test", "user_id", serviceUserUUIDString),
				),
			},
			{
				ResourceName:  "sotoon_iam_service_user_group_binding.test",
				ImportState:   true,
				ImportStateId: fmt.Sprintf("%s:%s:%s", serviceUserUUIDString, groupUUIDString, workspaceUUIDString),
			},
		},
	})
}
