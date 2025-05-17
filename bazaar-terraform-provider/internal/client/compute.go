package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"terraform-provider-sotoon/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ IClient[model.Compute] = &computeClient{}

const (
	ComputeCreateApiVersion = "compute.ravh.ir/v1"
	ComputeCreateKind       = "InstanceClaim"
)

type computeClient struct {
	client *baseClient
	config *config
}

func (cc *computeClient) List(ctx context.Context) ([]model.Compute, error) {
	url := cc.getUrl()

	resp, err := cc.client.get(url)

	var computes []model.Compute

	var body struct {
		Items []computeCrd `json:"items"`
	}

	if err = json.Unmarshal(resp, &body); err != nil {
		return nil, err
	}

	for _, item := range body.Items {
		computes = append(computes, item.toModel())
	}

	return computes, nil
}

func (cc *computeClient) Get(ctx context.Context, name string) (*model.Compute, error) {
	computes, err := cc.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, compute := range computes {
		if compute.Name.ValueString() == name {
			return &compute, nil
		}
	}

	return nil, errors.New("instance not found")
}

func (cc *computeClient) Create(ctx context.Context, compute *model.Compute) error {
	crd := computeToCrd(compute)
	if err := cc.client.post(cc.createUrl(), crd); err != nil {
		return err
	}

	s, err := cc.Get(ctx, crd.Name)
	if err != nil {
		return err
	}
	compute.PoweredOn = s.PoweredOn
	return nil
}

func (cc *computeClient) Update(ctx context.Context, new *model.Compute) error {
	old, err := cc.Get(ctx, new.Name.ValueString())
	if err != nil {
		return err
	}
	if err = cc.client.patch(cc.getUrlWithName(old.Name.ValueString()), computeToCrd(new), computeToCrd(old)); err != nil {
		return err
	}

	s, err := cc.Get(ctx, new.Name.ValueString())
	if err != nil {
		return err
	}

	new.PoweredOn = s.PoweredOn
	return nil
}

func (cc *computeClient) Delete(ctx context.Context, compute model.Compute) error {
	return cc.client.delete(cc.getUrlWithName(compute.Name.ValueString()))
}

func (cc *computeClient) getUrl() string {
	return fmt.Sprintf(
		"%s/machine/v1/%s/apis/compute.cafebazaar.cloud%%2Fv1alpha1/namespaces/%s/instances",
		cc.config.host,
		cc.config.zone,
		cc.config.workspaceName,
	)
}

func (cc *computeClient) createUrl() string {
	return fmt.Sprintf(
		"%s/machine/v1/%s/apis/compute.ravh.ir%%2Fv1/namespaces/%s/instanceclaims",
		cc.config.host,
		cc.config.zone,
		cc.config.workspaceName,
	)
}

func (cc *computeClient) getUrlWithName(name string) string {
	return fmt.Sprintf("%s/%s", cc.getUrl(), name)
}

type computeCrd struct {
	ApiVersion     string `json:"apiVersion"`
	Kind           string `json:"kind"`
	model.Metadata `json:"metadata"`
	Spec           computeSpec `json:"spec"`
}

func (crd *computeCrd) toModel() model.Compute {
	s := model.Compute{
		Name:       types.StringValue(crd.Metadata.Name),
		IAMEnabled: types.BoolValue(crd.Spec.IAMEnabled),
		Image:      types.StringValue(crd.Spec.Image),
		Username:   types.StringValue(crd.Spec.InitialUnixUser.Username),
		Size:       types.StringValue(crd.Spec.Size),
		PoweredOn:  types.BoolValue(crd.Spec.PoweredOn),
		Subnet:     types.StringValue(""),
		Volumes:    []model.Volume{},
		ExternalIP: types.StringValue(""),
	}

	for _, volume := range crd.Spec.Volumes {
		if volume.LocalDisk.Name != "" {
			s.Volumes = append(s.Volumes, model.Volume{
				LocalDisk: model.LocalDisk{
					DiskSize: types.Int64Value(volume.LocalDisk.DiskSize),
					Name:     types.StringValue(volume.LocalDisk.Name),
				},
			})
		}
		if volume.Pvc.Name != "" {
			s.Volumes = append(s.Volumes, model.Volume{
				Pvc: model.Pvc{
					Name: types.StringValue(volume.Pvc.Name),
				},
			})
		}
	}

	return s
}

