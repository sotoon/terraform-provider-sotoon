package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"terraform-provider-sotoon/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ IClient[model.PVC] = &pvcClient{}

const (
	PvcApiVersion = "v1"
	PvcKind       = "PersistentVolumeClaim"
)

type pvcClient struct {
	client *baseClient
	config *config
}

func (p *pvcClient) List(ctx context.Context) ([]model.PVC, error) {
	resp, err := p.client.get(p.getUrl())
	if err != nil {
		return nil, err
	}

	var pvcs []model.PVC

	var body struct {
		Items []pvcCrd `json:"items"`
	}

	if err = json.Unmarshal(resp, &body); err != nil {
		return nil, err
	}

	for _, pvc := range body.Items {
		pvcs = append(pvcs, pvc.toModel())
	}

	return pvcs, nil

}

func (p *pvcClient) Get(ctx context.Context, name string) (*model.PVC, error) {
	pvcList, err := p.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, pvc := range pvcList {
		if pvc.Name.ValueString() == name {
			return &pvc, nil
		}
	}

	return nil, errors.New("pvc not found")
}

func (p *pvcClient) Create(ctx context.Context, pvc *model.PVC) error {
	crd := pvcToCrd(pvc)
	return p.client.post(p.getUrl(), crd)
}

func (p *pvcClient) Update(ctx context.Context, new *model.PVC) error {
	old, err := p.Get(ctx, new.Name.ValueString())
	if err != nil {
		return err
	}
	return p.client.patch(p.getUrlWithName(old.Name.ValueString()), pvcToCrd(new), pvcToCrd(old))
}

func (p *pvcClient) Delete(ctx context.Context, pvc model.PVC) error {
	return p.client.delete(p.getUrlWithName(pvc.Name.ValueString()))
}

func (p *pvcClient) getUrl() string {
	return fmt.Sprintf(
		"%s/machine/v1/%s/api%%2Fv1/namespaces/%s/persistentvolumeclaims",
		p.config.host,
		p.config.zone,
		p.config.workspaceName,
	)
}

func (p *pvcClient) getUrlWithName(name string) string {
	return fmt.Sprintf("%s/%s", p.getUrl(), name)
}

type pvcCrd struct {
	ApiVersion string      `json:"apiVersion"`
	Kind       string      `json:"kind"`
	Metadata   pvcMetadata `json:"metadata"`
	Spec       pvcSpec     `json:"spec"`
	Status     pvcStatus   `json:"status"`
}

func (p *pvcCrd) toModel() model.PVC {
	return model.PVC{
		Name:        types.StringValue(p.Metadata.Name),
		Size:        types.StringValue(p.Status.Capacity.Storage),
		StorageTier: types.StringValue(p.Metadata.Annotations.Tier),
	}
}

func pvcToCrd(pvc *model.PVC) pvcCrd {
	return pvcCrd{
		ApiVersion: PvcApiVersion,
		Kind:       PvcKind,
		Metadata: pvcMetadata{
			Name: pvc.Name.ValueString(),
			Annotations: pvcAnnotation{
				Tier: pvc.StorageTier.ValueString(),
			},
		},
		Spec: pvcSpec{
			Resource: pvcSpecResource{
				Requests: pvcSpecResourceRequest{
					Storage: pvc.Size.ValueString(),
				},
			},
			StorageClassName: "general-purpose",
			AccessModes:      []string{"ReadWriteOnce"},
			VolumeMode:       "Block",
		},
	}
}

type pvcStatus struct {
	Capacity struct {
		Storage string `json:"storage"`
	} `json:"Capacity"`
}

type pvcSpec struct {
	Resource         pvcSpecResource `json:"resources"`
	StorageClassName string          `json:"storageClassName"`
	AccessModes      []string        `json:"accessModes"`
	VolumeMode       string          `json:"volumeMode"`
}

type pvcSpecResource struct {
	Requests pvcSpecResourceRequest `json:"requests"`
}

type pvcSpecResourceRequest struct {
	Storage string `json:"storage"`
}

type pvcMetadata struct {
	Name        string        `json:"name"`
	Annotations pvcAnnotation `json:"annotations"`
}

type pvcAnnotation struct {
	Tier string `json:"compute.cafebazaar.cloud/storage-tier"`
}
