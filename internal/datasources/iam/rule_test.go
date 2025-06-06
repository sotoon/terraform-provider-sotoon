package iam_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	uuid "github.com/satori/go.uuid"
)

func Test_RuleDataSourceWithSpecifiedWokspace(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	organizationUUID := uuid.NewV4()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()
	workspaceName := "test_workspace"

	ruleUUID := uuid.NewV4()
	ruleUUIDString := ruleUUID.String()
	ruleName := "test_rule"
	actions := []string{"GET"}
	path := "/compute.cafebazaar.cloud/instancetype*"
	service := "k8s"
	object := "rri:v1:cafebazaar.cloud:" + workspaceUUIDString + ":" + service + ":" + path

	mocks.IAMClient.EXPECT().GetWorkspace(&workspaceUUID).Return(&types.Workspace{
		UUID:         &workspaceUUID,
		Name:         workspaceName,
		Namespace:    workspaceName,
		Organization: &organizationUUID,
	}, nil).AnyTimes()

	mocks.IAMClient.EXPECT().GetRuleByName(ruleName, workspaceName).Return(&types.Rule{
		UUID:          &ruleUUID,
		WorkspaceUUID: &workspaceUUID,
		Name:          ruleName,
		Actions:       actions,
		Object:        object,
		Deny:          false,
	}, nil).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				data "sotoon_iam_rule" "test" {
					workspace_id = "%s"
					name = "%s"
				}`, workspaceUUIDString, ruleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "name", ruleName),
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "id", ruleUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "path", path),
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "service", service),
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "is_denial", "false"),
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "actions.#", "1"),
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "actions.0", "GET"),
				),
			},
		},
	})
}

func Test_RuleDataSourceWithUnspecifiedWokspace(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	workspaceUUIDString := utils.GLOBAL_WORKSPACE
	workspaceUUID, _ := uuid.FromString(workspaceUUIDString)
	workspaceName := "global"

	ruleUUID := uuid.NewV4()
	ruleUUIDString := ruleUUID.String()
	ruleName := "can-manage-kaas-engine-clusters"
	actions := []string{"*"}
	path := "/{datacenter}/infrastructure.cluster.x-k8s.io/v+/sotoonclusters*"
	service := "kaas-engine"
	object := "rri:v1:cafebazaar.cloud:" + workspaceUUIDString + ":" + service + ":" + path

	mocks.IAMClient.EXPECT().GetRuleByName(ruleName, workspaceName).Return(&types.Rule{
		UUID:          &ruleUUID,
		WorkspaceUUID: &workspaceUUID,
		Name:          ruleName,
		Actions:       actions,
		Object:        object,
		Deny:          false,
	}, nil).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				data "sotoon_iam_rule" "test" {
					name = "%s"
				}`, ruleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "name", ruleName),
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "id", ruleUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "path", path),
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "service", service),
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "is_denial", "false"),
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "actions.#", "1"),
					resource.TestCheckResourceAttr("data.sotoon_iam_rule.test", "actions.0", "*"),
				),
			},
		},
	})
}
