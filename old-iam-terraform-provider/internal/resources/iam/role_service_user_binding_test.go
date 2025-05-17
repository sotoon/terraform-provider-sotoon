package iam_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	uuid "github.com/satori/go.uuid"
)

func Test_RoleServiceUserBinding(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	roleUUID := uuid.NewV4()
	roleUUIDString := roleUUID.String()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()

	serviceUserUUID := uuid.NewV4()
	serviceUserUUIDString := serviceUserUUID.String()

	items := map[string]string{
		"cluster": "neda",
	}

	newItems := map[string]string{
		"cluster": "afra",
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					mocks.IAMClient.EXPECT().BindRoleToServiceUser(&workspaceUUID, &roleUUID, &serviceUserUUID, items).AnyTimes().Return(nil)
					mocks.IAMClient.EXPECT().GetBindedRoleToServiceUserItems(&workspaceUUID, &roleUUID, &serviceUserUUID).AnyTimes().Return(items, nil)
				},
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				resource "sotoon_iam_role_service_user_binding" "test" {
					service_user_id = "%s"
					role_id = "%s"
					workspace_id = "%s"
					items = {
						"cluster" = "neda",
					}
				}`, serviceUserUUID, roleUUIDString, workspaceUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_role_service_user_binding.test", "service_user_id", serviceUserUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_service_user_binding.test", "role_id", roleUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_service_user_binding.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_service_user_binding.test", "items.cluster", "neda"),
				),
			},
			{
				ResourceName:  "sotoon_iam_role_service_user_binding.test",
				ImportState:   true,
				ImportStateId: fmt.Sprintf("%s:%s:%s", roleUUIDString, serviceUserUUIDString, workspaceUUIDString),
			},
			{
				ExpectNonEmptyPlan: true,
				PreConfig: func() {
					mocks.IAMClient.EXPECT().BindRoleToServiceUser(&workspaceUUID, &roleUUID, &serviceUserUUID, gomock.Any()).AnyTimes().Return(nil)
					mocks.IAMClient.EXPECT().GetBindedRoleToServiceUserItems(&workspaceUUID, &roleUUID, &serviceUserUUID).AnyTimes().Return(newItems, nil)
					mocks.IAMClient.EXPECT().UnbindRoleFromServiceUser(&workspaceUUID, &roleUUID, &serviceUserUUID, gomock.Any()).AnyTimes().Return(nil)
				},
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				resource "sotoon_iam_role_service_user_binding" "test" {
					service_user_id = "%s"
					role_id = "%s"
					workspace_id = "%s"
					items = {
						"cluster" = "afra",
					}
				}`, serviceUserUUID, roleUUIDString, workspaceUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_role_service_user_binding.test", "service_user_id", serviceUserUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_service_user_binding.test", "role_id", roleUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_service_user_binding.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_service_user_binding.test", "items.cluster", "afra"),
				),
			},
		},
	})
}
