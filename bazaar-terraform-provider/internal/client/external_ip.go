package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-sotoon/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ IClient[model.ExternalIP] = &externalIPClient{}

type externalIPClient struct {
	*baseClient
	config *config
}

func (e *externalIPClient) List(ctx context.Context) ([]model.ExternalIP, error) {
	url := fmt.Sprintf(
		"%s/machine/v1/%s/apis/networking.cafebazaar.cloud%%2Fv1alpha1/namespaces/%s/externalips",
		e.config.host,
		e.config.zone,
		e.config.workspaceName,
	)

	ctx = tflog.SetField(ctx, "url", url)
	tflog.Info(ctx, "externalIP url")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating GET request:", err)
		return nil, err
	}

	resp, err := e.doRequest(req)
	if err != nil {
		fmt.Println("Error executing GET request:", err)
		return nil, err
	}

	var externalIps []model.ExternalIP

	var body struct {
		Items []struct {
			model.Metadata       `json:"metadata"`
			model.ExternalIPSpec `json:"spec"`
		} `json:"items"`
	}

	if err = json.Unmarshal(resp, &body); err != nil {
		return nil, err
	}

	for _, externalIp := range body.Items {
		ne := model.ExternalIP{
			Name:      types.StringValue(externalIp.Metadata.Name),
			GatewayIP: types.StringValue(externalIp.ExternalIPSpec.Gateway),
			IP:        types.StringValue(externalIp.ExternalIPSpec.IP),
			Reserved:  types.BoolValue(externalIp.ExternalIPSpec.Reserved),
		}

		externalIps = append(externalIps, ne)
	}

	return externalIps, nil

}

func (e *externalIPClient) Create(ctx context.Context, t *model.ExternalIP) error {
	url := fmt.Sprintf(
		"%s/machine/v1/%s/apis/networking.cafebazaar.cloud%%2Fv1alpha1/namespaces/%s/externalips",
		e.config.host,
		e.config.zone,
		e.config.workspaceName,
	)

	ctx = tflog.SetField(ctx, "url", url)
	tflog.Info(ctx, "externalIP url")

	body := struct {
		ApiVersion           string `json:"apiVersion"`
		Kind                 string `json:"kind"`
		model.Metadata       `json:"metadata"`
		model.ExternalIPSpec `json:"spec"`
	}{
		ApiVersion: "networking.cafebazaar.cloud/v1alpha1",
		Kind:       "ExternalIP",
		Metadata: model.Metadata{
			Name: t.Name.ValueString(),
		},
		ExternalIPSpec: model.ExternalIPSpec{
			Reserved: t.Reserved.ValueBool(),
		},
	}

	// Serialize the ComputeSpec struct to JSON
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		return err
	}

	_, err = e.doRequest(req)
	if err != nil {
		fmt.Println("Error executing POST request:", err)
		return err
	}
	return nil
}

func (e *externalIPClient) Get(ctx context.Context, name string) (*model.ExternalIP, error) {
	panic("")
}

func (e *externalIPClient) Update(ctx context.Context, t *model.ExternalIP) error {
	//TODO implement me
	panic("implement me")
}

func (e *externalIPClient) Delete(ctx context.Context, t model.ExternalIP) error {
	//TODO implement me
	panic("implement me")
}
