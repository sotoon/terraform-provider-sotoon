package iam_test

import (
	"fmt"
	"testing"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	uuid "github.com/satori/go.uuid"
)

func Test_ServiceUsersDataSourceJustWorkspace(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()

	serviceUser1UUID := uuid.NewV4()
	serviceUser2UUID := uuid.NewV4()
	serviceUser3UUID := uuid.NewV4()

	serviceUser1UUIDString := serviceUser1UUID.String()

	serviceUser1Name := "test service user 1"
	serviceUser2Name := "test service user 2"
	serviceUser3Name := "test service user 3"

	mocks.IAMClient.EXPECT().GetServiceUsers(&workspaceUUID).Return([]*types.ServiceUser{
		{
			UUID:      &serviceUser1UUID,
			Name:      serviceUser1Name,
			Workspace: workspaceUUIDString,
		},
		{
			UUID:      &serviceUser2UUID,
			Name:      serviceUser2Name,
			Workspace: workspaceUUIDString,
		},
		{
			UUID:      &serviceUser3UUID,
			Name:      serviceUser3Name,
			Workspace: workspaceUUIDString,
		},
	}, nil).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				data "sotoon_iam_service_users" "test" {
					workspace_id = "%s"
				}`, workspaceUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("data.sotoon_iam_service_users.test", "group_id"),
					resource.TestCheckResourceAttr("data.sotoon_iam_service_users.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_service_users.test", "service_users.#", "3"),
					resource.TestCheckResourceAttr("data.sotoon_iam_service_users.test", "service_users.0.id", serviceUser1UUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_service_users.test", "service_users.0.workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_service_users.test", "service_users.0.name", serviceUser1Name),
				),
			},
		},
	})
}

func Test_ServiceUsersDataSourceWithGroup(t *testing.T) {
	mocks, mockCtrl, providerFactory := provider.GetTestPrerequisites(t)
	defer mockCtrl.Finish()

	workspaceUUID := uuid.NewV4()
	workspaceUUIDString := workspaceUUID.String()

	groupUUID := uuid.NewV4()
	groupUUIDString := groupUUID.String()

	serviceUser1UUID := uuid.NewV4()
	serviceUser2UUID := uuid.NewV4()
	serviceUser3UUID := uuid.NewV4()

	serviceUser1UUIDString := serviceUser1UUID.String()

	serviceUser1Name := "test service user 1"
	serviceUser2Name := "test service user 2"
	serviceUser3Name := "test service user 3"

	mocks.IAMClient.EXPECT().GetAllGroupServiceUsers(&workspaceUUID, &groupUUID).Return([]*types.ServiceUser{
		{
			UUID:      &serviceUser1UUID,
			Name:      serviceUser1Name,
			Workspace: workspaceUUIDString,
		},
		{
			UUID:      &serviceUser2UUID,
			Name:      serviceUser2Name,
			Workspace: workspaceUUIDString,
		},
		{
			UUID:      &serviceUser3UUID,
			Name:      serviceUser3Name,
			Workspace: workspaceUUIDString,
		},
	}, nil).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: provider.TestProviderConfig + fmt.Sprintf(`
				data "sotoon_iam_service_users" "test" {
					workspace_id = "%s"
					group_id = "%s"
				}`, workspaceUUIDString, groupUUIDString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sotoon_iam_service_users.test", "group_id", groupUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_service_users.test", "workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_service_users.test", "service_users.#", "3"),
					resource.TestCheckResourceAttr("data.sotoon_iam_service_users.test", "service_users.0.id", serviceUser1UUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_service_users.test", "service_users.0.workspace_id", workspaceUUIDString),
					resource.TestCheckResourceAttr("data.sotoon_iam_service_users.test", "service_users.0.name", serviceUser1Name),
				),
			},
		},
	})
}
