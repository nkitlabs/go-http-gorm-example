package errors_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	apierror "github.com/nkitlabs/go-http-gorm-example/pkg/errors"
)

func TestWrapError(t *testing.T) {
	tcs := []struct {
		name string
		in   error
		out  string
	}{
		{
			name: "TestNormalApiError",
			in:   apierror.NewError(400, "test error"),
			out:  "test error",
		},
		{
			name: "TestNormalWrapApiError",
			in:   apierror.NewError(400, "test error").Wrap("wrapped"),
			out:  "wrapped: test error",
		},
		{
			name: "TestNormalWrapApiError",
			in:   apierror.NewError(400, "test error").Wrapf("wrapped %d", 1),
			out:  "wrapped 1: test error",
		},
		{
			name: "TestNormalWrapApiError",
			in:   errors.Wrap(apierror.NewError(400, "test error").Wrap("wrapped"), "errors wrap"),
			out:  "errors wrap: wrapped: test error",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			out := tc.in.Error()
			assert.Equal(t, tc.out, out)
		})
	}
}

func TestPanicWrappedErr(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "wrapped error: test panic", r.(error).Error())
		}
	}()

	err := apierror.NewError(400, "test panic")
	newErr := err.Wrap("wrapped error")

	panic(newErr)
}
