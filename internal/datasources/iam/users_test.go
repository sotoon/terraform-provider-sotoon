package iam_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	uuid "github.com/satori/go.uuid"
)

func Test_UsersDataSourceJustWorkspace(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()

	user1UUID := uuid.NewV4()
	user2UUID := uuid.NewV4()
	user3UUID := uuid.NewV4()

	user1UUIDString := user1UUID.String()

	user1Email := "test1@sotoon.ir"
	user2Email := "test2@sotoon.ir"
	user3Email := "test3@sotoon.ir"

	user1Name := "test user 1"
	user2Name := "test user 2"
	user3Name := "test user 3"

	mocks.IAMClient.EXPECT().GetWorkspaceUsers(&workspaceUUID).Return([]*types.User{
		{
			UUID:            &user1UUID,
			Name:            user1Name,
			Email:           user1Email,
			InvitationToken: "abcdefghijklmnopqrstuvwxyz0123456789",
		},
		{
			UUID:            &user2UUID,
			Name:            user2Name,
			Email:           user2Email,
			InvitationToken: "abcdefghijklmnopqrstuvwxyz0123456789",
		},
		{
			UUID:            &user3UUID,
			Name:            user3Name,
			Email:           user3Email,
			InvitationToken: "abcdefghijklmnopqrstuvwxyz0123456789",
		},
	}, nil).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				data "sotoon_iam_users" "test" {
					workspace_id = "%s"
				}`, workspaceUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("data.sotoon_iam_users.test", "group_id"),
					resource.TestCheckResourceAttr("data.sotoon_iam_users.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_users.test", "users.#", "3"),
					resource.TestCheckResourceAttr("data.sotoon_iam_users.test", "users.0.id", user1UUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_users.test", "users.0.workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_users.test", "users.0.name", user1Name),
					resource.TestCheckResourceAttr("data.sotoon_iam_users.test", "users.0.email", user1Email),
				),
			},
		},
	})
}

func Test_UsersDataSourceWithGroup(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()

	groupUUID := uuid.NewV4()
	groupUUIDString := groupUUID.String()

	user1UUID := uuid.NewV4()
	user2UUID := uuid.NewV4()
	user3UUID := uuid.NewV4()

	user1UUIDString := user1UUID.String()

	user1Email := "test1@sotoon.ir"
	user2Email := "test2@sotoon.ir"
	user3Email := "test3@sotoon.ir"

	user1Name := "test user 1"
	user2Name := "test user 2"
	user3Name := "test user 3"

	mocks.IAMClient.EXPECT().GetAllGroupUsers(&workspaceUUID, &groupUUID).Return([]*types.User{
		{
			UUID:            &user1UUID,
			Name:            user1Name,
			Email:           user1Email,
			InvitationToken: "abcdefghijklmnopqrstuvwxyz0123456789",
		},
		{
			UUID:            &user2UUID,
			Name:            user2Name,
			Email:           user2Email,
			InvitationToken: "abcdefghijklmnopqrstuvwxyz0123456789",
		},
		{
			UUID:            &user3UUID,
			Name:            user3Name,
			Email:           user3Email,
			InvitationToken: "abcdefghijklmnopqrstuvwxyz0123456789",
		},
	}, nil).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				data "sotoon_iam_users" "test" {
					workspace_id = "%s"
					group_id = "%s"
				}`, workspaceUUIDString, groupUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sotoon_iam_users.test", "group_id", groupUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_users.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_users.test", "users.#", "3"),
					resource.TestCheckResourceAttr("data.sotoon_iam_users.test", "users.0.id", user1UUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_users.test", "users.0.workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_users.test", "users.0.name", user1Name),
					resource.TestCheckResourceAttr("data.sotoon_iam_users.test", "users.0.email", user1Email),
				),
			},
		},
	})
}
