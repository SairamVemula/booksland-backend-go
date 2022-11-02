package utils

import "net/http"

type RestError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Error   string `json:"error"`
}

func BadRequest(message string) *RestError {
	return &RestError{
		Message: message,
		Code:    http.StatusBadRequest,
		Error:   "bad request",
	}
}

func NotFound(message string) *RestError {
	return &RestError{
		Message: message,
		Code:    http.StatusNotFound,
		Error:   "not found",
	}
}

func InternalErr(message string) *RestError {
	return &RestError{
		Message: message,
		Code:    http.StatusInternalServerError,
		Error:   "internal server error",
	}
}
