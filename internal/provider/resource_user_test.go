package provider

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iam "github.com/sotoon/sotoon-sdk-go/sdk/core/iam_v1"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

type fakeUserWorkspaceRemover struct {
	user       *iam.IamUser
	lookupErr  error
	removeErr  error
	removedIDs []string
}

func (f *fakeUserWorkspaceRemover) GetUserByEmail(context.Context, string) (*iam.IamUser, error) {
	return f.user, f.lookupErr
}

func (f *fakeUserWorkspaceRemover) RemoveUserFromWorkspace(_ context.Context, userID string) error {
	f.removedIDs = append(f.removedIDs, userID)
	return f.removeErr
}

func TestUnitResourceUserDeleteRemovesUserFromWorkspace(t *testing.T) {
	const userID = "22222222-2222-2222-2222-222222222222"
	fake := &fakeUserWorkspaceRemover{user: &iam.IamUser{Uuid: userID}}
	d := userResourceData(t)

	if diags := resourceUserDeleteWithClient(context.Background(), d, fake); diags.HasError() {
		t.Fatalf("delete returned diagnostics: %v", diags)
	}
	if len(fake.removedIDs) != 1 || fake.removedIDs[0] != userID {
		t.Fatalf("removed IDs = %#v, want %q", fake.removedIDs, userID)
	}
	if d.Id() != "" {
		t.Fatalf("resource ID was not cleared: %q", d.Id())
	}
}

func TestUnitResourceUserDeleteIsIdempotentWhenUserIsAbsent(t *testing.T) {
	fake := &fakeUserWorkspaceRemover{lookupErr: client.ErrNotFound}
	d := userResourceData(t)

	if diags := resourceUserDeleteWithClient(context.Background(), d, fake); diags.HasError() {
		t.Fatalf("delete returned diagnostics: %v", diags)
	}
	if len(fake.removedIDs) != 0 {
		t.Fatalf("remove was called for an absent user: %#v", fake.removedIDs)
	}
	if d.Id() != "" {
		t.Fatalf("resource ID was not cleared: %q", d.Id())
	}
}

func TestUnitResourceUserDeletePreservesStateOnRemovalFailure(t *testing.T) {
	const userID = "22222222-2222-2222-2222-222222222222"
	fake := &fakeUserWorkspaceRemover{
		user:      &iam.IamUser{Uuid: userID},
		removeErr: errors.New("permission denied"),
	}
	d := userResourceData(t)

	if diags := resourceUserDeleteWithClient(context.Background(), d, fake); !diags.HasError() {
		t.Fatal("delete succeeded despite removal failure")
	}
	if d.Id() == "" {
		t.Fatal("resource ID was cleared after removal failure")
	}
}

func userResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, resourceUser().Schema, map[string]interface{}{
		"email": "person@example.com",
	})
	d.SetId("person@example.com")
	return d
}
