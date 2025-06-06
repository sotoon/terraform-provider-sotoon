package iam_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	uuid "github.com/satori/go.uuid"
)

func Test_ServiceUserDataSource(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	organizationUUID := uuid.NewV4()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()
	workspaceName := "test_workspace"

	serviceUserUUID := uuid.NewV4()
	serviceUserUUIDString := serviceUserUUID.String()
	serviceUserName := "test service user"

	mocks.IAMClient.EXPECT().GetWorkspace(&workspaceUUID).Return(&types.Workspace{
		UUID:         &workspaceUUID,
		Name:         workspaceName,
		Namespace:    workspaceName,
		Organization: &organizationUUID,
	}, nil).AnyTimes()

	mocks.IAMClient.EXPECT().GetServiceUserByName(workspaceName, serviceUserName).Return(&types.ServiceUser{
		UUID:      &serviceUserUUID,
		Name:      serviceUserName,
		Workspace: workspaceUUIDString,
	}, nil).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				data "sotoon_iam_service_user" "test" {
					workspace_id = "%s"
					name = "%s"
				}`, workspaceUUIDString, serviceUserName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sotoon_iam_service_user.test", "id", serviceUserUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_service_user.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_service_user.test", "name", serviceUserName),
				),
			},
		},
	})
}
