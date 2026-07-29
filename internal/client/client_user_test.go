package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	uuid "github.com/satori/go.uuid"
	sdk "github.com/sotoon/sotoon-sdk-go/sdk"
	iam "github.com/sotoon/sotoon-sdk-go/sdk/core/iam_v1"
)

func TestUnitRemoveUserFromWorkspace(t *testing.T) {
	const (
		workspaceID = "11111111-1111-1111-1111-111111111111"
		userID      = "22222222-2222-2222-2222-222222222222"
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		wantPath := "/iam/v1/api/v1/workspace/" + workspaceID + "/user/" + userID + "/"
		if r.URL.Path != wantPath {
			t.Errorf("path = %s, want %s", r.URL.Path, wantPath)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	handler, err := iam.NewHandler(server.URL, "test-token")
	if err != nil {
		t.Fatalf("create IAM handler: %v", err)
	}
	c := &Client{Workspace: workspaceID, sotoonSdk: &sdk.SDK{Iam_v1: handler}}
	if err := c.RemoveUserFromWorkspace(context.Background(), userID); err != nil {
		t.Fatalf("remove user from workspace: %v", err)
	}
}

func TestUnitRemoveUserFromWorkspaceRejectsUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()
	handler, err := iam.NewHandler(server.URL, "test-token")
	if err != nil {
		t.Fatalf("create IAM handler: %v", err)
	}
	c := &Client{Workspace: "11111111-1111-1111-1111-111111111111", sotoonSdk: &sdk.SDK{Iam_v1: handler}}
	if err := c.RemoveUserFromWorkspace(context.Background(), "22222222-2222-2222-2222-222222222222"); err == nil {
		t.Fatal("remove user succeeded for HTTP 403")
	}
}

func TestUnitBulkAddUsersToRoleSendsItemsForEveryUser(t *testing.T) {
	const (
		workspaceID = "11111111-1111-1111-1111-111111111111"
		roleID      = "22222222-2222-2222-2222-222222222222"
	)
	userIDs := []string{"33333333-3333-3333-3333-333333333333", "44444444-4444-4444-4444-444444444444"}
	wantedItems := map[string]string{"namespace": "example"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		wantPath := "/iam/v1/api/v1/workspace/" + workspaceID + "/role/" + roleID + "/bulk-add-users/"
		if r.URL.Path != wantPath {
			t.Errorf("path = %s, want %s", r.URL.Path, wantPath)
		}
		var request iam.IamBulkAddUsersToRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if len(request.Users) != len(userIDs) || request.Items == nil || len(*request.Items) != len(userIDs) {
			t.Fatalf("request = %#v, want one item map per user", request)
		}
		for _, items := range *request.Items {
			if items["namespace"] != wantedItems["namespace"] {
				t.Fatalf("items = %#v, want %#v", items, wantedItems)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("[]"))
	}))
	defer server.Close()
	handler, err := iam.NewHandler(server.URL, "test-token")
	if err != nil {
		t.Fatalf("create IAM handler: %v", err)
	}
	c := &Client{Workspace: workspaceID, sotoonSdk: &sdk.SDK{Iam_v1: handler}}
	roleUUID := uuid.FromStringOrNil(roleID)
	if err := c.BulkAddUsersToRole(context.Background(), roleUUID, userIDs, wantedItems); err != nil {
		t.Fatalf("bulk add users to role: %v", err)
	}
}
