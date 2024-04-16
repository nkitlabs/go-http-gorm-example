package errors

import (
	"fmt"

	goerr "github.com/pkg/errors"
)

type WrappedError struct {
	// This error layer description.
	Msg string `json:"message"`
	// The underlying error that triggered this one.
	parent error `json:"-"`
}

func (e WrappedError) Error() string {
	return fmt.Sprintf("%s: %s", e.Msg, e.parent.Error())
}

func (e WrappedError) Cause() error {
	return e.parent
}

// Is reports whether any error in e's chain matches a target.
func (e *WrappedError) Is(target error) bool {
	if e == target {
		return true
	}

	w := e.Cause()
	for {
		if w == target {
			return true
		}

		x, ok := w.(causer)
		if ok {
			w = x.Cause()
		}
		if x == nil {
			return false
		}
	}
}

// Unwrap implements the built-in errors.Unwrap
func (e *WrappedError) Unwrap() error {
	return e.parent
}

// Wrap extends given error with an additional information. If err is nil, this returns nil,
// avoiding the need for an if statement when wrapping a error returned at the end of a function.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	// If this error does not carry the stacktrace information yet, attach
	// one. This should be done only once per error at the lowest frame
	// possible (most inner wrap).
	if stackTrace(err) == nil {
		err = goerr.WithStack(err)
	}

	return &WrappedError{
		parent: err,
		Msg:    msg,
	}
}

// Wrapf extends given error with an additional information. This function works like Wrap
// function with additional functionality of formatting the input as specified.
func Wrapf(err error, format string, args ...interface{}) error {
	desc := fmt.Sprintf(format, args...)
	return Wrap(err, desc)
}
