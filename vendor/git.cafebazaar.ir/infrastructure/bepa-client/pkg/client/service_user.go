package client

import (
	"net/http"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/routes"
	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	uuid "github.com/satori/go.uuid"
)

func (c *bepaClient) GetServiceUser(workspaceUUID, serviceUserUUID *uuid.UUID) (*types.ServiceUser, error) {
	replaceDict := map[string]string{
		serviceUserUUIDPlaceholder: serviceUserUUID.String(),
		workspaceUUIDPlaceholder:   workspaceUUID.String(),
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteServiceUserGetOne), replaceDict)

	service := &types.ServiceUser{}
	if err := c.Do(http.MethodGet, apiURL, 0, nil, service); err != nil {
		return nil, err
	}
	return service, nil
}

func (c *bepaClient) GetServiceUsers(workspaceUUID *uuid.UUID) ([]*types.ServiceUser, error) {
	replaceDict := map[string]string{
		workspaceUUIDPlaceholder: workspaceUUID.String(),
	}
	serviceUsers := []*types.ServiceUser{}
	apiURL := substringReplace(trimURLSlash(routes.RouteServiceUserGetALL), replaceDict)
	if err := c.Do(http.MethodGet, apiURL, 0, nil, &serviceUsers); err != nil {
		return nil, err
	}
	return serviceUsers, nil
}

func (c *bepaClient) DeleteServiceUser(workspaceUUID, serviceUserUUID *uuid.UUID) error {
	replaceDict := map[string]string{
		serviceUserUUIDPlaceholder: serviceUserUUID.String(),
		workspaceUUIDPlaceholder:   workspaceUUID.String(),
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteServiceUserDelete), replaceDict)
	return c.Do(http.MethodDelete, apiURL, 0, nil, nil)
}

func (c *bepaClient) GetServiceUserByName(workspaceName string, serviceUserName string) (*types.ServiceUser, error) {
	replaceDict := map[string]string{
		serviceUserNamePlaceholder: serviceUserName,
		workspaceNamePlaceholder:   workspaceName,
		userUUIDPlaceholder:        c.userUUID,
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteServiceUserGetByName), replaceDict)

	serviceUser := &types.ServiceUser{}
	if err := c.Do(http.MethodGet, apiURL, 0, nil, serviceUser); err != nil {
		return nil, err
	}
	return serviceUser, nil
}
func (c *bepaClient) CreateServiceUser(serviceUserName string, workspace *uuid.UUID) (*types.ServiceUser, error) {
	userRequest := &types.ServiceUserReq{
		Name:      serviceUserName,
		Workspace: workspace.String(),
	}
	replaceDict := map[string]string{
		workspaceUUIDPlaceholder: workspace.String(),
	}
	createdServiceUser := &types.ServiceUser{}
	apiURL := substringReplace(trimURLSlash(routes.RouteServiceUserCreate), replaceDict)
	if err := c.Do(http.MethodPost, apiURL, 0, userRequest, createdServiceUser); err != nil {
		return nil, err
	}
	return createdServiceUser, nil
}

func (c *bepaClient) CreateServiceUserToken(serviceUserUUID, workspaceUUID *uuid.UUID) (*types.ServiceUserToken, error) {
	replaceDict := map[string]string{
		serviceUserUUIDPlaceholder: serviceUserUUID.String(),
		workspaceUUIDPlaceholder:   workspaceUUID.String(),
	}
	ServiceUserToken := &types.ServiceUserToken{}
	apiURL := substringReplace(trimURLSlash(routes.RouteServiceUserTokenCreate), replaceDict)
	err := c.Do(http.MethodPost, apiURL, 0, nil, ServiceUserToken)
	return ServiceUserToken, err
}

func (c *bepaClient) GetAllServiceUserToken(serviceUserUUID, workspaceUUID *uuid.UUID) (*[]types.ServiceUserToken, error) {

	replaceDict := map[string]string{
		serviceUserUUIDPlaceholder: serviceUserUUID.String(),
		workspaceUUIDPlaceholder:   workspaceUUID.String(),
	}
	ServiceUserTokens := &[]types.ServiceUserToken{}
	apiURL := substringReplace(trimURLSlash(routes.RouteServiceUserTokenGetALL), replaceDict)
	err := c.Do(http.MethodGet, apiURL, 0, nil, ServiceUserTokens)
	return ServiceUserTokens, err
}

func (c *bepaClient) DeleteServiceUserToken(serviceUserUUID, workspaceUUID, serviceUserTokenUUID *uuid.UUID) error {

	replaceDict := map[string]string{
		serviceUserUUIDPlaceholder:      serviceUserUUID.String(),
		workspaceUUIDPlaceholder:        workspaceUUID.String(),
		serviceUserTokenUUIDPlaceholder: serviceUserTokenUUID.String(),
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteServiceUserTokenDelete), replaceDict)
	return c.Do(http.MethodDelete, apiURL, 0, nil, nil)

}
