package iam_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	uuid "github.com/satori/go.uuid"
)

func Test_RuleResource(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()

	ruleUUID := uuid.NewV4()
	ruleUUIDString := ruleUUID.String()

	name := "test"
	actions := []string{"GET"}
	path := "/compute.cafebazaar.cloud/instancetype*"
	service := "k8s"
	object := "rri:v1:cafebazaar.cloud:" + workspaceUUIDString + ":" + service + ":" + path

	secondName := "test2"
	secondActions := []string{"GET", "POST"}
	secondPath := "/compute.cafebazaar.cloud/a-change/instancetype*"
	secondService := "k8s"
	secondObject := "rri:v1:cafebazaar.cloud:" + workspaceUUIDString + ":" + secondService + ":" + secondPath

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					mocks.IAMClient.EXPECT().CreateRule(name, &workspaceUUID, actions, object, false).AnyTimes().Return(&types.Rule{
						UUID:          &ruleUUID,
						WorkspaceUUID: &workspaceUUID,
						Name:          name,
						Actions:       actions,
						Object:        object,
						Deny:          false,
					}, nil)

					mocks.IAMClient.EXPECT().GetRule(&ruleUUID, &workspaceUUID).AnyTimes().Return(&types.Rule{
						UUID:          &ruleUUID,
						WorkspaceUUID: &workspaceUUID,
						Name:          name,
						Actions:       actions,
						Object:        object,
						Deny:          false,
					}, nil)
				},
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				resource "sotoon_iam_rule" "test" {
					name         = "%s"
					actions      = ["GET"]
					service      = "k8s"
					path         = "%s"
					workspace_id = "%s"
					is_denial    = false
				}`, name, path, workspaceUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "name", name),
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "id", ruleUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "path", path),
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "service", service),
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "is_denial", "false"),
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "actions.#", "1"),
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "actions.0", "GET"),
				),
			},
			{
				ResourceName:      "sotoon_iam_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s:%s", ruleUUIDString, workspaceUUIDString),
			},
			{
				ExpectNonEmptyPlan: true,
				PreConfig: func() {
					mocks.IAMClient.EXPECT().GetRule(&ruleUUID, &workspaceUUID).AnyTimes().Return(&types.Rule{
						UUID:          &ruleUUID,
						WorkspaceUUID: &workspaceUUID,
						Name:          secondName,
						Actions:       secondActions,
						Object:        secondObject,
						Deny:          true,
					}, nil)
					mocks.IAMClient.EXPECT().UpdateRule(&ruleUUID, secondName, &workspaceUUID, secondActions, secondObject, true).AnyTimes().Return(&types.Rule{
						UUID:          &ruleUUID,
						WorkspaceUUID: &workspaceUUID,
						Name:          secondName,
						Actions:       secondActions,
						Object:        secondObject,
						Deny:          true,
					}, nil)
					mocks.IAMClient.EXPECT().DeleteRule(&ruleUUID, &workspaceUUID).AnyTimes().Return(nil)
				},
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				resource "sotoon_iam_rule" "test" {
					name         = "%s"
					actions      = ["GET", "POST"]
					service      = "k8s"
					path         = "%s"
					workspace_id = "%s"
					is_denial    = true
				}`, secondName, secondPath, workspaceUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "name", secondName),
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "id", ruleUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "path", secondPath),
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "service", secondService),
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "is_denial", "true"),
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "actions.#", "2"),
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "actions.0", "GET"),
					resource.TestCheckResourceAttr("sotoon_iam_rule.test", "actions.1", "POST"),
				),
			},
		},
	})
}
