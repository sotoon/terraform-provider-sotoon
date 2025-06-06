package client

import (
	"net/http"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/routes"
	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	uuid "github.com/satori/go.uuid"
)

func (c *bepaClient) CreateMyUserTokenWithToken(secret string) (*types.UserToken, error) {
	userTokenreq := &types.UserTokenReq{
		Secret: secret,
	}
	replaceDict := map[string]string{
		userUUIDPlaceholder: c.userUUID,
	}

	userToken := &types.UserToken{}
	apiURL := substringReplace(trimURLSlash(routes.RouteUserTokenCreateByToken), replaceDict)
	err := c.Do(http.MethodPost, apiURL, 0, userTokenreq, userToken)
	return userToken, err
}

func (c *bepaClient) GetMyUserToken(userTokenUUID *uuid.UUID) (*types.UserToken, error) {

	replaceDict := map[string]string{
		userUUIDPlaceholder:      c.userUUID,
		userTokenUUIDPlaceholder: userTokenUUID.String(),
	}

	userToken := &types.UserToken{}
	apiURL := substringReplace(trimURLSlash(routes.RouteUserTokenGetOne), replaceDict)
	err := c.Do(http.MethodGet, apiURL, 0, nil, userToken)
	return userToken, err
}
func (c *bepaClient) GetAllMyUserTokens() (*[]types.UserToken, error) {

	replaceDict := map[string]string{
		userUUIDPlaceholder: c.userUUID,
	}

	userTokens := &[]types.UserToken{}
	apiURL := substringReplace(trimURLSlash(routes.RouteUserTokenGetAll), replaceDict)
	err := c.Do(http.MethodGet, apiURL, 0, nil, userTokens)
	return userTokens, err
}

func (c *bepaClient) DeleteMyUserToken(userTokenUUID *uuid.UUID) error {

	replaceDict := map[string]string{
		userUUIDPlaceholder:      c.userUUID,
		userTokenUUIDPlaceholder: userTokenUUID.String(),
	}

	apiURL := substringReplace(trimURLSlash(routes.RouteUserTokenDelete), replaceDict)
	return c.Do(http.MethodDelete, apiURL, 0, nil, nil)

}
