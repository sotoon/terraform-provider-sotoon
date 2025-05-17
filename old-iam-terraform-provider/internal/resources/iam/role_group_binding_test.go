package iam_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	uuid "github.com/satori/go.uuid"
)

func Test_RoleGroupResource(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()

	groupUUID := uuid.NewV4()
	groupUUIDString := groupUUID.String()

	roleUUID := uuid.NewV4()
	roleUUIDString := roleUUID.String()

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
					mocks.IAMClient.EXPECT().BindRoleToGroup(&workspaceUUID, &roleUUID, &groupUUID, items).Times(1).Return(nil)
					mocks.IAMClient.EXPECT().GetBindedRoleToGroupItems(&workspaceUUID, &roleUUID, &groupUUID).AnyTimes().Return(items, nil)
				},
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				resource "sotoon_iam_role_group_binding" "test" {
					workspace_id = "%s"
					role_id = "%s"
					group_id = "%s"
					items = {
						"cluster" = "neda",
					}
				}`, workspaceUUIDString, roleUUIDString, groupUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_role_group_binding.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_group_binding.test", "role_id", roleUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_group_binding.test", "group_id", groupUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_group_binding.test", "items.cluster", "neda"),
				),
			},
			{
				ResourceName:  "sotoon_iam_role_group_binding.test",
				ImportState:   true,
				ImportStateId: fmt.Sprintf("%s:%s:%s", roleUUIDString, groupUUIDString, workspaceUUIDString),
			},
			{
				ExpectNonEmptyPlan: true,
				PreConfig: func() {
					mocks.IAMClient.EXPECT().BindRoleToGroup(&workspaceUUID, &roleUUID, &groupUUID, newItems).AnyTimes().Return(nil)
					mocks.IAMClient.EXPECT().GetBindedRoleToGroupItems(&workspaceUUID, &roleUUID, &groupUUID).AnyTimes().Return(newItems, nil)
					mocks.IAMClient.EXPECT().UnbindRoleFromGroup(&workspaceUUID, &roleUUID, &groupUUID, gomock.Any()).AnyTimes().Return(nil)
				},
				Config: provider.TestProviderConfig + fmt.Sprintf(`
					resource "sotoon_iam_role_group_binding" "test" {
						workspace_id = "%s"
						role_id = "%s"
						group_id = "%s"
						items = {
							"cluster" = "afra",
						}
					}`, workspaceUUIDString, roleUUIDString, groupUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_role_group_binding.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_group_binding.test", "role_id", roleUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_group_binding.test", "group_id", groupUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_group_binding.test", "items.cluster", "afra"),
				),
			},
		},
	})
}
