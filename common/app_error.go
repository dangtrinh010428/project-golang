package common

import "net/http"

type AppError struct {
	StatusCode int    `json:"status_code"`
	RootErr    error  `json:"-"`
	Message    string `json:"message"`
	Log        string `json:"log"`
	Key        string `json:"key"`
}

func NewFullErrorResponse(StatusCode int, RootErr error, Message string, Log string, Key string) *AppError {
	return &AppError{
		StatusCode: StatusCode,
		RootErr:    RootErr,
		Message:    Message,
		Log:        Log,
		Key:        Key,
	}
}

func NewErrorResponse(RootErr error, Message string, Log string, Key string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		RootErr:    RootErr,
		Message:    Message,
		Log:        Log,
		Key:        Key,
	}
}

func NewUnauthorized(RootErr error, Message string, Log string, Key string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		RootErr:    RootErr,
		Message:    Message,
		Log:        Log,
		Key:        Key,
	}
}

func (e *AppError) RootError() error {
	if err, ok := e.RootErr.(*AppError); ok {
		return err.RootError()
	}

	return e.RootErr
}

func (e *AppError) Error() string {
	return e.RootError().Error()
}
