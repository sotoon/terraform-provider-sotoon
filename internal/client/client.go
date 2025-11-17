package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	uuid "github.com/satori/go.uuid"
	sdk "github.com/sotoon/sotoon-sdk-go/sdk"
	iam "github.com/sotoon/sotoon-sdk-go/sdk/core/iam_v1"
	"github.com/sotoon/sotoon-sdk-go/sdk/interceptors"
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
	UserID         string
	WorkspaceUUID  *uuid.UUID
	HTTPClient     *http.Client
	sotoonSdk      *sdk.SDK
}

type logger struct {
}

func (a *logger) BeforeRequest(data interceptors.InterceptorData) (interceptors.InterceptorData, error) {
	var body []byte
	if data.Request != nil && data.Request.Body != nil {
		body, _ = io.ReadAll(data.Request.Body)
		data.Request.Body = io.NopCloser(bytes.NewReader(body))
	}

	tflog.Info(data.Ctx, "BeforeRequest", map[string]interface{}{
		"id":     data.ID,
		"method": data.Request.Method,
		"url":    data.Request.URL,
		"body":   string(body),
	})
	return data, nil
}

func (a *logger) AfterResponse(data interceptors.InterceptorData) (interceptors.InterceptorData, error) {
	body, _ := io.ReadAll(data.Response.Body)
	tflog.Info(data.Ctx, "AfterResponse", map[string]interface{}{
		"id":       data.ID,
		"method":   data.Request.Method,
		"url":      data.Request.URL,
		"response": string(body),
	})
	data.Response.Body = io.NopCloser(bytes.NewReader(body))
	return data, nil
}

// NewClient creates a new unified API client for both Compute and IAM.
func NewClient(host, token, workspace, userID string, shouldLog bool) (*Client, error) {
	if host == "" || token == "" || workspace == "" || userID == "" {
		return nil, fmt.Errorf("host, token, workspace, and userID must not be empty")
	}

	interceptorsArray := make([]interceptors.Interceptor, 0, 4)

	if shouldLog {
		interceptorsArray = append(interceptorsArray, &logger{})
	}

	interceptorsArray = append(interceptorsArray,
		interceptors.NewTreatAsErrorInterceptor(
			interceptors.NewTreatAsErrorInterceptor_ErrorDetectorAll(),
		),
		interceptors.NewCircuitBreakerInterceptor(interceptors.CircuteBreakerForJust429, false),
		interceptors.NewRetryInterceptor(
			interceptors.NewDefaultInterceptorTransport(token),
			interceptors.NewRetryInterceptor_ExponentialBackoff(time.Second, time.Second*10),
			interceptors.NewRetryInterceptor_RetryDeciderAll(15),
		),
	)

	sotoonSdk, err := sdk.NewSDK(token, sdk.WithInterceptor(interceptorsArray...))

	if err != nil {
		return nil, fmt.Errorf("failed to create sotoon sdk: %w", err)
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
		sotoonSdk:      sotoonSdk,
		UserID:         userID,
	}, nil
}

// --- IAM User Functions ---

