package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	uuid "github.com/satori/go.uuid"
	iamclient "github.com/sotoon/iam-client/pkg/client"
	"github.com/sotoon/iam-client/pkg/types"
)

var ErrNotFound = errors.New("resource not found")

// API response structs

type ObjectMetadata struct {
	UID       string `json:"uid,omitempty"`
	Name      string `json:"name"`
	Workspace string `json:"workspace"`
}

// Client is a unified wrapper for both Compute and IAM APIs.
type Client struct {
	ComputeBaseURL string
	APIToken       string
	Workspace      string
	WorkspaceUUID  *uuid.UUID
	HTTPClient     *http.Client
	IAMClient      iamclient.Client
}

// NewClient creates a new unified API client for both Compute and IAM.
func NewClient(host, token, workspace, userID string) (*Client, error) {
	if host == "" || token == "" || workspace == "" || userID == "" {
		return nil, fmt.Errorf("host, token, workspace, and userID must not be empty")
	}

	iam, err := iamclient.NewClient(token, "https://bepa.sotoon.ir", workspace, userID, 2)
	if err != nil {
		return nil, fmt.Errorf("failed to create sotoon iam client: %w", err)
	}
	workspaceUUID, err := uuid.FromString(workspace)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace_uuid format: %w", err)
	}

	return &Client{
		ComputeBaseURL: fmt.Sprintf("%s/compute/v2/thr1/workspaces/%s", host, workspace),
		APIToken:       token,
		Workspace:      workspace,
		WorkspaceUUID:  &workspaceUUID,
		HTTPClient:     &http.Client{Timeout: 30 * time.Second},
		IAMClient:      iam,
	}, nil
}

// --- IAM User Functions ---

func (c *Client) InviteUser(ctx context.Context, email string) (*types.InvitationInfo, error) {

	invitationInfo, err := c.IAMClient.InviteUser(c.WorkspaceUUID, email)

	if err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal array into Go value of type") {
			tflog.Debug(ctx, "Successfully invited user (ignoring known unmarshal error)")
			return nil, nil
		}
		// Log the actual error before returning
		tflog.Error(ctx, "Failed to invite user", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return invitationInfo, nil
}

