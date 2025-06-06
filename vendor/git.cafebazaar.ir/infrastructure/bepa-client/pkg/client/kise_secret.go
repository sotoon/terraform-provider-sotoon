package client

import (
	"net/http"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/routes"
	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	uuid "github.com/satori/go.uuid"
)

func (c *bepaClient) DeleteUserKiseSecret(kiseSecretUUID *uuid.UUID) error {
	replaceDict := map[string]string{
		userUUIDPlaceholder:       c.userUUID,
		kiseSecretUUIDPlaceholder: kiseSecretUUID.String(),
		workspaceUUIDPlaceholder:  c.defaultWorkspace,
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteKiseSecretDelete), replaceDict)

	err := c.Do(http.MethodDelete, apiURL, 0, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *bepaClient) GetAllUserKiseSecret() ([]*types.KiseSecret, error) {
	replaceDict := map[string]string{
		userUUIDPlaceholder:      c.userUUID,
		workspaceUUIDPlaceholder: c.defaultWorkspace,
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteKiseSecretGetAll), replaceDict)

	kiseSecret := []*types.KiseSecret{}
	err := c.Do(http.MethodGet, apiURL, 0, nil, &kiseSecret)
	if err != nil {
		return nil, err
	}
	return kiseSecret, nil
}

func (c *bepaClient) CreateKiseSecretForDefaultUser() (*types.KiseSecret, error) {
	replaceDict := map[string]string{
		userUUIDPlaceholder:      c.userUUID,
		workspaceUUIDPlaceholder: c.defaultWorkspace,
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteKiseSecretCreate), replaceDict)

	createdKiseSecret := &types.KiseSecret{}
	if err := c.Do(http.MethodPost, apiURL, 0, nil, createdKiseSecret); err != nil {
		return nil, err
	}
	return createdKiseSecret, nil
}
