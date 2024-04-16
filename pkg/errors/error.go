package errors

import (
	"net/http"
)

var (
	ErrInternal     = NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	ErrInvalidInput = NewError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	ErrNotFound     = NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
)

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"-"`
}

func NewError(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

func NewNotFoundError(message string) *Error {
	return ErrNotFound.WithMessage(message)
}

func (e *Error) WithMessage(message string) *Error {
	e.Message = message
	return e
}

// ToError converts an error to an Error object.
func ToError(err error) *Error {
	if err == nil {
		return nil
	}

	for {
		if e, ok := err.(*Error); ok {
			return e
		}

		if e, ok := err.(Error); ok {
			return &e
		}

		if c, ok := err.(causer); ok {
			err = c.Cause()
		} else {
			return ErrInternal
		}
	}
}

func (e Error) Error() string {
	return e.Message
}

// Is check if given error instance is of a given kind/type. This involves
// unwrapping given error using the Cause method if available.
func (e *Error) Is(err error) bool {
	// Reflect usage is necessary to correctly compare with
	// a nil implementation of an error.
	if e == nil {
		return isNilErr(err)
	}

	for {
		if err == e {
			return true
		}

		// If this is a collection of errors, this function must return
		// true if at least one is from the group match.
		if u, ok := err.(unpacker); ok {
			for _, er := range u.Unpack() {
				if e.Is(er) {
					return true
				}
			}
		}

		if c, ok := err.(causer); ok {
			err = c.Cause()
		} else {
			return false
		}
	}
}

// Wrap extends this error with an additional information. It's a handy function to call
// Wrap with this errors package.
func (e Error) Wrap(msg string) error { return Wrap(e, msg) }

// Wrapf extends this error with an additional information. It's a handy function to call
// Wrapf with this errors package.
func (e Error) Wrapf(desc string, args ...interface{}) error { return Wrapf(e, desc, args...) }
