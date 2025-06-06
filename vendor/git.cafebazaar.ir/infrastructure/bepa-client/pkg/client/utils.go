package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"github.com/spf13/viper"
)

var (
	ErrNotMatched          = errors.New("not matched")
	ErrForbidden           = errors.New("forbidden")
	ErrNotFound            = errors.New("not exists")
	ErrBadRequest          = errors.New("bad request")
	ErrInternalServerError = errors.New("internal server error")
)

// HTTPResponseError is a type for errors on http requests based on status code
type HTTPResponseError struct {
	StatusCode  int
	IsFaulty    bool
	internalErr error
}

func (re *HTTPResponseError) Error() string {
	return re.internalErr.Error()
}

func createHTTPResponseError(statusCode int, internalErr error) *HTTPResponseError {
	return &HTTPResponseError{
		StatusCode:  statusCode,
		IsFaulty:    false,
		internalErr: internalErr,
	}
}

func createFaultyHTTPResponseError(statusCode int, internalErr error) *HTTPResponseError {
	return &HTTPResponseError{
		StatusCode:  statusCode,
		IsFaulty:    true,
		internalErr: internalErr,
	}
}

func ensureStatusOK(resp *http.Response, successCode int) error {
	httpStatusCodeRange := int(resp.StatusCode / 100)
	if successCode == 0 && httpStatusCodeRange == 2 || resp.StatusCode == successCode {
		return nil
	}

	switch httpStatusCodeRange {
	case 2:
		return createFaultyHTTPResponseError(resp.StatusCode, ErrNotMatched)
	case 4:
		switch resp.StatusCode {
		case http.StatusNotFound:
			return createHTTPResponseError(resp.StatusCode, ErrNotFound)
		case http.StatusForbidden:
			return createHTTPResponseError(resp.StatusCode, ErrForbidden)
		case http.StatusBadRequest:
			return createHTTPResponseError(resp.StatusCode, ErrBadRequest)

		}
	case 5:
		return createFaultyHTTPResponseError(resp.StatusCode, ErrInternalServerError)

	}
	var jerr types.ResponseError
	if err := json.NewDecoder(resp.Body).Decode(&jerr); err != nil {
		return createFaultyHTTPResponseError(resp.StatusCode, err)
	}
	return createHTTPResponseError(resp.StatusCode, errors.New(jerr.Error))
}

func substringReplace(str string, dict map[string]string) string {
	for pattern, value := range dict {
		str = strings.Replace(str, pattern, value, -1)
	}
	return str
}

func CreateKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\",", key, value)
	}
	return b.String()
}

func trimURLSlash(url string) string {
	return strings.TrimPrefix(url, "/")
}

func persistClientConfigFile() error {
	return viper.WriteConfigAs(viper.ConfigFileUsed())
}

func AddItemsAsQueryParams(url string, items map[string]string) string {
	url += "?"
	for key, item := range items {
		url += key + "=" + item + "&"
	}
	return url[:len(url)-1]

}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}
