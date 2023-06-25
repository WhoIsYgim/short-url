package errs

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	Errors = map[error]struct {
		Code    int
		Message string
	}{
		InternalServerError: {
			http.StatusInternalServerError,
			InternalServerError.Error(),
		},
		BadRequest: {
			http.StatusBadRequest,
			BadRequest.Error(),
		},
		UnableToCreateLink: {
			http.StatusConflict,
			UnableToCreateLink.Error(),
		},
		LinkNotFound: {
			http.StatusNotFound,
			LinkNotFound.Error(),
		},
	}
)

var (
	InternalServerError = errors.New("internal server error")
	BadRequest          = errors.New("bad request")
	UnableToCreateLink  = errors.New("unable to create link")

	LinkNotFound = errors.New("link not found")
)

type AppError struct {
	err           error
	internalError error
}

func (ae AppError) Error() string {
	if ae.internalError == nil {
		return fmt.Sprintf("[error]: %s", ae.err)
	}
	return fmt.Sprintf("[error]: %s", ae.internalError)
}

func (ae AppError) Unwrap() error {
	return ae.err
}

func NewAppError(err, internal error) *AppError {
	return &AppError{
		err:           err,
		internalError: internal,
	}
}

func BadRequestError() *AppError {
	return NewAppError(BadRequest, nil)
}

func InternalError(err error) *AppError {
	return NewAppError(InternalServerError, err)
}

func NotFoundError() *AppError {
	return NewAppError(LinkNotFound, nil)
}
