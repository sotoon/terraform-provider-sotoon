package client

import (
	"errors"
	"net/http"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/routes"
	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	uuid "github.com/satori/go.uuid"
)

func (c *bepaClient) GetWorkspaces() ([]*types.Workspace, error) {
	apiURL := trimURLSlash(routes.RouteWorkspaceGetAll)

	workspaces := []*types.Workspace{}
	err := c.Do(http.MethodGet, apiURL, 0, nil, &workspaces)
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

// todo deprecate and remove this functoin
func (c *bepaClient) GetWorkspaceByName(name string) (*types.Workspace, error) {
	replaceDict := map[string]string{
		workspaceNamePlaceholder: name,
		userUUIDPlaceholder:      c.userUUID,
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteUserGetOneWorkspaceByName), replaceDict)

	workspace := &types.Workspace{}
	err := c.Do(http.MethodGet, apiURL, 0, nil, &workspace)
	if err != nil {
		return nil, err
	}
	return workspace, nil
}

func (c *bepaClient) GetWorkspaceByNameAndOrgName(name string, organizationName string) (*types.WorkspaceWithOrganization, error) {
	replaceDict := map[string]string{
		userUUIDPlaceholder: c.userUUID,
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteUserGetOneWorkspace), replaceDict)

	parameters := map[string]string{
		"name":     name,
		"org_name": organizationName,
	}

	var workspacesSingleArray []types.WorkspaceWithOrganization

	err := c.DoWithParams(http.MethodGet, apiURL, parameters, 0, nil, &workspacesSingleArray)
	if err != nil {
		return nil, err
	}
	if len(workspacesSingleArray) == 1 {
		return &workspacesSingleArray[0], nil
	} else {
		return nil, errors.New("No workspace found")
	}
}

func (c *bepaClient) GetWorkspace(uuid *uuid.UUID) (*types.Workspace, error) {
	replaceDict := map[string]string{
		workspaceUUIDPlaceholder: uuid.String(),
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteWorkspaceGetOne), replaceDict)

	workspace := &types.Workspace{}
	err := c.Do(http.MethodGet, apiURL, 0, nil, &workspace)
	if err != nil {
		return nil, err
	}
	return workspace, nil
}

func (c *bepaClient) GetMyWorkspaces() ([]*types.WorkspaceWithOrganization, error) {
	replaceDict := map[string]string{
		userUUIDPlaceholder: c.userUUID,
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteUserGetAllWorkspaces), replaceDict)

	workspaces := []*types.WorkspaceWithOrganization{}
	err := c.Do(http.MethodGet, apiURL, 0, nil, &workspaces)
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (c *bepaClient) GetWorkspaceUsers(uuid *uuid.UUID) ([]*types.User, error) {
	replaceDict := map[string]string{
		workspaceUUIDPlaceholder: uuid.String(),
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteWorkspaceGetUsers), replaceDict)

	users := []*types.User{}
	err := c.Do(http.MethodGet, apiURL, 0, nil, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (c *bepaClient) CreateWorkspace(name string) (*types.Workspace, error) {
	workspaceRequest := &types.WorkspaceReq{
		Name: name,
	}

	createdWorkspace := &types.Workspace{}
	apiURL := trimURLSlash(routes.RouteWorkspaceCreate)
	err := c.Do(http.MethodPost, apiURL, 0, workspaceRequest, &createdWorkspace)
	if err != nil {
		return nil, err
	}
	return createdWorkspace, nil
}

func (c *bepaClient) GetWorkspaceRules(uuid *uuid.UUID) ([]*types.Rule, error) {
	replaceDict := map[string]string{
		workspaceUUIDPlaceholder: uuid.String(),
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteWorkspaceGetAllRules), replaceDict)
	rules := []*types.Rule{}
	err := c.Do(http.MethodGet, apiURL, 0, nil, &rules)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

func (c *bepaClient) GetWorkspaceRoles(uuid *uuid.UUID) ([]*types.Role, error) {
	replaceDict := map[string]string{
		workspaceUUIDPlaceholder: uuid.String(),
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteWorkspaceGetAllRoles), replaceDict)
	roles := []*types.Role{}
	err := c.Do(http.MethodGet, apiURL, 0, nil, &roles)
	if err != nil {
		return nil, err
	}
	return roles, nil

}

func (c *bepaClient) DeleteWorkspace(uuid *uuid.UUID) error {
	replaceDict := map[string]string{
		workspaceUUIDPlaceholder: uuid.String(),
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteWorkspaceDelete), replaceDict)

	return c.Do(http.MethodDelete, apiURL, 0, nil, nil)
}
