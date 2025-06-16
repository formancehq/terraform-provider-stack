package sdk

import (
	"errors"
	"net/http"
	"strings"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/sdkerrors"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func HandleStackError(err error, resp *http.Response, diag *diag.Diagnostics) {
	var details []string

	traceparent := resp.Header.Get("traceparent")
	if traceparent != "" {
		details = append(details, "Traceparent: "+traceparent)
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
	case *sdkerrors.V3ErrorResponse:
		details = append(details, "Error Code: "+string(e.ErrorCode))
		if e.Details != nil {
			details = append(details, "Details: "+*e.Details)
		}

		if e.ErrorMessage != "" {
			details = append(details, "Message: "+e.ErrorMessage)
		}
	case *sdkerrors.V2Error:
		details = append(details, "Error Code: "+string(e.ErrorCode))
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
