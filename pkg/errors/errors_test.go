package errors_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	apierror "github.com/nkitlabs/go-http-gorm-example/pkg/errors"
)

func TestToError(t *testing.T) {
	ref_err := apierror.NewError(400, "test error")
	ref_err_wrapper := apierror.NewError(400, "test error").Wrap("test error wrapper")
	ref_err_wrapper2 := errors.Wrap(apierror.NewError(400, "test error").Wrap("test error wrapper"), "test")
	expected_err := apierror.NewError(400, "test error")

	tcs := []struct {
		name string
		in   error
		out  *apierror.Error
	}{
		{
			name: "testNormalErrorPointer",
			in:   ref_err,
			out:  expected_err,
		},
		{
			name: "testNormalErrorObject",
			in:   ref_err,
			out:  expected_err,
		},
		{
			name: "testWrappedError",
			in:   ref_err_wrapper,
			out:  expected_err,
		},
		{
			name: "testWrappedError2",
			in:   ref_err_wrapper2,
			out:  expected_err,
		},
		{
			name: "testNoApiError",
			in:   errors.New("test"),
			out:  apierror.ErrInternal,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			errResp := apierror.ToError(tc.in)
			assert.Equal(t, *tc.out, *errResp)
		})
	}
}
