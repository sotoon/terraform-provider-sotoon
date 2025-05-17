package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/snorwin/jsonpatch"
	"io/ioutil"
	"net/http"
	"time"
)

type baseClient struct {
	httpClient *http.Client
	token      string
}

func newBaseClient(token string) *baseClient {
	return &baseClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		token:      token,
	}
}

func (bc *baseClient) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bc.token))

	res, err := bc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (bc *baseClient) get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return bc.doRequest(req)
}

func (bc *baseClient) post(url string, body any) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	_, err = bc.doRequest(req)
	return err
}

func (bc *baseClient) patch(url string, new, old any) error {
	patch, err := jsonpatch.CreateJSONPatch(new, old)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(patch.Raw()))
	if err != nil {
		return err
	}

	req.Header.Set("content-type", "application/json-patch+json")

	_, err = bc.doRequest(req)
	return err
}

func (bc *baseClient) delete(url string) error {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	_, err = bc.doRequest(req)
	return err
}
