package sdk

import (
	"context"
	"fmt"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/sdkerrors"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"go.opentelemetry.io/otel/trace"
)

func HandleStackError(ctx context.Context, err error, diag *diag.Diagnostics) {
	sharedError := &sdkerrors.V2ErrorResponse{
		ErrorCode:    "INTERNAL",
		ErrorMessage: "unexpected error",
	}
	switch e := err.(type) {
	case *sdkerrors.V2ErrorResponse:
		sharedError = e
	}
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		diag.AddError("traceparent", fmt.Sprintf("%s-%s", span.SpanContext().TraceID(), span.SpanContext().SpanID()))
	}
	diag.AddError(
		string(sharedError.ErrorCode),
		sharedError.ErrorMessage,
	)
}
