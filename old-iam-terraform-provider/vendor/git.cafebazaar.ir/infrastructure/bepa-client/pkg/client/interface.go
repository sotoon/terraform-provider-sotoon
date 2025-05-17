package client

import (
	"net/url"
	"time"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	uuid "github.com/satori/go.uuid"
)

// Client represents bepa client interface
type Client interface {
	GetOrganizations() ([]*types.Organization, error)
	GetOrganization(*uuid.UUID) (*types.Organization, error)
	GetOrganizationWorkspaces(*uuid.UUID) ([]*types.Workspace, error)
	GetOrganizationWorkspace(*uuid.UUID, *uuid.UUID) (*types.Workspace, error)

	GetWorkspaces() ([]*types.Workspace, error)
	GetWorkspaceByName(name string) (*types.Workspace, error)
	GetWorkspaceByNameAndOrgName(name string, organizationName string) (*types.WorkspaceWithOrganization, error)
	GetWorkspace(uuid *uuid.UUID) (*types.Workspace, error)
	CreateWorkspace(name string) (*types.Workspace, error)
	DeleteWorkspace(uuid *uuid.UUID) error
	GetWorkspaceUsers(uuid *uuid.UUID) ([]*types.User, error)
	GetWorkspaceRoles(uuid *uuid.UUID) ([]*types.Role, error)
	GetWorkspaceRules(uuid *uuid.UUID) ([]*types.Rule, error)
	AddUserToWorkspace(userUUID, workspaceUUID *uuid.UUID) error
	RemoveUserFromWorkspace(userUUID, workspaceUUID *uuid.UUID) error
	SetConfigDefaultWorkspace(uuid *uuid.UUID) error

	CreateRole(roleName string, workspaceUUID *uuid.UUID) (*types.Role, error)
	UpdateRole(roleUUID *uuid.UUID, roleName string, workspaceUUID *uuid.UUID) (*types.Role, error)
	GetRole(roleUUID, workspaceUUID *uuid.UUID) (*types.RoleRes, error)
	GetRoleByName(roleName, workspaceName string) (*types.RoleRes, error)
	GetAllRoles() ([]*types.Role, error)
	GetRoleUsers(roleUUID, workspaceUUID *uuid.UUID) ([]*types.User, error)
	GetRoleRules(roleUUID, workspaceUUID *uuid.UUID) ([]*types.Rule, error)
	DeleteRole(roleUUID, workspaceUUID *uuid.UUID) error
	BindRoleToUser(workspaceUUID, roleUUID, userUUID *uuid.UUID, items map[string]string) error
	UnbindRoleFromUser(workspaceUUID, roleUUID, userUUID *uuid.UUID, items map[string]string) error
	GetBindedRoleToUserItems(workspaceUUID, roleUUID, userUUID *uuid.UUID) (map[string]string, error)
	GetBindedRoleToServiceUserItems(workspaceUUID, roleUUID, userUUID *uuid.UUID) (map[string]string, error)
	GetBindedRoleToGroupItems(workspaceUUID, roleUUID, userUUID *uuid.UUID) (map[string]string, error)

	GetRule(ruleUUID, workspaceUUID *uuid.UUID) (*types.Rule, error)
	GetRuleByName(ruleName, workspaceName string) (*types.Rule, error)
	CreateRule(ruleName string, workspaceUUID *uuid.UUID, ruleActions []string, object string, deny bool) (*types.Rule, error)
	DeleteRule(ruleUUID, workspaceUUID *uuid.UUID) error
	GetAllRules() ([]*types.Rule, error)
	GetAllUserRules(userUUID *uuid.UUID) ([]*types.Rule, error)
	BindRuleToRole(roleUUID, ruleUUID, workspaceUUID *uuid.UUID) error
	UnbindRuleFromRole(roleUUID, ruleUUID, workspaceUUID *uuid.UUID) error
	GetRuleRoles(ruleUUID, workspaceUUID *uuid.UUID) ([]*types.Role, error)
	UpdateRule(ruleUUID *uuid.UUID, ruleName string, workspaceUUID *uuid.UUID, ruleActions []string, object string, deny bool) (*types.Rule, error)

	CreateUser(userName, email, password string) (*types.User, error)
	GetUser(userUUID *uuid.UUID) (*types.User, error)
	GetMySelf() (*types.User, error)
	DeleteMySelf() error
	GetUserByEmail(email string, workspaceUUID *uuid.UUID) (*types.User, error)
	GetUserByName(userName string, workspaceUUID *uuid.UUID) (*types.User, error)
	GetUsers() ([]*types.User, error)
	DeleteUser(userUUID *uuid.UUID) error
	UpdateUser(userUUID *uuid.UUID, name, email, password string) error
	SetMyPassword(password string) error
	SetMyEmail(email string) error
	SetMyName(name string) error
	GetSecret(userUUID *uuid.UUID) (*types.UserSecret, error)
	RevokeSecret(userUUID *uuid.UUID) error
	SuspendUserInWorkspace(workspaceUUID *uuid.UUID, userUUID *uuid.UUID) error
	ActivateUserInWorkspace(workspaceUUID *uuid.UUID, userUUID *uuid.UUID) error
	InviteUser(workspaceUUID *uuid.UUID, email string) (*types.InvitationInfo, error)
	JoinByInvitationToken(name, password, invitationToken string) (*types.User, error)
	GetMyWorkspaces() ([]*types.WorkspaceWithOrganization, error)
	GetUserRoles(userUUID *uuid.UUID) ([]*types.RoleBinding, error)
	CreateUserTokenByCreds(email, password string) (*types.UserToken, error)
	SetConfigDefaultUserData(context, token, userUUID, email string) error
	SetCurrentContext(context string) error
	SuspendUser(userUUID *uuid.UUID) error
	ActivateUser(userUUID *uuid.UUID) error

	CreatePublicKeyForDefaultUser(title, keyType, key string) (*types.PublicKey, error)
	GetOneDefaultUserPublicKey(publicKeyUUID *uuid.UUID) (*types.PublicKey, error)
	GetAllDefaultUserPublicKeys() ([]*types.PublicKey, error)
	DeleteDefaultUserPublicKey(publicKeyUUID *uuid.UUID) error
	CreatePublicKeyFromFileForDefaultUser(title, fileAdd string) (*types.PublicKey, error)
	VerifyPublicKey(keyType string, key string, workspaceUUID string, username string, hostname string) (bool, error)

	GetAllUserKiseSecret() ([]*types.KiseSecret, error)
	DeleteUserKiseSecret(KiseSecretUUID *uuid.UUID) error
	CreateKiseSecretForDefaultUser() (*types.KiseSecret, error)

	Authorize(identity, userType, action, object string) error
	Identify(token string) (*types.UserRes, error)

	Do(method, path string, successCode int, req interface{}, resp interface{}) error
	SetAccessToken(token string)
	SetDefaultWorkspace(workspace string)
	SetUser(userUUID string)

	CreateMyUserTokenWithToken(secret string) (*types.UserToken, error)
	GetMyUserToken(UserTokenUUID *uuid.UUID) (*types.UserToken, error)
	GetAllMyUserTokens() (*[]types.UserToken, error)
	DeleteMyUserToken(UserTokenUUID *uuid.UUID) error

	GetAllServices() (*[]types.Service, error)
	GetService(name string) (*types.Service, error)

	DeleteServiceUserToken(serviceUserUUID, workspaceUUID, serviceUserTokenUUID *uuid.UUID) error
	GetAllServiceUserToken(serviceUserUUID, workspaceUUID *uuid.UUID) (*[]types.ServiceUserToken, error)
	CreateServiceUserToken(serviceUserUUID, workspaceUUID *uuid.UUID) (*types.ServiceUserToken, error)
	CreateServiceUser(serviceUserName string, workspace *uuid.UUID) (*types.ServiceUser, error)
	GetServiceUserByName(workspaceName string, serviceUserName string) (*types.ServiceUser, error)
	DeleteServiceUser(workspaceUUID, serviceUserUUID *uuid.UUID) error
	GetServiceUsers(workspaceUUID *uuid.UUID) ([]*types.ServiceUser, error)
	GetServiceUser(workspaceUUID, serviceUserUUID *uuid.UUID) (*types.ServiceUser, error)
	BindRoleToServiceUser(workspaceUUID, roleUUID, serviceUserUUID *uuid.UUID, items map[string]string) error
	UnbindRoleFromServiceUser(workspaceUUID, roleUUID, serviceUserUUID *uuid.UUID, items map[string]string) error
	GetRoleServiceUsers(roleUUID, workspaceUUID *uuid.UUID) ([]*types.ServiceUser, error)

	GetGroup(workspaceUUID, groupUUID *uuid.UUID) (*types.Group, error)
	GetAllGroups(workspaceUUID *uuid.UUID) ([]*types.Group, error)
	DeleteGroup(workspaceUUID, groupUUID *uuid.UUID) error
	GetGroupByName(workspaceName string, groupName string) (*types.Group, error)
	CreateGroup(groupName string, workspace *uuid.UUID) (*types.GroupRes, error)
	GetGroupUser(workspaceUUID, groupUUID, userUUID *uuid.UUID) (*types.User, error)
	GetAllGroupUsers(workspaceUUID, groupUUID *uuid.UUID) ([]*types.User, error)
	GetAllGroupServiceUsers(workspaceUUID, groupUUID *uuid.UUID) ([]*types.ServiceUser, error)
	UnbindUserFromGroup(workspaceUUID, groupUUID, userUUID *uuid.UUID) error
	BindGroup(groupName string, workspace, groupUUID, userUUID *uuid.UUID) error
	GetRoleGroups(roleUUID, workspaceUUID *uuid.UUID) ([]*types.Group, error)
	BindRoleToGroup(workspaceUUID, roleUUID, groupUUID *uuid.UUID, items map[string]string) error
	UnbindRoleFromGroup(workspaceUUID, roleUUID, groupUUID *uuid.UUID, items map[string]string) error
	BindServiceUserToGroup(worspaceUUID, groupUUID, serviceUserUUID *uuid.UUID) error
	UnbindServiceUserFromGroup(worspaceUUID, groupUUID, serviceUserUUID *uuid.UUID) error
	GetGroupServiceUser(worspaceUUID, groupUUID, serviceUserUUID *uuid.UUID) (*types.ServiceUser, error)

	GetServerURL() string

	GetAllDefaultBackupKeys() ([]*types.BackupKey, error)
	GetOneDefaultBackupKey(BackupKeyUUID *uuid.UUID) (*types.BackupKey, error)
	DeleteDefaultWorkspaceBackupKey(backupKeyUUID *uuid.UUID) error
	CreateBackupKeyForDefaultWorkspace(title, keyType, key string) (*types.BackupKey, error)
	CreateBackupKeyFromFileForDefaultUser(title, fileAdd string) (*types.BackupKey, error)

	GetBepaURL() (*url.URL, error)
}

type Cache interface {
	Get(string) (interface{}, bool)
	Set(string, interface{}, time.Duration)
}
