package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"errors"
	uuid "github.com/satori/go.uuid"
	iamclient "github.com/sotoon/iam-client/pkg/client"
	"github.com/sotoon/iam-client/pkg/types"
    "strings"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
    WorkspaceUUID  *uuid.UUID // Added to store the workspace UUID for invitations
	HTTPClient     *http.Client
	IAMClient      iamclient.Client 
}

// NewClient creates a new unified API client for both Compute and IAM.
func NewClient(host, token, workspace string) (*Client, error) {
	if host == "" || token == "" || workspace == "" {
		return nil, fmt.Errorf("host, token, and workspace must not be empty")
	}

	iam, err := iamclient.NewClient(token, "https://bepa.sotoon.ir", workspace , "" , 2)
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

// Helper function to create and send requests
func (c *Client) sendComputeRequest(ctx context.Context, method, path string, payload interface{}) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonPayload)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.ComputeBaseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.HTTPClient.Do(req)
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

func (c *Client) CreateGroup(ctx context.Context, name, description string) (*types.Group, error) {
	group, err := c.IAMClient.CreateGroup(name, c.WorkspaceUUID)
	if err != nil {
		return nil, err
	}

	return &types.Group{
		UUID: group.UUID,
		Name: group.Name,
	}, nil
}

func (c *Client) DeleteGroup(ctx context.Context, groupID string) error {
	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		return fmt.Errorf("invalid group ID format: %w", err)
	}
	return c.IAMClient.DeleteGroup(c.WorkspaceUUID, &groupUUID)
}

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
