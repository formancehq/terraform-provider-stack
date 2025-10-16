package otlp

import (
	"go.opentelemetry.io/otel"
)

var Tracer = otel.Tracer("github.com/formancehq/terraform-provider-stack")
