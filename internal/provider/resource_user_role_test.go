package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	iam "github.com/sotoon/sotoon-sdk-go/sdk/core/iam_v1"
)

type userRoleAddCall struct {
	userIDs []string
	items   map[string]string
}

type fakeUserRoleBindingClient struct {
	users       []iam.IamUserWithRoleItems
	addCalls    []userRoleAddCall
	unbindCalls []string
}

func (f *fakeUserRoleBindingClient) GetRoleUsers(context.Context, *uuid.UUID) ([]iam.IamUserWithRoleItems, error) {
	return f.users, nil
}

func (f *fakeUserRoleBindingClient) BulkAddUsersToRole(_ context.Context, _ uuid.UUID, userIDs []string, items map[string]string) error {
	f.addCalls = append(f.addCalls, userRoleAddCall{userIDs: append([]string(nil), userIDs...), items: cloneRoleItems(items)})
	for _, userID := range userIDs {
		index := -1
		for i := range f.users {
			if f.users[i].Uuid == userID {
				index = i
				break
			}
		}
		if index == -1 {
			f.users = append(f.users, iam.IamUserWithRoleItems{Uuid: userID})
			index = len(f.users) - 1
		}
		f.users[index].Items = append(f.users[index].Items, cloneRoleItems(items))
	}
	return nil
}

func (f *fakeUserRoleBindingClient) UnbindRoleFromUser(_ context.Context, _ *uuid.UUID, userID *uuid.UUID) error {
	f.unbindCalls = append(f.unbindCalls, userID.String())
	remaining := f.users[:0]
	for _, user := range f.users {
		if user.Uuid != userID.String() {
			remaining = append(remaining, user)
		}
	}
	f.users = remaining
	return nil
}

func TestUnitUserRoleCreateAddsDistinctItemsForExistingRole(t *testing.T) {
	const userID = "11111111-1111-1111-1111-111111111111"
	first := map[string]string{"namespace": "first"}
	fake := &fakeUserRoleBindingClient{users: []iam.IamUserWithRoleItems{{Uuid: userID, Items: []map[string]string{first}}}}
	d := userRoleResourceData(t, userID, map[string]interface{}{"namespace": "second"})
	if diags := resourceUserRoleCreateWithClient(context.Background(), d, fake); diags.HasError() {
		t.Fatalf("create returned diagnostics: %v", diags)
	}
	if len(fake.addCalls) != 1 || fake.addCalls[0].items["namespace"] != "second" {
		t.Fatalf("add calls = %#v, want one call for the second item set", fake.addCalls)
	}
	if len(fake.users) != 1 || len(fake.users[0].Items) != 2 {
		t.Fatalf("remote user items = %#v, want both item sets", fake.users)
	}
}

func TestUnitUserRoleDeletePreservesOtherItemSets(t *testing.T) {
	const userID = "11111111-1111-1111-1111-111111111111"
	first := map[string]string{"namespace": "first"}
	second := map[string]string{"namespace": "second"}
	fake := &fakeUserRoleBindingClient{users: []iam.IamUserWithRoleItems{{Uuid: userID, Items: []map[string]string{first, second}}}}
	d := userRoleResourceData(t, userID, map[string]interface{}{"namespace": "second"})
	d.SetId("managed-binding")
	if diags := resourceUserRoleDeleteWithClient(context.Background(), d, fake); diags.HasError() {
		t.Fatalf("delete returned diagnostics: %v", diags)
	}
	if len(fake.unbindCalls) != 1 {
		t.Fatalf("unbind calls = %#v, want one", fake.unbindCalls)
	}
	if len(fake.users) != 1 || len(fake.users[0].Items) != 1 || !roleItemsEqual(fake.users[0].Items[0], first) {
		t.Fatalf("remote user items = %#v, want only first item set", fake.users)
	}
	if d.Id() != "" {
		t.Fatalf("resource ID was not cleared: %q", d.Id())
	}
}

func TestUnitUserRoleBindingsHaveDistinctIDsForDistinctItems(t *testing.T) {
	const userID = "11111111-1111-1111-1111-111111111111"
	first := userRoleResourceData(t, userID, map[string]interface{}{"namespace": "first"})
	second := userRoleResourceData(t, userID, map[string]interface{}{"namespace": "second"})
	fake := &fakeUserRoleBindingClient{users: []iam.IamUserWithRoleItems{{
		Uuid: userID,
		Items: []map[string]string{
			{"namespace": "first"},
			{"namespace": "second"},
		},
	}}}
	if diags := resourceUserRoleReadWithClient(context.Background(), first, fake); diags.HasError() {
		t.Fatalf("read first binding: %v", diags)
	}
	if diags := resourceUserRoleReadWithClient(context.Background(), second, fake); diags.HasError() {
		t.Fatalf("read second binding: %v", diags)
	}
	if first.Id() == second.Id() {
		t.Fatalf("distinct item sets produced the same resource ID: %q", first.Id())
	}
}

func userRoleResourceData(t *testing.T, userID string, items map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceUserRole().Schema, map[string]interface{}{
		"role_id":  "22222222-2222-2222-2222-222222222222",
		"user_ids": []interface{}{userID},
		"items":    items,
	})
}

func cloneRoleItems(items map[string]string) map[string]string {
	cloned := make(map[string]string, len(items))
	for key, value := range items {
		cloned[key] = value
	}
	return cloned
}
