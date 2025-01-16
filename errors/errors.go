package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	BAD_REQUEST_ERROR = []ErrorContent{
		{
			Status: http.StatusBadRequest,
			Title:  "Bad Request",
		},
	}

	NOT_FOUND_ERROR = []ErrorContent{
		{
			Status: http.StatusNotFound,
			Title:  "Not Found",
		},
	}

	METHOD_NOT_ALLOWED_ERROR = []ErrorContent{
		{
			Status: http.StatusMethodNotAllowed,
			Title:  "Method Not Allowed",
		},
	}

	UNSUPPORTED_MEDIA_TYPE_ERROR = []ErrorContent{
		{
			Status: http.StatusUnsupportedMediaType,
			Title:  "Unsupported Media Type",
		},
	}

	INTERNAL_SERVER_ERROR = []ErrorContent{
		{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
		},
	}
)

type ErrorContent struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}

type HTTPError struct {
	Status int          `json:"-"`
	Header *http.Header `json:"-"`

	Errors *[]ErrorContent `json:"errors"`
}

func (e *HTTPError) Error() string {
	if e.Errors == nil {
		e.Errors = &INTERNAL_SERVER_ERROR
	}

	return fmt.Sprintf("HTTP %d: %s", e.Status, (*e.Errors)[0].Title)
}

func (e *HTTPError) SetHeader(key, value string) {
	e.Header.Set(key, value)
}

func (e *HTTPError) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	for key, values := range *e.Header {
		for _, value := range values {
			w.Header().Set(key, value)
		}
	}

	w.WriteHeader(e.Status)

	if e.Errors == nil {
		e.Errors = &INTERNAL_SERVER_ERROR
	}

	json.NewEncoder(w).Encode(e)
}

func NewHTTPError(status int, errors *[]ErrorContent) *HTTPError {
	return &HTTPError{
		Status: status,
		Header: &http.Header{},
		Errors: errors,
	}
}
