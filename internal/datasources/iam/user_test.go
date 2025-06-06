package iam_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	uuid "github.com/satori/go.uuid"
)

func Test_UserDataSource(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()

	userUUID := uuid.NewV4()
	userUUIDString := userUUID.String()

	userEmail := "test@sotoon.ir"
	userName := "test user"

	mocks.IAMClient.EXPECT().GetUserByEmail(userEmail, &workspaceUUID).Return(&types.User{
		UUID:            &userUUID,
		Name:            userName,
		Email:           userEmail,
		InvitationToken: "abcdefghijklmnopqrstuvwxyz0123456789",
	}, nil).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				data "sotoon_iam_user" "test" {
					workspace_id = "%s"
					email = "%s"
				}`, workspaceUUIDString, userEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sotoon_iam_user.test", "id", userUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_user.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_user.test", "name", userName),
					resource.TestCheckResourceAttr("data.sotoon_iam_user.test", "email", userEmail),
				),
			},
		},
	})
}
