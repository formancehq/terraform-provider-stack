package sdk

import (
	"context"
	"errors"
	"strings"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/sdkerrors"
	"github.com/formancehq/terraform-provider-stack/pkg"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func HandleStackError(ctx context.Context, err error, diag *diag.Diagnostics) {
	var details []string

	resp := pkg.ResponseFromContext(ctx)
	if resp != nil {
		traceparent := resp.Header.Get("Traceparent")
		if traceparent != "" {
			details = append(details, "Traceparent: "+traceparent)
		}
	}

	switch e := err.(type) {
	case *sdkerrors.V2ErrorResponse:
		details = append(details, "Error Code: "+string(e.ErrorCode))
		if e.Details != nil {
			details = append(details, "Details: "+*e.Details)
		}
		if e.ErrorMessage != "" {
			details = append(details, "Message: "+e.ErrorMessage)
		}
	default:
		details = append(details, errors.ErrUnsupported.Error(), err.Error())
	}

	if len(details) > 1 {
		diag.AddError(details[0], strings.Join(details[1:], "\r\n"))
	} else if len(details) == 1 {
		diag.AddError(details[0], "")
	} else {
		diag.AddError("Unexpected error", "An unexpected error occurred")
	}
}
