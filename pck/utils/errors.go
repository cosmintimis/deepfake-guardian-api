package utils

import "net/http"

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *CustomError) Error() string {
	return e.Message
}

var ErrMissingID = &CustomError{
	Code:    http.StatusBadRequest,
	Message: "missing id parameter",
}

var ErrMediaNotFound = &CustomError{
	Code:    http.StatusNotFound,
	Message: "media not found",
}
