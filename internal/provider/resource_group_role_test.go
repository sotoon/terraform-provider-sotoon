package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	iam "github.com/sotoon/sotoon-sdk-go/sdk/core/iam_v1"
)

type fakeGroupRoleBindingClient struct {
	roles       []iam.IamRole
	addCalls    [][]iam.IamRoleItem
	unbindCalls []string
}

func (f *fakeGroupRoleBindingClient) GetWorkspaceGroupRoleList(context.Context, *uuid.UUID, *uuid.UUID) ([]iam.IamRole, error) {
	return f.roles, nil
}

func (f *fakeGroupRoleBindingClient) BulkAddRolesToGroup(_ context.Context, _ *uuid.UUID, roles []iam.IamRoleItem) error {
	f.addCalls = append(f.addCalls, roles)
	for _, added := range roles {
		index := -1
		for i := range f.roles {
			if f.roles[i].Uuid == added.RoleUuid {
				index = i
				break
			}
		}
		if index == -1 {
			f.roles = append(f.roles, iam.IamRole{Uuid: added.RoleUuid})
			index = len(f.roles) - 1
		}
		if added.ItemsList != nil {
			for _, items := range *added.ItemsList {
				f.roles[index].Items = append(f.roles[index].Items, cloneRoleItems(items))
			}
		}
	}
	return nil
}

func (f *fakeGroupRoleBindingClient) UnbindRoleFromGroup(_ context.Context, roleID, _ *uuid.UUID) error {
	f.unbindCalls = append(f.unbindCalls, roleID.String())
	remaining := f.roles[:0]
	for _, role := range f.roles {
		if role.Uuid != roleID.String() {
			remaining = append(remaining, role)
		}
	}
	f.roles = remaining
	return nil
}

func TestUnitGroupRoleCreateAddsDistinctItemsForExistingRole(t *testing.T) {
	const roleID = "22222222-2222-2222-2222-222222222222"
	first := map[string]string{"namespace": "first"}
	fake := &fakeGroupRoleBindingClient{roles: []iam.IamRole{{Uuid: roleID, Items: []map[string]string{first}}}}
	d := groupRoleResourceData(t, roleID, map[string]interface{}{"namespace": "second"})
	workspaceID := uuid.FromStringOrNil("33333333-3333-3333-3333-333333333333")
	if diags := resourceGroupRoleCreateWithClient(context.Background(), d, fake, &workspaceID); diags.HasError() {
		t.Fatalf("create returned diagnostics: %v", diags)
	}
	if len(fake.addCalls) != 1 || len(fake.roles) != 1 || len(fake.roles[0].Items) != 2 {
		t.Fatalf("add calls = %#v, remote roles = %#v; want both item sets", fake.addCalls, fake.roles)
	}
}

func TestUnitGroupRoleDeletePreservesOtherItemSets(t *testing.T) {
	const roleID = "22222222-2222-2222-2222-222222222222"
	first := map[string]string{"namespace": "first"}
	second := map[string]string{"namespace": "second"}
	fake := &fakeGroupRoleBindingClient{roles: []iam.IamRole{{Uuid: roleID, Items: []map[string]string{first, second}}}}
	d := groupRoleResourceData(t, roleID, map[string]interface{}{"namespace": "second"})
	d.SetId("managed-binding")
	workspaceID := uuid.FromStringOrNil("33333333-3333-3333-3333-333333333333")
	if diags := resourceGroupRoleDeleteWithClient(context.Background(), d, fake, &workspaceID); diags.HasError() {
		t.Fatalf("delete returned diagnostics: %v", diags)
	}
	if len(fake.unbindCalls) != 1 {
		t.Fatalf("unbind calls = %#v, want one", fake.unbindCalls)
	}
	if len(fake.roles) != 1 || len(fake.roles[0].Items) != 1 || !roleItemsEqual(fake.roles[0].Items[0], first) {
		t.Fatalf("remote role items = %#v, want only first item set", fake.roles)
	}
	if d.Id() != "" {
		t.Fatalf("resource ID was not cleared: %q", d.Id())
	}
}

func TestUnitGroupRoleBindingsHaveDistinctIDsForDistinctItems(t *testing.T) {
	const roleID = "22222222-2222-2222-2222-222222222222"
	first := groupRoleResourceData(t, roleID, map[string]interface{}{"namespace": "first"})
	second := groupRoleResourceData(t, roleID, map[string]interface{}{"namespace": "second"})
	fake := &fakeGroupRoleBindingClient{roles: []iam.IamRole{{
		Uuid: roleID,
		Items: []map[string]string{
			{"namespace": "first"},
			{"namespace": "second"},
		},
	}}}
	workspaceID := uuid.FromStringOrNil("33333333-3333-3333-3333-333333333333")
	if diags := resourceGroupRoleReadWithClient(context.Background(), first, fake, &workspaceID); diags.HasError() {
		t.Fatalf("read first binding: %v", diags)
	}
	if diags := resourceGroupRoleReadWithClient(context.Background(), second, fake, &workspaceID); diags.HasError() {
		t.Fatalf("read second binding: %v", diags)
	}
	if first.Id() == second.Id() {
		t.Fatalf("distinct item sets produced the same resource ID: %q", first.Id())
	}
}

func groupRoleResourceData(t *testing.T, roleID string, items map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceGroupRole().Schema, map[string]interface{}{
		"group_id": "11111111-1111-1111-1111-111111111111",
		"role_ids": []interface{}{roleID},
		"items":    items,
	})
}
