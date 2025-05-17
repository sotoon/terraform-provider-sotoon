package types

import "fmt"

type ResponseError struct {
	Error    string   `json:"message"`
	Invalids []string `json:"invalids,omitempty"`
}

type RequestExecutionError struct {
	Err        error
	StatusCode int
	Data       []byte
}

func (ree *RequestExecutionError) Error() string {
	return ree.Err.Error()
}

type HealthCheckResponse struct {
	ServerUrl string
	Err       error
}

func (hcr *HealthCheckResponse) String() string {
	return fmt.Sprintf("url: %s, error: %s", hcr.ServerUrl, hcr.Err)
}
