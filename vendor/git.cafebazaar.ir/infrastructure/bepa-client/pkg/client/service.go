package client

import (
	"net/http"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/routes"
	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
)

func (c *bepaClient) GetService(name string) (*types.Service, error) {

	replaceDict := map[string]string{
		serviceNamePlaceholder: name,
	}

	service := &types.Service{}
	apiURL := substringReplace(trimURLSlash(routes.RouteServiceGetOne), replaceDict)
	err := c.Do(http.MethodGet, apiURL, 0, nil, service)
	return service, err
}
func (c *bepaClient) GetAllServices() (*[]types.Service, error) {

	services := &[]types.Service{}
	apiURL := substringReplace(trimURLSlash(routes.RouteServiceGetAll), nil)
	err := c.Do(http.MethodGet, apiURL, 0, nil, services)
	return services, err
}
