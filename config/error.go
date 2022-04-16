package config

import (
	"errors"
	"log"
	"net/http"
)

var (
	ErrBadRequest = errors.New("Bad request")
	ErrInternal = errors.New("Internal error")
	ErrInvalidAPICall = errors.New("Invalid API call")
	ErrNotAuthenticated = errors.New("Not Authenticated")
	ErrResourceNotFound = errors.New("Resource not found")
)

type ErrorResponse struct {
	ErrorCode int
	Cause string
}

const (
	ErrorCodeInternal = 0
	ErrorCodeInvalidJSONBody = 30
	ErrorCodeInvalidCredentials = 201
	ErrorCodeEntityNotFound = 404
	ErrorCodeValidation = 500
)

type serverError struct {
	code int
	cause string
	errorType error
}

func (e serverError) Error() string {
	return e.cause
}

var (
	MapErrorTypeToHTTPStatus = mapErrorTypeToHTTPStatus
	IsError = isError
	NewError = newError
)

func mapErrorTypeToHTTPStatus(err error) int {
	switch err {
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrInternal:
		return http.StatusInternalServerError
	case ErrInvalidAPICall, ErrResourceNotFound:
		return http.StatusNotFound
	case ErrNotAuthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func isError(errorType error) (bool, int, string, error) {
	err, isError := errorType.(serverError)
	if !isError {
		return false, 0, "", errorType
	}
	return true, err.code, err.cause, err.errorType
}

func newError(cause string, code int, errorType, err error) error {
	if err != nil {
		log.Printf("error: %v: %v", cause, err)
	} else {
		log.Printf("error: %v", cause)
	}
	return serverError{code, cause, errorType}
}