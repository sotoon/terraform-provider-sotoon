package iam_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/utils"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	uuid "github.com/satori/go.uuid"
)

func Test_RoleResource(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	roleUUID := uuid.NewV4()
	roleUUIDString := roleUUID.String()
	roleName := "test_role"
	roleNameAfterEdit := "edited_test_role"

	organizationUUID := uuid.NewV4()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()
	workspaceName := "test_wroskapce_name"

	rule1UUID := uuid.NewV4()
	rule1UUIDString := rule1UUID.String()
	rule1Name := "rule_1_name"

	rule2UUID := uuid.NewV4()
	rule2UUIDString := rule2UUID.String()
	rule2Name := "rule_2_name"

	rule3UUID := uuid.NewV4()
	rule3UUIDString := rule2UUID.String()
	rule3Name := "rule_3_name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					mocks.IAMClient.EXPECT().CreateRole(roleName, &workspaceUUID).AnyTimes().Return(&types.Role{
						UUID:      &roleUUID,
						Workspace: &workspaceUUID,
						Name:      roleName,
					}, nil)

					mocks.IAMClient.EXPECT().BindRuleToRole(&roleUUID, gomock.Any(), &workspaceUUID).AnyTimes().Return(nil)
					mocks.IAMClient.EXPECT().GetRole(&roleUUID, &workspaceUUID).AnyTimes().Return(&types.RoleRes{
						UUID: &roleUUID,
						Name: roleName,
						Workspace: &types.Workspace{
							UUID:         &workspaceUUID,
							Name:         workspaceName,
							Namespace:    workspaceName,
							Organization: &organizationUUID,
						},
					}, nil)
					mocks.IAMClient.EXPECT().GetRoleRules(&roleUUID, &workspaceUUID).AnyTimes().Return([]*types.Rule{
						{
							UUID:          &rule1UUID,
							WorkspaceUUID: &workspaceUUID,
							Name:          rule1Name,
							Actions:       []string{},
							Object:        utils.CreateRRI(workspaceUUIDString, "service", "/path1"),
							Deny:          false,
						},
						{
							UUID:          &rule2UUID,
							WorkspaceUUID: &workspaceUUID,
							Name:          rule2Name,
							Actions:       []string{},
							Object:        utils.CreateRRI(workspaceUUIDString, "service", "/path2"),
							Deny:          false,
						},
					}, nil)
				},
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				resource "sotoon_iam_role" "test" {
					name = "%s"
					workspace_id = "%s"
					rules = [
						{id = "%s"},
						{id = "%s"},
					]
				}`, roleName, workspaceUUIDString, rule1UUIDString, rule2UUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_role.test", "name", roleName),
					resource.TestCheckResourceAttr("sotoon_iam_role.test", "id", roleUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role.test", "rules.#", "2"),
				),
			},
			{
				ResourceName:      "sotoon_iam_role.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s:%s", roleUUIDString, workspaceUUIDString),
			},
			{
				ExpectNonEmptyPlan: true,
				PreConfig: func() {
					mocks.IAMClient.EXPECT().BindRuleToRole(&roleUUID, &rule3UUID, &workspaceUUID).AnyTimes().Return(nil)
					mocks.IAMClient.EXPECT().UnbindRuleFromRole(&roleUUID, &rule2UUID, &workspaceUUID).AnyTimes().Return(nil)
					mocks.IAMClient.EXPECT().UpdateRole(&roleUUID, roleNameAfterEdit, &workspaceUUID).AnyTimes().Return(&types.Role{
						UUID:      &roleUUID,
						Workspace: &workspaceUUID,
						Name:      roleNameAfterEdit,
					}, nil)

					mocks.IAMClient.EXPECT().DeleteRole(&roleUUID, &workspaceUUID).AnyTimes().Return(nil)

					mocks.IAMClient.EXPECT().GetRole(&roleUUID, &workspaceUUID).AnyTimes().Return(&types.RoleRes{
						UUID: &roleUUID,
						Name: roleName,
						Workspace: &types.Workspace{
							UUID:         &workspaceUUID,
							Name:         workspaceName,
							Namespace:    workspaceName,
							Organization: &organizationUUID,
						},
					}, nil)
					mocks.IAMClient.EXPECT().GetRoleRules(&roleUUID, &workspaceUUID).AnyTimes().Return([]*types.Rule{
						{
							UUID:          &rule1UUID,
							WorkspaceUUID: &workspaceUUID,
							Name:          rule1Name,
							Actions:       []string{},
							Object:        utils.CreateRRI(workspaceUUIDString, "service", "/path1"),
							Deny:          false,
						},
						{
							UUID:          &rule3UUID,
							WorkspaceUUID: &workspaceUUID,
							Name:          rule3Name,
							Actions:       []string{},
							Object:        utils.CreateRRI(workspaceUUIDString, "service", "/path3"),
							Deny:          false,
						},
					}, nil)
				},
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				resource "sotoon_iam_role" "test" {
					name = "%s"
					workspace_id = "%s"
					rules = [
						{id = "%s"},
						{id = "%s"},
					]
				}`, roleNameAfterEdit, workspaceUUIDString, rule1UUIDString, rule3UUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_role.test", "name", roleNameAfterEdit),
					resource.TestCheckResourceAttr("sotoon_iam_role.test", "id", roleUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_role.test", "rules.#", "2"),
				),
			},
		},
	})
}