func (c *Client) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	user, err := c.IAMClient.GetUserByEmail(email, c.WorkspaceUUID)
	if err != nil {
		if err.Error() == "User not found" {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (c *Client) GetUserByID(ctx context.Context, userID string) (*types.User, error) {
	userUUID, err := uuid.FromString(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}
	user, err := c.IAMClient.GetUser(&userUUID)
	if err != nil {
		if err.Error() == "User not found" {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	userUUID, err := uuid.FromString(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}
	return c.IAMClient.DeleteUser(&userUUID)
}

func (c *Client) GetWorkspaceUsers(ctx context.Context, workspaceID *uuid.UUID) ([]*types.User, error) {
	return c.IAMClient.GetWorkspaceUsers(workspaceID)
}

func (c *Client) GetWorkspaceGroups(ctx context.Context, workspaceID *uuid.UUID) ([]*types.Group, error) {
	return c.IAMClient.GetAllGroups(workspaceID)
}

func (c *Client) GetWorkspaceGroupUsersList(ctx context.Context, workspaceID, groupID *uuid.UUID) ([]*types.User, error) {
	return c.IAMClient.GetAllGroupUserList(workspaceID, groupID)
}

func (c *Client) GetWorkspaceGroupRoleList(ctx context.Context, workspaceID, groupID *uuid.UUID) ([]*types.Role, error) {
	return c.IAMClient.GetWorkspaceGroupRoleList(*workspaceID, *groupID)
}

func (c *Client) GetAllGroupServiceUserList(ctx context.Context, workspaceID, groupID *uuid.UUID) ([]*types.ServiceUser, error) {
	return c.IAMClient.GetAllGroupServiceUserList(workspaceID, groupID)
}

func (c *Client) GetWorkspaceGroupDetail(ctx context.Context, workspaceID, groupID uuid.UUID) (*types.Group, error) {
	return c.IAMClient.GetWorkspaceGroupDetail(workspaceID, groupID)
}

func (c *Client) CreateGroup(ctx context.Context, name, description string) (*types.Group, error) {
	group, err := c.IAMClient.CreateGroup(name, description, c.WorkspaceUUID)
	if err != nil {
		return nil, err
	}

	return &types.Group{
		UUID:        group.UUID,
		Name:        group.Name,
		Description: group.Description,
	}, nil
}

func (c *Client) DeleteGroup(ctx context.Context, groupID string) error {
	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		return fmt.Errorf("invalid group ID format: %w", err)
	}
	return c.IAMClient.DeleteGroup(c.WorkspaceUUID, &groupUUID)
}

// --- IAM User-Token Functions ---

func (c *Client) CreateMyUserToken(name string, expiresAt *time.Time) (*types.UserToken, error) {
	return c.IAMClient.CreateMyUserToken(name, expiresAt)
}

func (c *Client) GetMyUserToken(tokenUUID *uuid.UUID) (*types.UserToken, error) {
	return c.IAMClient.GetMyUserToken(tokenUUID)
}

func (c *Client) GetAllMyUserTokenList() (*[]types.UserToken, error) {
	return c.IAMClient.GetAllMyUserTokenList()
}

func (c *Client) GetUser(userUUID *uuid.UUID) (*types.User, error) {
	return c.IAMClient.GetUser(userUUID)
}

func (c *Client) DeleteMyUserToken(tokenUUID *uuid.UUID) error {
	return c.IAMClient.DeleteMyUserToken(tokenUUID)
}

// --- IAM Public-Key Functions ---

func (c *Client) CreateMyUserPublicKey(title, keyType, key string) (*types.PublicKey, error) {
	return c.IAMClient.CreateMyUserPublicKey(title, keyType, key)
}

func (c *Client) GetOneDefaultUserPublicKey(keyUUID *uuid.UUID) (*types.PublicKey, error) {
	return c.IAMClient.GetOneDefaultUserPublicKey(keyUUID)
}

func (c *Client) GetAllMyUserPublicKeyList() ([]*types.PublicKey, error) {
	return c.IAMClient.GetAllMyUserPublicKeyList()
}

func (c *Client) DeleteDefaultUserPublicKey(keyUUID *uuid.UUID) error {
	return c.IAMClient.DeleteMyUserPublicKey(keyUUID)
}

// --- Group Functions ---

func (c *Client) UpdateGroup(ctx context.Context, groupID string, name string, description string) error {
	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		tflog.Error(ctx, "Failed to parse groupID string to UUID", map[string]interface{}{
			"groupID": groupID,
			"error":   err.Error(),
		})
		return err
	}
	err = c.IAMClient.UpdateGroup(*c.WorkspaceUUID, groupUUID, &name, &description, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) RemoveUserFromGroup(ctx context.Context, groupID string, userID string) error {
	tflog.Debug(ctx, "Attempting to remove user from group", map[string]interface{}{"userID": userID, "groupID": groupID})

	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		return fmt.Errorf("invalid group ID format: %w", err)
	}
	userUUID, err := uuid.FromString(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	// Call the UnbindUserFromGroup function with pointers to the UUIDs.
	err = c.IAMClient.UnbindUserFromGroup(c.WorkspaceUUID, &groupUUID, &userUUID)
	if err != nil {
		tflog.Error(ctx, "Failed to remove user from group via client", map[string]interface{}{
			"userID":  userID,
			"groupID": groupID,
			"error":   err.Error(),
		})
		return err
	}

	tflog.Info(ctx, "Successfully removed user from group via client", map[string]interface{}{"userID": userID, "groupID": groupID})
	return nil
}

func (c *Client) BindRoleToGroup(workspaceUUID *uuid.UUID, roleUUID *uuid.UUID, groupUUID *uuid.UUID, items map[string]string) error {
	return c.IAMClient.BindRoleToGroup(workspaceUUID, roleUUID, groupUUID, items)
}

func (c *Client) BulkAddRolesToGroup(workspaceUUID *uuid.UUID, groupUUID *uuid.UUID, rolesWithItems []types.RoleWithItems) error {
	return c.IAMClient.BulkAddRolesToGroup(*workspaceUUID, *groupUUID, rolesWithItems)
}

func (c *Client) UnbindRoleFromGroup(workspaceUUID, roleUUID, groupUUID *uuid.UUID, items map[string]string) error {
	return c.IAMClient.UnbindRoleFromGroup(workspaceUUID, roleUUID, groupUUID, items)
}

func (c *Client) BulkAddServiceUsersToGroup(workspaceUUID, groupUUID uuid.UUID, serviceUserUUIDs []uuid.UUID) ([]*types.GroupServiceUser, error) {
	return c.IAMClient.BulkAddServiceUsersToGroup(workspaceUUID, groupUUID, serviceUserUUIDs)
}

func (c *Client) UnbindServiceUserFromGroup(workspaceUUID, groupUUID, serviceUserUUID *uuid.UUID) error {
	return c.IAMClient.UnbindServiceUserFromGroup(workspaceUUID, groupUUID, serviceUserUUID)
}

func (c *Client) GetServiceUsers(workspaceUUID *uuid.UUID) ([]*types.ServiceUser, error) {
	return c.IAMClient.GetServiceUsers(workspaceUUID)
}

func (c *Client) GetServiceUser(workspaceUUID, serviceUserUUID *uuid.UUID) (*types.ServiceUser, error) {
	return c.IAMClient.GetServiceUser(workspaceUUID, serviceUserUUID)
}

func (c *Client) GetWorkspaceServiceUserList(workspaceUUID uuid.UUID) ([]*types.ServiceUserWithCompactRole, error) {
	return c.IAMClient.GetWorkspaceServiceUserList(workspaceUUID)
}

func (c *Client) GetWorkspaceServiceUserDetail(workspaceUUID, serviceUserUUID uuid.UUID) (*types.ServiceUserWithCompactRole, error) {
	return c.IAMClient.GetWorkspaceServiceUserDetail(workspaceUUID, serviceUserUUID)
}

func (c *Client) CreateServiceUser(serviceUserName, description string, workspace *uuid.UUID) (*types.ServiceUser, error) {
	return c.IAMClient.CreateServiceUser(serviceUserName, description, workspace)
}

func (c *Client) DeleteServiceUser(workspaceUUID, serviceUserUUID *uuid.UUID) error {
	return c.IAMClient.DeleteServiceUser(workspaceUUID, serviceUserUUID)
}

func (c *Client) UpdateServiceUser(workspaceUUID, serviceUserUUID uuid.UUID, name, description string) (*types.ServiceUser, error) {
	return c.IAMClient.UpdateServiceUser(workspaceUUID, serviceUserUUID, name, description)
}

func (c *Client) GetWorkspaceServiceUserTokenList(serviceUserUUID, workspaceUUID *uuid.UUID) (*[]types.ServiceUserToken, error) {
	return c.IAMClient.GetWorkspaceServiceUserTokenList(serviceUserUUID, workspaceUUID)
}

func (c *Client) CreateServiceUserToken(serviceUserUUID, workspaceUUID *uuid.UUID, name string, expiresAt *time.Time) (*types.ServiceUserToken, error) {
	return c.IAMClient.CreateServiceUserToken(serviceUserUUID, workspaceUUID, name, expiresAt)
}

func (c *Client) DeleteServiceUserToken(serviceUserUUID, workspaceUUID, serviceUserTokenUUID *uuid.UUID) error {
	return c.IAMClient.DeleteServiceUserToken(serviceUserUUID, workspaceUUID, serviceUserTokenUUID)
}

func (c *Client) GetWorkspaceServiceUserPublicKeyList(workspaceUUID, serviceUserUUID uuid.UUID) ([]*types.ServiceUserPublicKey, error) {
	return c.IAMClient.GetWorkspaceServiceUserPublicKeyList(workspaceUUID, serviceUserUUID)
}

func (c *Client) CreateServiceUserPublicKey(workspaceUUID, serviceUserUUID uuid.UUID, name, publicKey string) (*types.ServiceUserPublicKey, error) {
	return c.IAMClient.CreateServiceUserPublicKey(workspaceUUID, serviceUserUUID, name, publicKey)
}

func (c *Client) DeleteServiceUserPublicKey(workspaceUUID, serviceUserUUID, publicKeyUUID uuid.UUID) error {
	return c.IAMClient.DeleteServiceUserPublicKey(workspaceUUID, serviceUserUUID, publicKeyUUID)
}

func (c *Client) UnbindRoleFromServiceUser(workspaceUUID, roleUUID, serviceUserUUID *uuid.UUID, items map[string]string) error {
	return c.IAMClient.UnbindRoleFromServiceUser(workspaceUUID, roleUUID, serviceUserUUID, items)
}

func (c *Client) GetRoleServiceUsers(roleUUID, workspaceUUID *uuid.UUID) ([]*types.ServiceUser, error) {
	return c.IAMClient.GetRoleServiceUsers(roleUUID, workspaceUUID)
}

func (c *Client) BulkAddServiceUsersToRole(workspaceUUID, roleUUID uuid.UUID, serviceUserUUIDs []uuid.UUID) error {
	return c.IAMClient.BulkAddServiceUsersToRole(workspaceUUID, roleUUID, serviceUserUUIDs)
}

// --- IAM Role Functions ---

func (c *Client) GetWorkspaceRoles(ctx context.Context) ([]*types.Role, error) {
	return c.IAMClient.GetWorkspaceRoles(c.WorkspaceUUID)
}

func (c *Client) CreateRole(ctx context.Context, name, description string) (*types.RoleWithCompactWorkspace, error) {
	return c.IAMClient.CreateRole(name, description, c.WorkspaceUUID)
}

func (c *Client) GetRole(ctx context.Context, roleUUID *uuid.UUID) (*types.RoleRes, error) {
	return c.IAMClient.GetRole(roleUUID, c.WorkspaceUUID)
}

func (c *Client) GetRoleByName(ctx context.Context, roleName string) (*types.RoleRes, error) {
	return c.IAMClient.GetRoleByName(roleName, c.Workspace)
}

func (c *Client) DeleteRole(ctx context.Context, roleID string) error {
	id, err := uuid.FromString(roleID)
	if err != nil {
		return fmt.Errorf("invalid role ID %q: %w", roleID, err)
	}
	return c.IAMClient.DeleteRole(&id, c.WorkspaceUUID)
}

func (c *Client) UpdateRole(ctx context.Context, roleUUID *uuid.UUID, name string) (*types.Role, error) {
	return c.IAMClient.UpdateRole(roleUUID, name, c.WorkspaceUUID)
}

func (c *Client) BulkAddRulesToRole(ctx context.Context, roleUUID uuid.UUID, ruleUUIDs []uuid.UUID) error {
	return c.IAMClient.BulkAddRulesToRole(*c.WorkspaceUUID, roleUUID, ruleUUIDs)
}

func (c *Client) GetRoleRules(ctx context.Context, roleUUID *uuid.UUID) ([]*types.Rule, error) {
	return c.IAMClient.GetRoleRules(roleUUID, c.WorkspaceUUID)
}

func (c *Client) UnbindRuleFromRole(ctx context.Context, roleUUID *uuid.UUID, ruleUUID *uuid.UUID) error {
	return c.IAMClient.UnbindRuleFromRole(roleUUID, ruleUUID, c.WorkspaceUUID)
}

// --- IAM Rule Functions ---

func (c *Client) GetWorkspaceRules(ctx context.Context) ([]*types.Rule, error) {
	return c.IAMClient.GetWorkspaceRules(c.WorkspaceUUID)
}

func (c *Client) GetRule(ctx context.Context, ruleUUID *uuid.UUID) (*types.Rule, error) {
	return c.IAMClient.GetRule(ruleUUID, c.WorkspaceUUID)
}
