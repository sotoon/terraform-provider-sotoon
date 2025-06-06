package client

import (
	"errors"
	"fmt"
	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/routes"
	"net/http"
	"net/url"
)

func checkResponse(response *http.Response) (bool, error) {
	if response.StatusCode != http.StatusOK {
		return false, errors.New("invalid response code")
	}
	// todo: check healthz response keys and values
	return true, nil
}

func healthCheck(c *bepaClient, serverUrl *url.URL) error {
	client := &http.Client{Timeout: c.bepaTimeout}
	healthCheckPath, err := url.Parse(trimURLSlash(routes.RouteHealthCheck))
	if err != nil {
		return err
	}
	healthCheckFullPath := serverUrl.ResolveReference(healthCheckPath)
	httpRequest, err := http.NewRequest(http.MethodGet, healthCheckFullPath.String(), nil)
	if err != nil {
		return err
	}
	httpRequest.Header.Add("authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()

	ok, err := checkResponse(httpResponse)
	if err != nil || ok != true {
		return err
	}
	return nil
}
