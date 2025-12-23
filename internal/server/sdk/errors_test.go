package sdk

import (
	"context"
	"errors"
	"testing"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/require"
)

func TestHandleSDKError(t *testing.T) {
	for _, tt := range []struct {
		name     string
		err      error
		expected diag.Diagnostic
	}{
		{
			name: "Test case 1",
			err:  errors.New(""),
			expected: diag.NewErrorDiagnostic(
				"INTERNAL",
				"",
			),
		},
		{
			name: "Error string",
			err:  errors.New(`{"errorCode":"VALIDATION","errorMessage":"invalid config: polling period invalid: polling period cannot be lower than minimum of 20m0s: validation error: validation error"}`),
			expected: diag.NewErrorDiagnostic(
				"VALIDATION",
				"invalid config: polling period invalid: polling period cannot be lower than minimum of 20m0s: validation error: validation error",
			),
		},
		{
			name: "SDKError case",
			err: &sdkerrors.SDKError{
				Body: `{"errorCode":"SOME_ERROR","errorMessage":"An error occurred"}`,
			},
			expected: diag.NewErrorDiagnostic("SOME_ERROR", "An error occurred"),
		},
		{
			name: "V2ErrorResponse case",
			err: &sdkerrors.V2ErrorResponse{
				ErrorCode:    "SOME_ERROR",
				ErrorMessage: "An error occurred",
			},
			expected: diag.NewErrorDiagnostic("SOME_ERROR", "An error occurred"),
		},
		{
			name: "V2BulkResponse case",
			err: &sdkerrors.V2BulkResponse{
				ErrorCode:    pointer.For[shared.V2ErrorsEnum]("SOME_ERROR"),
				ErrorMessage: pointer.For("An error occurred"),
			},
			expected: diag.NewErrorDiagnostic("SOME_ERROR", "An error occurred"),
		},
		{
			name: "v3error response case",
			err: &sdkerrors.ErrorResponse{
				ErrorMessage: "a message",
				ErrorCode:    "123",
			},
			expected: diag.NewErrorDiagnostic("123", "a message"),
		},
		{
			name:     "invalid error type",
			err:      errors.New("some random error"),
			expected: diag.NewErrorDiagnostic("INTERNAL", "some random error"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			diag := make(diag.Diagnostics, 0)
			HandleStackError(context.Background(), tt.err, &diag)
			require.Len(t, diag, 1)
			require.Equal(t, tt.expected, diag[0])
		})
	}
}
