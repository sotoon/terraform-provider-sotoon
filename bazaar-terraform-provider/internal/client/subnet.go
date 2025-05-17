package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-sotoon/internal/model"
)

var _ IClient[model.Subnet] = &subnetClient{}

const (
	SubnetApiVersion = "networking.cafebazaar.cloud/v1alpha1"
	SubnetKind       = "Subnet"
)

type subnetClient struct {
	client *baseClient
	config *config
}

func (sc *subnetClient) List(ctx context.Context) ([]model.Subnet, error) {
	url := sc.getUrl()

	resp, err := sc.client.get(url)

	var subnets []model.Subnet

	var body struct {
		Items []subnetCrd `json:"items"`
	}

	if err = json.Unmarshal(resp, &body); err != nil {
		return nil, err
	}

	for _, item := range body.Items {
		subnets = append(subnets, item.toModel())
	}

	return subnets, nil
}

func (sc *subnetClient) Get(ctx context.Context, name string) (*model.Subnet, error) {
	subnets, err := sc.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, subnet := range subnets {
		if subnet.Name.ValueString() == name {
			return &subnet, nil
		}
	}

	return nil, errors.New("subnet not found")
}

func (sc *subnetClient) Create(ctx context.Context, subnet *model.Subnet) error {
	crd := subnetToCrd(subnet)
	if err := sc.client.post(sc.getUrl(), crd); err != nil {
		return err
	}

	s, err := sc.Get(ctx, crd.Name)
	if err != nil {
		return err
	}

	subnet.GatewayIP = s.GatewayIP
	subnet.Routes = s.Routes
	return nil
}

func (sc *subnetClient) Update(ctx context.Context, new *model.Subnet) error {
	old, err := sc.Get(ctx, new.Name.ValueString())
	if err != nil {
		return err
	}
	if err = sc.client.patch(sc.getUrlWithName(old.Name.ValueString()), subnetToCrd(new), subnetToCrd(old)); err != nil {
		return err
	}

	s, err := sc.Get(ctx, new.Name.ValueString())
	if err != nil {
		return err
	}

	new.GatewayIP = s.GatewayIP
	new.Routes = s.Routes

	return nil
}

func (sc *subnetClient) Delete(ctx context.Context, subnet model.Subnet) error {
	return sc.client.delete(sc.getUrlWithName(subnet.Name.ValueString()))
}

func (sc *subnetClient) getUrl() string {
	return fmt.Sprintf(
		"%s/machine/v1/%s/apis/networking.cafebazaar.cloud%%2Fv1alpha1/namespaces/%s/subnets",
		sc.config.host,
		sc.config.zone,
		sc.config.workspaceName,
	)
}

func (sc *subnetClient) getUrlWithName(name string) string {
	return fmt.Sprintf("%s/%s", sc.getUrl(), name)
}

type subnetCrd struct {
	ApiVersion     string `json:"apiVersion"`
	Kind           string `json:"kind"`
	model.Metadata `json:"metadata"`
	Spec           subnetSpec `json:"spec"`
}

func (crd *subnetCrd) toModel() model.Subnet {
	s := model.Subnet{
		Name:      types.StringValue(crd.Metadata.Name),
		GatewayIP: types.StringValue(crd.Spec.GatewayIP),
		Cidr:      types.StringValue(crd.Spec.Cidr),
	}

	for _, route := range crd.Spec.Routes {
		s.Routes = append(s.Routes, model.SubnetRoute{
			To:         types.StringValue(route.To),
			ExternalIP: types.StringValue(route.Via.ExternalIP),
		})
	}

	return s
}

func subnetToCrd(subnet *model.Subnet) subnetCrd {
	crd := subnetCrd{
		ApiVersion: SubnetApiVersion,
		Kind:       SubnetKind,
		Metadata:   model.Metadata{Name: subnet.Name.ValueString()},
		Spec: subnetSpec{
			Cidr:   subnet.Cidr.ValueString(),
			Routes: []subnetRoute{},
		},
	}
	for _, route := range subnet.Routes {
		crd.Spec.Routes = append(crd.Spec.Routes, subnetRoute{
			To: route.To.ValueString(),
			Via: via{
				ExternalIP: route.ExternalIP.ValueString(),
			},
		})
	}
	return crd
}

type subnetSpec struct {
	Cidr      string        `json:"cidr"`
	GatewayIP string        `json:"gatewayIP"`
	Routes    []subnetRoute `json:"routes"`
}

type subnetRoute struct {
	To  string `json:"to"`
	Via via    `json:"via"`
}

type via struct {
	ExternalIP string `json:"externalIP"`
}
