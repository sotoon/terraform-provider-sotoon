package client

import (
	"net/http"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/routes"
	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
)

func (c *bepaClient) Authorize(identity, userType, action, object string) error {
	c.log("authorizing %v", identity)
	req, err := c.NewRequest(http.MethodGet, trimURLSlash(routes.RouteAuthz), nil)

	if err != nil {
		return err
	}

	query := req.URL.Query()
	query.Set("identity", identity)
	query.Set("user_type", userType)
	query.Set("object", object)
	query.Set("action", action)

	req.URL.RawQuery = query.Encode()
	data, statusCode, errRes := proccessRequest(req, 0)
	if errRes == nil {
		c.log("user %v is authorized", identity)
		return nil
	}

	c.log("user %v is not authorized", identity)
	return &types.RequestExecutionError{
		Err:        errRes,
		StatusCode: statusCode,
		Data:       data,
	}
}
