package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application-level error with a type and HTTP status code.
type AppError struct {
	Type       string // error type: validation, auth, forbidden, conflict, notfound, etc.
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return e.Message
}

// Error constructors
func ErrValidation(msg string) *AppError {
	return &AppError{Type: "validation", Message: msg, StatusCode: http.StatusBadRequest}
}

func ErrUnauthorized(msg string) *AppError {
	return &AppError{Type: "unauthorized", Message: msg, StatusCode: http.StatusUnauthorized}
}

func ErrForbidden(msg string) *AppError {
	return &AppError{Type: "forbidden", Message: msg, StatusCode: http.StatusForbidden}
}

func ErrNotFound(msg string) *AppError {
	return &AppError{Type: "notfound", Message: msg, StatusCode: http.StatusNotFound}
}

func ErrConflict(msg string) *AppError {
	return &AppError{Type: "conflict", Message: msg, StatusCode: http.StatusConflict}
}

func ErrInternal(msg string) *AppError {
	return &AppError{Type: "internal", Message: msg, StatusCode: http.StatusInternalServerError}
}

// ToHTTP converts an error to HTTP status code and message.
// If err is an AppError, uses its code; otherwise returns 500.
func ToHTTP(err error) (int, map[string]interface{}) {
	if err == nil {
		return http.StatusOK, nil
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr.StatusCode, map[string]interface{}{
			"error": appErr.Type,
			"message": appErr.Message,
		}
	}

	// Default to internal server error
	return http.StatusInternalServerError, map[string]interface{}{
		"error":   "internal",
		"message": fmt.Sprintf("internal server error: %v", err),
	}
}
