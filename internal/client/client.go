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

type ExternalIP struct {
	Metadata struct {
		Name      string `json:"name"`
		UID       string `json:"uid"`
		Workspace string `json:"workspace"`
	} `json:"metadata"`
	Spec struct {
		BoundTo *struct {
			Kind string `json:"kind"`
			Name string `json:"name"`
		} `json:"boundTo"`
		IP string `json:"ip"`
	} `json:"spec"`
}

type InstanceSpec struct {
	IAMEnabled  bool `json:"iamEnabled"`
	ImageSource struct {
		Image string `json:"image"`
	} `json:"imageSource"`
	Interfaces []struct {
		Link string `json:"link"`
		Name string `json:"name"`
	} `json:"interfaces"`
	PoweredOn bool   `json:"poweredOn"`
	Type      string `json:"type"`
}

type Instance struct {
	Metadata ObjectMetadata `json:"metadata"`
	Spec     InstanceSpec   `json:"spec"`
}

type InstanceList struct {
	Items []Instance `json:"items"`
}

type Image struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	OsType      string `json:"osType"`
	Description string `json:"description"`
}

type CreateUserRequest struct {
	Email    string
	Name     string
	Password string
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

	iam, err := iamclient.NewClient(token, "https://bepa.sotoon.ir", workspace , "")
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

// --- Compute Functions ---

func (c *Client) CreateInstance(ctx context.Context, name string, spec InstanceSpec) (*Instance, error) {
	// payload := map[string]interface{}{"metadata": map[string]string{"name": name}, "spec": spec}
	payload := map[string]interface{}{
		"apiVersion": "compute/v2",
		"kind":       "Instance",
		"metadata":   map[string]string{"name": name, "workspace": c.Workspace},
		"spec":       spec,
	}

	resp, err := c.sendComputeRequest(ctx, http.MethodPost, "/instances", payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create instance (status: %s): %s", resp.Status, string(bodyBytes))
	}

	var instance Instance
	if err := json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return nil, err
	}
	return &instance, nil
}

func (c *Client) GetInstance(ctx context.Context, name string) (*Instance, error) {
	resp, err := c.sendComputeRequest(ctx, http.MethodGet, "/instances/"+name, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get instance (status: %s): %s", resp.Status, string(bodyBytes))
	}
	var instance Instance
	if err := json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return nil, err
	}
	return &instance, nil
}

func (c *Client) UpdateInstance(ctx context.Context, name string, spec InstanceSpec) (*Instance, error) {
	payload := map[string]interface{}{"metadata": map[string]string{"name": name}, "spec": spec}
	resp, err := c.sendComputeRequest(ctx, http.MethodPut, "/instances/"+name, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update instance (status: %s): %s", resp.Status, string(bodyBytes))
	}
	var instance Instance
	if err := json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return nil, err
	}
	return &instance, nil
}

func (c *Client) DeleteInstance(ctx context.Context, name string) error {
	resp, err := c.sendComputeRequest(ctx, http.MethodDelete, "/instances/"+name, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete instance (status: %s): %s", resp.Status, string(bodyBytes))
	}
	return nil
}

func (c *Client) ListImages(ctx context.Context) ([]Image, error) {
	resp, err := c.sendComputeRequest(ctx, http.MethodGet, "/images", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list images: %s", resp.Status)
	}
	var imageList struct {
		Items []Image `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&imageList); err != nil {
		return nil, err
	}
	return imageList.Items, nil
}

func (c *Client) ListInstances(ctx context.Context) ([]Instance, error) {
	resp, err := c.sendComputeRequest(ctx, http.MethodGet, "/instances", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list instances (status: %s): %s", resp.Status, string(bodyBytes))
	}

	var list InstanceList
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed to decode instance list: %w", err)
	}

	return list.Items, nil
}

// --- External IP Functions ---

func (c *Client) CreateExternalIP(ctx context.Context, name string) (*ExternalIP, error) {
	payload := map[string]interface{}{
        "apiVersion": "compute/v2",
        "kind":       "ExternalIP",
        "metadata":   map[string]string{"name": name, "workspace": c.Workspace},
        "spec":       map[string]bool{"reserved": true},
    }

	resp, err := c.sendComputeRequest(ctx, http.MethodPost, "/external-ips", payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create external IP: %s ", resp.Status)
	}
	var ip ExternalIP
	if err := json.NewDecoder(resp.Body).Decode(&ip); err != nil {
		return nil, err
	}
	return &ip, nil
}

func (c *Client) GetExternalIP(ctx context.Context, name string) (*ExternalIP, error) {
	resp, err := c.sendComputeRequest(ctx, http.MethodGet, "/external-ips/"+name, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get external IP: %s", resp.Status)
	}
	var ip ExternalIP
	if err := json.NewDecoder(resp.Body).Decode(&ip); err != nil {
		return nil, err
	}
	return &ip, nil
}

func (c *Client) DeleteExternalIP(ctx context.Context, name string) error {
	resp, err := c.sendComputeRequest(ctx, http.MethodDelete, "/external-ips/"+name, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete external IP: %s", resp.Status)
	}
	return nil
}

func (c *Client) ListExternalIPs(ctx context.Context) ([]ExternalIP, error) {
	resp, err := c.sendComputeRequest(ctx, http.MethodGet, "/external-ips", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list external IPs: %s", resp.Status)
	}
	var ipList struct {
		Items []ExternalIP `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ipList); err != nil {
		return nil, err
	}
	return ipList.Items, nil
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

