package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// API response structs
type ObjectMetadata struct {
    UID       string `json:"uid,omitempty"`
    Name      string `json:"name"`
    Workspace string `json:"workspace"`
}

type ExternalIPSpec struct {
    Reserved bool `json:"reserved"`
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
		IP       string `json:"ip"`
	} `json:"spec"`
}

type LinkSpec struct {
    ExternalIPRef struct {
        Name string `json:"name"`
    } `json:"externalIPRef"`
    SubnetName string `json:"subnetName"`
    VPCName    string `json:"vpcName"`
}

type Link struct {
    Metadata ObjectMetadata `json:"metadata"`
    Spec     LinkSpec       `json:"spec"`
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

// Image data source struct
type Image struct {
    Name        string `json:"name"`
    Version     string `json:"version"`
    OsType      string `json:"osType"`
    Description string `json:"description"`
}

// Client struct
type Client struct {
    BaseURL    string
    APIToken   string
    Workspace  string
    HTTPClient *http.Client
}

// NewClient creates a new API client
func NewClient(host, token, workspace string) (*Client, error) {
    return &Client{
        BaseURL:    fmt.Sprintf("%s/compute/v2/thr1/workspaces/%s", host, workspace),
        APIToken:   token,
        Workspace:  workspace,
        HTTPClient: &http.Client{Timeout: 10 * time.Second},
    }, nil
}
type InstanceList struct {
	Items []Instance `json:"items"`
}

// Helper function to create and send requests
func (c *Client) sendRequest(method, path string, payload interface{}) (*http.Response, error) {
    var body io.Reader
    if payload != nil {
        jsonPayload, err := json.Marshal(payload)
        if err != nil {
            return nil, err
        }
        body = bytes.NewBuffer(jsonPayload)
    }

    req, err := http.NewRequest(method, c.BaseURL+path, body)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "Bearer "+c.APIToken)
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")

    return c.HTTPClient.Do(req)
}

// CRUD for ExternalIP
func (c *Client) CreateExternalIP(name string) (*ExternalIP, error) {
    payload := map[string]interface{}{
        "apiVersion": "compute/v2",
        "kind":       "ExternalIP",
        "metadata":   map[string]string{"name": name, "workspace": c.Workspace},
        "spec":       map[string]bool{"reserved": true},
    }

    resp, err := c.sendRequest(http.MethodPost, "/external-ips", payload)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return nil, fmt.Errorf("failed to create external IP: %s", resp.Status)
    }

    var ip ExternalIP
    if err := json.NewDecoder(resp.Body).Decode(&ip); err != nil {
        return nil, err
    }
    return &ip, nil
}

func (c *Client) GetExternalIP(name string) (*ExternalIP, error) {
    resp, err := c.sendRequest(http.MethodGet, "/external-ips/"+name, nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusNotFound {
        return nil, nil // Not found is not an error for Read operations
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

func (c *Client) DeleteExternalIP(name string) error {
    resp, err := c.sendRequest(http.MethodDelete, "/external-ips/"+name, nil)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to delete external IP: %s", resp.Status)
    }
    return nil
}

// CRUD for Instance
func (c *Client) CreateInstance(name string, spec InstanceSpec) (*Instance, error) {
	payload := map[string]interface{}{
		"apiVersion": "compute/v2",
		"kind":       "Instance",
		"metadata":   map[string]string{"name": name, "workspace": c.Workspace},
		"spec":       spec,
	}

	resp, err := c.sendRequest(http.MethodPost, "/instances", payload)
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

// GetInstance retrieves a specific instance by its name (which is its resourceId).
func (c *Client) GetInstance(name string) (*Instance, error) {
	resp, err := c.sendRequest(http.MethodGet, "/instances/"+name, nil)
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

// UpdateInstance updates an existing instance.
func (c *Client) UpdateInstance(name string, spec InstanceSpec) (*Instance, error) {
    payload := map[string]interface{}{
		"apiVersion": "compute/v2",
		"kind":       "Instance",
		"metadata":   map[string]string{"name": name, "workspace": c.Workspace},
		"spec":       spec,
	}
	resp, err := c.sendRequest(http.MethodPut, "/instances/"+name, payload)
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


// DeleteInstance deletes a specific instance by its name.
func (c *Client) DeleteInstance(name string) error {
	resp, err := c.sendRequest(http.MethodDelete, "/instances/"+name, nil)
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

// List Images
func (c *Client) ListImages() ([]Image, error) {
    resp, err := c.sendRequest(http.MethodGet, "/images", nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to list images: %s", resp.Status)
    }

    var images []Image
    if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
        return nil, err
    }
    return images, nil
}

func (c *Client) ListInstances() ([]Instance, error) {
	resp, err := c.sendRequest(http.MethodGet, "/instances", nil)
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

func (c *Client) ListExternalIPs() ([]ExternalIP, error) {
	apiURL := fmt.Sprintf("%s//external-ips", c.BaseURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)

	resp, err := c.HTTPClient.Do(req)
    
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode  )
	}

	var apiResponse struct {
		Items []ExternalIP `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}
    
	return apiResponse.Items, nil
}