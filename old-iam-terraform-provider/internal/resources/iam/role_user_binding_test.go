package iam_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	uuid "github.com/satori/go.uuid"
)

func Test_RoleRuleBindingResourceWithItems(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	roleUUID := uuid.NewV4()
	roleUUIDString := roleUUID.String()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()

	userUUID := uuid.NewV4()
	userUUIDString := userUUID.String()

	items := map[string]string{
		"cluster": "neda",
		"zone":    "test",
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
					mocks.IAMClient.EXPECT().BindRoleToUser(&workspaceUUID, &roleUUID, &userUUID, items).AnyTimes().Return(nil)
					mocks.IAMClient.EXPECT().GetBindedRoleToUserItems(&workspaceUUID, &roleUUID, &userUUID).AnyTimes().Return(items, nil)
				},
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				resource "sotoon_iam_role_user_binding" "test" {
					user_id = "%s"
					role_id = "%s"
					workspace_id = "%s"
					items = {
						"cluster" = "neda",
						"zone" = "test"
					}
				}`, userUUIDString, roleUUIDString, workspaceUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "user_id", userUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "role_id", roleUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "items.cluster", "neda"),
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "items.zone", "test"),
				),
			},
			{
				ResourceName:  "sotoon_iam_role_user_binding.test",
				ImportState:   true,
				ImportStateId: fmt.Sprintf("%s:%s:%s", roleUUIDString, userUUIDString, workspaceUUIDString),
			},
			{
				ExpectNonEmptyPlan: true,
				PreConfig: func() {
					mocks.IAMClient.EXPECT().BindRoleToUser(&workspaceUUID, &roleUUID, &userUUID, gomock.Any()).AnyTimes().Return(nil)
					mocks.IAMClient.EXPECT().GetBindedRoleToUserItems(&workspaceUUID, &roleUUID, &userUUID).AnyTimes().Return(newItems, nil)
					mocks.IAMClient.EXPECT().UnbindRoleFromUser(&workspaceUUID, &roleUUID, &userUUID, gomock.Any()).AnyTimes().Return(nil)
				},
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				resource "sotoon_iam_role_user_binding" "test" {
					user_id = "%s"
					role_id = "%s"
					workspace_id = "%s"
					items = {
						"cluster" = "afra",
					}
				}`, userUUIDString, roleUUIDString, workspaceUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "user_id", userUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "role_id", roleUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "items.cluster", "afra"),
					resource.TestCheckNoResourceAttr("sotoon_iam_role_user_binding.test", "items.zone"),
				),
			},
		},
	})
}

func Test_RoleRuleBindingResourceChangeUser(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	roleUUID := uuid.NewV4()
	roleUUIDString := roleUUID.String()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()

	userUUID1 := uuid.NewV4()
	userUUIDString1 := userUUID1.String()

	userUUID2 := uuid.NewV4()
	userUUIDString2 := userUUID2.String()

	items := map[string]string{
		"cluster": "neda",
		"zone":    "test",
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					mocks.IAMClient.EXPECT().BindRoleToUser(&workspaceUUID, &roleUUID, &userUUID1, items).AnyTimes().Return(nil)
					mocks.IAMClient.EXPECT().GetBindedRoleToUserItems(&workspaceUUID, &roleUUID, &userUUID1).AnyTimes().Return(items, nil)
				},
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				resource "sotoon_iam_role_user_binding" "test" {
					user_id = "%s"
					role_id = "%s"
					workspace_id = "%s"
					items = {
						"cluster" = "neda",
						"zone" = "test"
					}
				}`, userUUIDString1, roleUUIDString, workspaceUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "user_id", userUUIDString1),
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "role_id", roleUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "items.cluster", "neda"),
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "items.zone", "test"),
				),
			},
			{
				ResourceName:  "sotoon_iam_role_user_binding.test",
				ImportState:   true,
				ImportStateId: fmt.Sprintf("%s:%s:%s", roleUUIDString, userUUIDString1, workspaceUUIDString),
			},
			{
				ExpectNonEmptyPlan: true,
				PreConfig: func() {
					mocks.IAMClient.EXPECT().BindRoleToUser(&workspaceUUID, &roleUUID, gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
					mocks.IAMClient.EXPECT().GetBindedRoleToUserItems(&workspaceUUID, &roleUUID, gomock.Any()).AnyTimes().Return(items, nil)
					mocks.IAMClient.EXPECT().UnbindRoleFromUser(&workspaceUUID, &roleUUID, gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
				},
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				resource "sotoon_iam_role_user_binding" "test" {
					user_id = "%s"
					role_id = "%s"
					workspace_id = "%s"
					items = {
						"cluster" = "afra",
					}
				}`, userUUIDString2, roleUUIDString, workspaceUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "user_id", userUUIDString2),
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "role_id", roleUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role_user_binding.test", "items.cluster", "afra"),
					resource.TestCheckNoResourceAttr("sotoon_iam_role_user_binding.test", "items.zone"),
				),
			},
		},
	})
}