type computeCreateCrd struct {
	ApiVersion     string `json:"apiVersion"`
	Kind           string `json:"kind"`
	model.Metadata `json:"metadata"`
	Spec           computeCreateSpec `json:"spec"`
}

func computeToCrd(compute *model.Compute) computeCreateCrd {
	crd := computeCreateCrd{
		ApiVersion: ComputeCreateApiVersion,
		Kind:       ComputeCreateKind,
		Metadata:   model.Metadata{Name: compute.Name.ValueString()},
		Spec: computeCreateSpec{
			Name:       compute.Name.ValueString(),
			IAMEnabled: compute.IAMEnabled.ValueBool(),
			Image:      compute.Image.ValueString(),
			Username:   compute.Username.ValueString(),
			Type:       compute.Size.ValueString(),
			SubnetName: compute.Subnet.ValueString(),
			Disks:      []volume{},
			LinkExternalIP: linkExternalIP{
				Name: compute.ExternalIP.ValueString(),
			},
		},
	}
	for _, v := range compute.Volumes {
		if v.LocalDisk.Name.ValueString() != "" {
			crd.Spec.Disks = append(crd.Spec.Disks, volume{
				LocalDisk: localDisk{
					DiskSize: v.LocalDisk.DiskSize.ValueInt64(),
					Name:     v.LocalDisk.Name.ValueString(),
				},
			})
		}
		if v.Pvc.Name.ValueString() != "" {
			crd.Spec.Disks = append(crd.Spec.Disks, volume{
				Pvc: pvc{
					Name: v.Pvc.Name.ValueString(),
				},
			})
		}
		if v.RemoteDisk.Name.ValueString() != "" {
			crd.Spec.Disks = append(crd.Spec.Disks, volume{
				RemoteDisk: remoteDisk{
					Name: v.RemoteDisk.Name.ValueString(),
					Size: v.RemoteDisk.Size.ValueInt64(),
					Tier: v.RemoteDisk.Tier.ValueString(),
				},
			})
		}
	}
	return crd
}

type computeSpec struct {
	IAMEnabled      bool            `json:"iamEnabled"`
	Image           string          `json:"image"`
	InitialUnixUser initialUnixUser `json:"initialUnixUser"`
	Size            string          `json:"type"`
	Volumes         []volume        `json:"volumes"`
	PoweredOn       bool            `json:"poweredOn"`
}

type computeCreateSpec struct {
	Name           string         `json:"name"`
	IAMEnabled     bool           `json:"iamEnabled"`
	Image          string         `json:"image"`
	Username       string         `json:"username"`
	Type           string         `json:"type"`
	SubnetName     string         `json:"subnetName"`
	Disks          []volume       `json:"disks"`
	LinkExternalIP linkExternalIP `json:"linkExternalIP"`
}

type initialUnixUser struct {
	Username string `json:"username"`
}

type volume struct {
	LocalDisk  localDisk  `json:"localDisk"`
	Pvc        pvc        `json:"pvc"`
	RemoteDisk remoteDisk `json:"remoteDisk"`
}

type localDisk struct {
	DiskSize int64  `json:"diskSizeGb"`
	Name     string `json:"name"`
}

type pvc struct {
	Name string `json:"claimName"`
}

type remoteDisk struct {
	Name string `json:"name"`
	Size int64  `json:"gbSize"`
	Tier string `json:"tier"`
}

type linkExternalIP struct {
	Name string `json:"name"`
}
