package client

import (
	"context"
	"terraform-provider-sotoon/internal/model"
)

type resource interface {
	model.Subnet | model.Compute | model.ExternalIP | model.PVC
}

type IClient[T resource] interface {
	List(context.Context) ([]T, error)
	Get(context.Context, string) (*T, error)
	Create(context.Context, *T) error
	Update(context.Context, *T) error
	Delete(context.Context, T) error
}
