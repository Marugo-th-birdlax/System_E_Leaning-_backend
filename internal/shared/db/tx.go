package errors

import "net/http"

type AppError struct {
	Code    string
	Message string
	Status  int
}

func (e *AppError) Error() string { return e.Message }

func New(code, msg string, status int) *AppError {
	return &AppError{Code: code, Message: msg, Status: status}
}

var (
	ErrNotFound   = New("NOT_FOUND", "resource not found", http.StatusNotFound)
	ErrBadRequest = New("BAD_REQUEST", "bad request", http.StatusBadRequest)
	ErrInternal   = New("INTERNAL_ERROR", "internal error", http.StatusInternalServerError)
)
