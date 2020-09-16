package httperror

import "errors"

type HTTPError struct {
	Message interface{} `json:"message"`
}

var ErrNotFound = errors.New("record not found")
