package client

import (
	"terraform-provider-sotoon/internal/model"
)

type Client struct {
	Subnet     IClient[model.Subnet]
	ExternalIP IClient[model.ExternalIP]
	PVC        IClient[model.PVC]
	Compute    IClient[model.Compute]
}

func NewClient(host, workspaceName, workspaceID, zone, token string) *Client {
	config := &config{
		host:          host,
		workspaceName: workspaceName,
		workspaceID:   workspaceID,
		zone:          zone,
	}
	baseClient := newBaseClient(token)
	return &Client{
		Subnet: &subnetClient{
			client: baseClient,
			config: config,
		},
		ExternalIP: &externalIPClient{
			baseClient: baseClient,
			config:     config,
		},
		PVC: &pvcClient{
			client: baseClient,
			config: config,
		},
		Compute: &computeClient{
			client: baseClient,
			config: config,
		},
	}
}