func (c *Client) InviteUser(ctx context.Context, email string) (*iam.IamUserInvitation, error) {

	res, err := c.sotoonSdk.Iam_v1.InviteUsersToWorkspaceWithResponse(ctx, c.Workspace, iam.IamInviteRequest{Emails: []string{email}})

	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 200 {
		return res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetUserByEmail(ctx context.Context, email string) (*iam.IamUser, error) {
	res, err := c.sotoonSdk.Iam_v1.ListWorkspaceUsersWithResponse(ctx, c.Workspace, &iam.ListWorkspaceUsersParams{Email: &email})
	if err != nil {
		return nil, err
	}
	if res.JSON200 != nil && len(*res.JSON200) > 0 {
		return &(*res.JSON200)[0], nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetWorkspaceUsers(ctx context.Context, workspaceID *uuid.UUID) ([]iam.IamUser, error) {
	res, err := c.sotoonSdk.Iam_v1.ListWorkspaceUsersWithResponse(ctx, workspaceID.String(), nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return *res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetWorkspaceUserByUUID(ctx context.Context, workspaceID *uuid.UUID, userID string) (*iam.IamUserWorkspaceDetailedUser, error) {
	res, err := c.sotoonSdk.Iam_v1.GetDetailedWorkspaceUserWithResponse(ctx, workspaceID.String(), userID)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetWorkspaceUserByEmail(ctx context.Context, workspaceID *uuid.UUID, email string) (*iam.IamUser, error) {
	res, err := c.sotoonSdk.Iam_v1.ListWorkspaceUsersWithResponse(ctx, workspaceID.String(), &iam.ListWorkspaceUsersParams{Email: &email})
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		if len(*res.JSON200) == 0 {
			return nil, ErrNotFound
		}
		return &(*res.JSON200)[0], nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetWorkspaceGroups(ctx context.Context, workspaceID *uuid.UUID) ([]iam.IamGroup, error) {
	res, err := c.sotoonSdk.Iam_v1.ListGroupsWithResponse(ctx, workspaceID.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return *res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetWorkspaceGroupUsersList(ctx context.Context, workspaceID, groupID *uuid.UUID) ([]iam.IamUser, error) {
	res, err := c.sotoonSdk.Iam_v1.ListGroupUsersWithResponse(ctx, workspaceID.String(), groupID.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return *res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetWorkspaceGroupRoleList(ctx context.Context, workspaceID, groupID *uuid.UUID) ([]iam.IamRole, error) {
	res, err := c.sotoonSdk.Iam_v1.ListGroupRolesWithResponse(ctx, workspaceID.String(), groupID.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return *res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetAllGroupServiceUserList(ctx context.Context, workspaceID, groupID *uuid.UUID) ([]iam.IamServiceUser, error) {
	res, err := c.sotoonSdk.Iam_v1.ListGroupServiceUsersWithResponse(ctx, workspaceID.String(), groupID.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return *res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetWorkspaceGroupDetail(ctx context.Context, workspaceID, groupID uuid.UUID) (*iam.IamGroupDetail, error) {
	res, err := c.sotoonSdk.Iam_v1.GetDetailedGroupWithResponse(ctx, workspaceID.String(), groupID.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) CreateGroup(ctx context.Context, name, description string) (*iam.IamGroup, error) {
	res, err := c.sotoonSdk.Iam_v1.CreateGroupWithResponse(ctx, c.Workspace,
		iam.IamRequestCreateGroup{
			Description: &description,
			Name:        name,
		},
	)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 201 {
		return res.JSON201, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) DeleteGroup(ctx context.Context, groupID string) error {
	_, err := c.sotoonSdk.Iam_v1.DeleteGroupWithResponse(ctx, c.Workspace, groupID)
	return err
}

// --- IAM User-Token Functions ---

func (c *Client) CreateMyUserToken(ctx context.Context, name string, expiresAt *time.Time) (*iam.IamUserToken, error) {
	res, err := c.sotoonSdk.Iam_v1.CreateUserTokenWithResponse(
		ctx, c.UserID,
		iam.IamReuqestUserTokenCreate{
			Name:      name,
			ExpiresAt: *expiresAt,
			IsHashed:  true,
		},
	)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 201 {
		return res.JSON201, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetMyUserToken(ctx context.Context, tokenUUID *uuid.UUID) (*iam.IamUserToken, error) {
	res, err := c.sotoonSdk.Iam_v1.ListUserTokensWithResponse(ctx, c.UserID)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 200 {
		for _, token := range *res.JSON200 {
			if token.Uuid == tokenUUID.String() {
				return &token, nil
			}
		}
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetAllMyUserTokenList(ctx context.Context) ([]iam.IamUserToken, error) {
	res, err := c.sotoonSdk.Iam_v1.ListUserTokensWithResponse(ctx, c.UserID)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return *res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetUserDetailed(ctx context.Context, userUUID *uuid.UUID) (*iam.IamUserWorkspaceDetailedUser, error) {
	res, err := c.sotoonSdk.Iam_v1.GetDetailedWorkspaceUserWithResponse(ctx, c.Workspace, userUUID.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetUser(ctx context.Context, userUUID *uuid.UUID) (*iam.IamUser, error) {
	res, err := c.sotoonSdk.Iam_v1.GetUserWithResponse(ctx, userUUID.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) DeleteMyUserToken(ctx context.Context, tokenUUID *uuid.UUID) error {
	_, err := c.sotoonSdk.Iam_v1.DeleteUserTokenWithResponse(ctx, c.UserID, tokenUUID.String())
	return err
}

// --- IAM Public-Key Functions ---

func (c *Client) CreateMyUserPublicKey(ctx context.Context, title, key string) (*iam.IamUserPublicKey, error) {
	res, err := c.sotoonSdk.Iam_v1.CreateUserPublicKeyWithResponse(ctx, c.UserID,
		iam.IamRequestCreateUserPublicKey{
			Title: title,
			Key:   key,
		})
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 201 {
		return res.JSON201, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})

	return nil, ErrNotFound
}

func (c *Client) GetUserPublicKey(ctx context.Context, keyUUID *uuid.UUID) (*iam.IamUserPublicKey, error) {
	res, err := c.sotoonSdk.Iam_v1.ListUserPublicKeysWithResponse(ctx, c.UserID)
	if err != nil {
		return nil, err
	}
	for _, key := range *res.JSON200 {
		if key.Uuid == keyUUID.String() {
			return &key, nil
		}
	}
	return nil, ErrNotFound

}

func (c *Client) GetAllMyUserPublicKeyList(ctx context.Context) ([]iam.IamUserPublicKey, error) {
	res, err := c.sotoonSdk.Iam_v1.ListUserPublicKeysWithResponse(ctx, c.UserID)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return *res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) DeleteUserPublicKey(ctx context.Context, keyUUID *uuid.UUID) error {
	_, err := c.sotoonSdk.Iam_v1.DeleteUserPublicKeyWithResponse(ctx, c.UserID, keyUUID.String())
	return err
}

// --- Group Functions ---

func (c *Client) UpdateGroup(ctx context.Context, groupID string, name string, description string) error {
	_, err := c.sotoonSdk.Iam_v1.UpdateGroupWithResponse(ctx, c.Workspace, groupID,
		iam.IamRequestCreateGroup{
			Name:        name,
			Description: &description,
		})
	return err
}

func (c *Client) GetAllGroupUserList(ctx context.Context, groupUUID *uuid.UUID) ([]iam.IamUser, error) {
	res, err := c.sotoonSdk.Iam_v1.ListGroupUsersWithResponse(ctx, c.Workspace, groupUUID.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return *res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) BulkAddUsersToGroup(ctx context.Context, groupUUID uuid.UUID, uuids []string) ([]iam.IamServiceUserGroup, error) {
	res, err := c.sotoonSdk.Iam_v1.BulkAddUsersToGroupWithResponse(ctx, c.Workspace, groupUUID.String(), iam.IamBulkAddUsersRequest{Users: uuids})
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 201 {
		return *res.JSON201, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) RemoveUserFromGroup(ctx context.Context, groupID string, userID string) error {
	// Call the UnbindUserFromGroup function with pointers to the UUIDs.
	_, err := c.sotoonSdk.Iam_v1.RemoveUserFromGroupWithResponse(ctx, c.Workspace, groupID, userID)
	return err
}

func (c *Client) BulkAddRolesToGroup(ctx context.Context, groupUUID *uuid.UUID, rolesWithItems []iam.IamRoleItem) error {
	_, err := c.sotoonSdk.Iam_v1.BulkAddRolesToGroupWithResponse(ctx, c.Workspace,
		groupUUID.String(), iam.IamBulkAddRolesRequest{
			Roles: rolesWithItems,
		})
	return err

}

func (c *Client) UnbindRoleFromGroup(ctx context.Context, roleUUID, groupUUID *uuid.UUID) error {
	_, err := c.sotoonSdk.Iam_v1.RemoveRoleFromGroupWithResponse(ctx, c.Workspace, roleUUID.String(), groupUUID.String())
	return err

}

func (c *Client) BulkAddServiceUsersToGroup(ctx context.Context, groupUUID uuid.UUID, serviceUserUUIDs []string) ([]iam.IamServiceUserGroup, error) {
	res, err := c.sotoonSdk.Iam_v1.BulkAddServiceUsersToGroupWithResponse(ctx, c.Workspace, groupUUID.String(), iam.IamBulkAddServiceUsersRequest{ServiceUsers: serviceUserUUIDs})
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 201 {
		return *res.JSON201, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) UnbindServiceUserFromGroup(ctx context.Context, groupUUID, serviceUserUUID *uuid.UUID) error {
	_, err := c.sotoonSdk.Iam_v1.RemoveServiceUserFromGroupWithResponse(ctx, c.Workspace, groupUUID.String(), serviceUserUUID.String())
	return err
}

func (c *Client) GetServiceUsers(ctx context.Context) ([]iam.IamServiceUser, error) {
	res, err := c.sotoonSdk.Iam_v1.ListServiceUsersWithResponse(ctx, c.Workspace)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return *res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetServiceUser(ctx context.Context, serviceUserUUID *uuid.UUID) (*iam.IamServiceUserDetailed, error) {
	res, err := c.sotoonSdk.Iam_v1.GetDetailedServiceUserWithResponse(ctx, c.Workspace, serviceUserUUID.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetWorkspaceServiceUserDetail(ctx context.Context, workspaceUUID, serviceUserUUID uuid.UUID) (*iam.IamServiceUserDetailed, error) {
	res, err := c.sotoonSdk.Iam_v1.GetDetailedServiceUserWithResponse(ctx, workspaceUUID.String(), serviceUserUUID.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 200 {
		return res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) CreateServiceUser(ctx context.Context, serviceUserName, description string) (*iam.IamServiceUser, error) {
	res, err := c.sotoonSdk.Iam_v1.CreateServiceUserWithResponse(ctx,
		c.Workspace, iam.IamServiceUserCreate{
			Name:        serviceUserName,
			Description: &description,
		})
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 201 {
		return res.JSON201, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) DeleteServiceUser(ctx context.Context, serviceUserUUID *uuid.UUID) error {
	_, err := c.sotoonSdk.Iam_v1.DeleteServiceUserWithResponse(ctx, c.Workspace, serviceUserUUID.String())
	return err
}

func (c *Client) UpdateServiceUser(ctx context.Context, serviceUserUUID uuid.UUID, name, description string) (*iam.IamServiceUser, error) {
	res, err := c.sotoonSdk.Iam_v1.UpdateServiceUserWithResponse(ctx, c.Workspace,
		serviceUserUUID.String(),
		iam.IamServiceUser{
			Name:        name,
			Description: description,
		})
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetWorkspaceServiceUserTokenList(ctx context.Context, serviceUserUUID, workspaceUUID *uuid.UUID) (*[]iam.IamServiceUserToken, error) {
	res, err := c.sotoonSdk.Iam_v1.ListServiceUserTokensWithResponse(ctx, workspaceUUID.String(), serviceUserUUID.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 200 {
		return res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) CreateServiceUserToken(ctx context.Context, serviceUserUUID *uuid.UUID, name string, expiresAt *time.Time) (*iam.IamServiceUserTokenWithSecret, error) {
	res, err := c.sotoonSdk.Iam_v1.CreateServiceUserTokenWithResponse(ctx, c.Workspace, serviceUserUUID.String(),
		iam.IamServiceUserTokenWithSecret{
			Name:      name,
			ExpiresAt: expiresAt,
		},
	)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 201 {
		return res.JSON201, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) DeleteServiceUserToken(ctx context.Context, serviceUserUUID, serviceUserTokenUUID *uuid.UUID) error {
	_, err := c.sotoonSdk.Iam_v1.DeleteServiceUserTokenWithResponse(ctx, c.Workspace, serviceUserUUID.String(), serviceUserTokenUUID.String())
	return err
}

func (c *Client) GetWorkspaceServiceUserPublicKeyList(ctx context.Context, workspaceUUID, serviceUserUUID uuid.UUID) ([]iam.IamServiceUserPublicKey, error) {
	res, err := c.sotoonSdk.Iam_v1.ListServiceUserPublicKeysWithResponse(ctx,
		workspaceUUID.String(),
		serviceUserUUID.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return *res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) CreateServiceUserPublicKey(ctx context.Context, serviceUserUUID uuid.UUID, name, publicKey string) (*iam.IamServiceUserPublicKey, error) {
	res, err := c.sotoonSdk.Iam_v1.CreateServiceUserPublicKeyWithResponse(ctx, c.Workspace, serviceUserUUID.String(),
		iam.IamServiceUserPublicKeyCreate{
			Key:   publicKey,
			Title: name,
		})
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 201 {
		return res.JSON201, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) DeleteServiceUserPublicKey(ctx context.Context, serviceUserUUID, publicKeyUUID uuid.UUID) error {
	_, err := c.sotoonSdk.Iam_v1.DeleteServiceUserPublicKey(ctx, c.Workspace, serviceUserUUID.String(), publicKeyUUID.String())
	return err
}

func (c *Client) UnbindRoleFromServiceUser(ctx context.Context, roleUUID, serviceUserUUID *uuid.UUID) error {
	_, err := c.sotoonSdk.Iam_v1.RemoveRoleFromServiceUserWithResponse(ctx, c.Workspace, roleUUID.String(), serviceUserUUID.String())
	return err
}

func (c *Client) GetRoleServiceUsers(ctx context.Context, roleUUID *uuid.UUID) ([]iam.IamServiceUserWithRoleItems, error) {
	res, err := c.sotoonSdk.Iam_v1.ListRolesServiceUsersWithResponse(ctx, c.Workspace, roleUUID.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return *res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) BulkAddServiceUsersToRole(ctx context.Context, roleUUID uuid.UUID, serviceUserUUIDs []string, items map[string]any) error {
	var itemsString *[]map[string]string
	if items != nil {
		converted := convertMapAnyToString(items)
		temp := make([]map[string]string, len(serviceUserUUIDs))
		for index := range serviceUserUUIDs {
			temp[index] = converted
		}
		itemsString = &temp
	}
	_, err := c.sotoonSdk.Iam_v1.BulkAddServiceUsersToRoleWithResponse(ctx, c.Workspace, roleUUID.String(),
		iam.IamBulkAddServiceUsersToRoleRequest{
			ServiceUsers: serviceUserUUIDs,
			Items:        itemsString,
		})
	return err
}

// --- IAM Role Functions ---

func (c *Client) GetWorkspaceRoles(ctx context.Context, worksapceUUID string) ([]iam.IamRole, error) {
	res, err := c.sotoonSdk.Iam_v1.ListRolesWithResponse(ctx, worksapceUUID, nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return *res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) CreateRole(ctx context.Context, name, description string) (*iam.IamMinimalRoleWithTime, error) {
	res, err := c.sotoonSdk.Iam_v1.CreateRoleWithResponse(ctx, c.Workspace,
		iam.IamCreateRole{
			Name:          name,
			DescriptionEn: description,
			Workspace:     c.Workspace,
		},
	)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 201 {
		return res.JSON201, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetRole(ctx context.Context, roleUUID *uuid.UUID) (*iam.IamRole, error) {
	res, err := c.sotoonSdk.Iam_v1.GetRoleWithResponse(ctx, c.Workspace, roleUUID.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) GetRoleByName(ctx context.Context, roleName string) (*iam.IamRole, error) {
	res, err := c.sotoonSdk.Iam_v1.ListRolesWithResponse(ctx, c.Workspace, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 200 {
		for _, role := range *res.JSON200 {
			if role.Name == roleName {
				return &role, nil
			}
		}
		return nil, ErrNotFound
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound

}

func (c *Client) DeleteRole(ctx context.Context, roleID string) error {
	_, err := c.sotoonSdk.Iam_v1.DeleteRoleWithResponse(ctx, c.Workspace, roleID)
	return err
}

func (c *Client) BulkAddRulesToRole(ctx context.Context, roleUUID uuid.UUID, ruleUUIDs []string) error {
	_, err := c.sotoonSdk.Iam_v1.BulkAddRulesToRoleWithResponse(
		ctx, c.Workspace, roleUUID.String(),
		iam.IamBulkAddRulesRequest{
			RulesUuidList: ruleUUIDs,
		})
	return err
}

func (c *Client) GetRoleRules(ctx context.Context, roleUUID *uuid.UUID) ([]iam.IamRule, error) {
	res, err := c.sotoonSdk.Iam_v1.ListRoleRulesWithResponse(ctx, c.Workspace, roleUUID.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return *res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

func (c *Client) UnbindRuleFromRole(ctx context.Context, roleUUID *uuid.UUID, ruleUUID *uuid.UUID) error {
	_, err := c.sotoonSdk.Iam_v1.RemoveRuleFromRoleWithResponse(ctx, c.Workspace, roleUUID.String(), ruleUUID.String())
	return err
}

func (c *Client) GetRoleUsers(ctx context.Context, roleUUID *uuid.UUID) ([]iam.IamUserWithRoleItems, error) {
	res, err := c.sotoonSdk.Iam_v1.ListRoleUsersWithResponse(ctx, c.Workspace, roleUUID.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return *res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}

// convertMapAnyToString converts a map[string]any to map[string]string
// by converting each value to its string representation using fmt.Sprintf
func convertMapAnyToString(input map[string]any) map[string]string {
	result := make(map[string]string, len(input))
	for key, value := range input {
		result[key] = fmt.Sprintf("%v", value)
	}
	return result
}

func (c *Client) BulkAddUsersToRole(ctx context.Context, roleUUID uuid.UUID, uuids []string, items map[string]any) error {
	var itemsString *[]map[string]string
	if items != nil {
		converted := convertMapAnyToString(items)
		temp := make([]map[string]string, len(uuids))
		for index := range uuids {
			temp[index] = converted
		}
		itemsString = &temp
	}
	_, err := c.sotoonSdk.Iam_v1.BulkAddUsersToRoleWithResponse(ctx, c.Workspace, roleUUID.String(),
		iam.IamBulkAddUsersToRoleRequest{
			Users: uuids,
			Items: itemsString,
		})
	return err
}

func (c *Client) UnbindRoleFromUser(ctx context.Context, roleUUID *uuid.UUID, userUUID *uuid.UUID) error {
	_, err := c.sotoonSdk.Iam_v1.RemoveRoleFromUser(ctx, c.Workspace, roleUUID.String(), userUUID.String())
	return err
}

// --- IAM Rule Functions ---

func (c *Client) GetWorkspaceRules(ctx context.Context, workspace string) ([]iam.IamRule, error) {
	res, err := c.sotoonSdk.Iam_v1.ListRulesWithResponse(ctx, workspace)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == 200 {
		return *res.JSON200, nil
	}
	tflog.Warn(ctx, "this should not happen", map[string]interface{}{"statusCode": res.StatusCode()})
	return nil, ErrNotFound
}
