package client

import (
	"net/http"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/routes"
	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
)

func (c *bepaClient) Identify(token string) (*types.UserRes, error) {
	idenReq := &types.UserTokenReq{
		Secret: token,
	}

	userRes := &types.UserRes{}
	err := c.Do(http.MethodPost, trimURLSlash(routes.RouteUserTokenIdentify), 0, idenReq, userRes)
	if err != nil {
		return nil, err
	}

	return userRes, nil
}
