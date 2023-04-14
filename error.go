package gofantasy

import (
	"fmt"
)

// Known errors
var (
	ErrNotImplemented = fmt.Errorf("method not implemented")
	ErrBadRequest     = fmt.Errorf("the API doesn’t understand the request. Something is missing")
	ErrUnauthorized   = fmt.Errorf("the API key is missing or misspelled")
	ErrForbidden      = fmt.Errorf("the API key doesn’t have the roles required to perform the request")
	ErrNotFound       = fmt.Errorf("the API understands the request but a parameter is missing or misspelled")
	ErrInternalServer = fmt.Errorf("something went wrong on the server’s side")
	ErrUnknown        = fmt.Errorf("an unknown error was returned")
)

type HTTPError struct {
	StatusCode     int    `json:"statusCode"`
	Status         string `json:"status"`
	Message        string `json:"message"`
	ReportErrorURL string `json:"reportErrorUrl"`
}

type RequestError struct {
	StatusCode int
	Err        error
}

func (e *RequestError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("status code %d", e.StatusCode)
}

type ErrorResponse struct {
	Error *HTTPError `json:"error,omitempty"`
}

func (m *HTTPError) Error() string {
	return m.Message
}

func (m *HTTPError) Cause() error {
	switch m.StatusCode {
	case 400:
		return ErrBadRequest
	case 401:
		return ErrUnauthorized
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFound
	case 500:
		return ErrInternalServer
	}
	return ErrUnknown
}
