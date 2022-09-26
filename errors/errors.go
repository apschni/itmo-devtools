package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type Type string

const (
	Authorization        Type = "AUTHORIZATION"        // Authentication Failures - 401
	BadRequest           Type = "BADREQUEST"           // Validation errors / BadInput - 400
	Conflict             Type = "CONFLICT"             // Already exists (eg, create account with existent email) - 409
	Internal             Type = "INTERNAL"             // Server and fallback errors - 500
	NotFound             Type = "NOTFOUND"             // For not finding resource - 404
	PayloadTooLarge      Type = "PAYLOADTOOLARGE"      // For uploading tons of JSON, or an image over the limit - 413
	UnsupportedMediaType Type = "UNSUPPORTEDMEDIATYPE" // For http 415
	ServiceUnavailable   Type = "SERVICE_UNAVAILABLE"  // For long running handlers
	Forbidden            Type = "FORBIDDEN"            // For http 403
)

type Error struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Status() int {
	switch e.Type {
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case PayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	case Forbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}

func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: fmt.Sprintf("Bad request. Reason: %v", reason),
	}
}

func NewConflict(name string, value string) *Error {
	return &Error{
		Type:    Conflict,
		Message: fmt.Sprintf("resource: %v with value: %v already exists", name, value),
	}
}

func NewForbidden() *Error {
	return &Error{Type: Forbidden,
		Message: "You don't have permission to access this resource.",
	}
}

func NewInternal() *Error {
	return &Error{
		Type:    Internal,
		Message: "Internal server error.",
	}
}

func NewNotFound(name string, value string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("resource: %v with value: %v not found", name, value),
	}
}

func NewPayloadTooLarge(maxBodySize int64, contentLength int64) *Error {
	return &Error{
		Type:    PayloadTooLarge,
		Message: fmt.Sprintf("Max payload size of %v exceeded. Actual payload size: %v", maxBodySize, contentLength),
	}
}

func NewUnsupportedMediaType(reason string) *Error {
	return &Error{
		Type:    UnsupportedMediaType,
		Message: reason,
	}
}

func NewServiceUnavailable() *Error {
	return &Error{
		Type:    ServiceUnavailable,
		Message: fmt.Sprintf("Service unavailable or timed out"),
	}
}
