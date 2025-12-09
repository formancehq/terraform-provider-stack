package resources

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-stack/pkg/tracing"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var (
	_ resource.Resource                     = &ResourceTracer{}
	_ resource.ResourceWithConfigure        = &ResourceTracer{}
	_ resource.ResourceWithImportState      = &ResourceTracer{}
	_ resource.ResourceWithValidateConfig   = &ResourceTracer{}
	_ resource.ResourceWithConfigValidators = &ResourceTracer{}
)
var (
	ErrValidateConfig = fmt.Errorf("error during ValidateConfig")
	ErrSchema         = fmt.Errorf("error during Schema")
	ErrConfigure      = fmt.Errorf("error during Configure")
	ErrCreate         = fmt.Errorf("error during Create")
	ErrRead           = fmt.Errorf("error during Read")
	ErrUpdate         = fmt.Errorf("error during Update")
	ErrDelete         = fmt.Errorf("error during Delete")
	ErrImportState    = fmt.Errorf("error during ImportState")
)

func injectTraceContext(ctx context.Context, res any, funcName string) context.Context {
	name := reflect.TypeOf(res).Elem().Name()
	ctx = logging.ContextWithField(ctx, "resource", strings.ToLower(name))
	ctx = logging.ContextWithField(ctx, "operation", strings.ToLower(funcName))

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ctx
	}

	// TODO: implement a logger hook to automatically add trace context to logs
	headerCarrier := propagation.MapCarrier{}
	propagation.TraceContext{}.Inject(ctx, headerCarrier)
	for k, v := range headerCarrier {
		ctx = logging.ContextWithField(ctx, k, v)
	}

	span.SetAttributes(
		attribute.String("resource", strings.ToLower(name)),
		attribute.String("operation", strings.ToLower(funcName)),
	)
	return ctx
}

type ResourceTracer struct {
	tracer          trace.Tracer
	logger          logging.Logger
	underlyingValue any
}

func NewResourceTracer(tracer trace.Tracer, logger logging.Logger, res any) func() resource.Resource {
	return func() resource.Resource {
		return &ResourceTracer{
			tracer:          tracer,
			logger:          logger,
			underlyingValue: res,
		}
	}

}

// ConfigValidators implements resource.ResourceWithConfigValidators.
func (r *ResourceTracer) ConfigValidators(context.Context) []resource.ConfigValidator {
	operation := "ConfigValidators"
	ctx := logging.ContextWithLogger(context.Background(), r.logger)
	var validators []resource.ConfigValidator
	if v, ok := r.underlyingValue.(resource.ResourceWithConfigValidators); ok {
		_ = tracing.TraceError(ctx, r.tracer, operation, func(ctx context.Context) error {
			ctx = injectTraceContext(ctx, v, operation)
			logging.FromContext(ctx).Debug("call")
			defer logging.FromContext(ctx).Debug("completed")
			validators = v.ConfigValidators(ctx)
			return nil
		})
	}
	return validators
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (r *ResourceTracer) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	operation := "ValidateConfig"
	ctx = logging.ContextWithLogger(ctx, r.logger)
	if v, ok := r.underlyingValue.(resource.ResourceWithValidateConfig); ok {
		_ = tracing.TraceError(ctx, r.tracer, operation, func(ctx context.Context) error {
			ctx = injectTraceContext(ctx, v, operation)
			logging.FromContext(ctx).Debug("call")
			defer logging.FromContext(ctx).Debug("completed")
			v.ValidateConfig(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				return ErrValidateConfig
			}
			return nil
		})
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *ResourceTracer) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	operation := "ImportState"
	ctx = logging.ContextWithLogger(ctx, r.logger)
	if v, ok := r.underlyingValue.(resource.ResourceWithImportState); ok {
		_ = tracing.TraceError(ctx, r.tracer, operation, func(ctx context.Context) error {
			ctx = injectTraceContext(ctx, v, operation)
			logging.FromContext(ctx).Debug("call")
			defer logging.FromContext(ctx).Debug("completed")
			v.ImportState(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				return ErrImportState
			}
			return nil
		})
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *ResourceTracer) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx = logging.ContextWithLogger(ctx, r.logger)
	operation := "Configure"
	if v, ok := r.underlyingValue.(resource.ResourceWithConfigure); ok {
		_ = tracing.TraceError(ctx, r.tracer, operation, func(ctx context.Context) error {
			ctx = injectTraceContext(ctx, v, operation)
			logging.FromContext(ctx).Debug("call")
			defer logging.FromContext(ctx).Debug("completed")
			v.Configure(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				return ErrConfigure
			}
			return nil
		})
	}
}

// Create implements resource.Resource.
func (r *ResourceTracer) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx = logging.ContextWithLogger(ctx, r.logger)
	operation := "Create"
	if v, ok := r.underlyingValue.(resource.Resource); ok {
		_ = tracing.TraceError(ctx, r.tracer, operation, func(ctx context.Context) error {
			ctx = injectTraceContext(ctx, v, operation)
			logging.FromContext(ctx).Debug("call")
			defer logging.FromContext(ctx).Debug("completed")
			v.Create(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				return ErrCreate
			}
			return nil
		})
	}
}

// Delete implements resource.Resource.
func (r *ResourceTracer) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx = logging.ContextWithLogger(ctx, r.logger)
	operation := "Delete"
	if v, ok := r.underlyingValue.(resource.Resource); ok {
		_ = tracing.TraceError(ctx, r.tracer, operation, func(ctx context.Context) error {
			ctx = injectTraceContext(ctx, v, operation)
			logging.FromContext(ctx).Debug("call")
			defer logging.FromContext(ctx).Debug("completed")
			v.Delete(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				return ErrDelete
			}
			return nil
		})
	}
}

// Metadata implements resource.Resource.
func (r *ResourceTracer) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	ctx = logging.ContextWithLogger(ctx, r.logger)
	operation := "Metadata"
	if v, ok := r.underlyingValue.(resource.Resource); ok {
		_ = tracing.TraceError(ctx, r.tracer, operation, func(ctx context.Context) error {
			ctx = injectTraceContext(ctx, v, operation)
			logging.FromContext(ctx).Debug("call")
			defer logging.FromContext(ctx).Debug("completed")
			v.Metadata(ctx, req, resp)
			return nil
		})
	}
}

// Read implements resource.Resource.
func (r *ResourceTracer) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx = logging.ContextWithLogger(ctx, r.logger)
	operation := "Read"
	if v, ok := r.underlyingValue.(resource.Resource); ok {
		_ = tracing.TraceError(ctx, r.tracer, operation, func(ctx context.Context) error {
			ctx = injectTraceContext(ctx, v, operation)
			logging.FromContext(ctx).Debug("call")
			defer logging.FromContext(ctx).Debug("completed")
			v.Read(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				return ErrRead
			}
			return nil
		})
	}
}

// Schema implements resource.Resource.
func (r *ResourceTracer) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx = logging.ContextWithLogger(ctx, r.logger)
	operation := "Schema"
	if v, ok := r.underlyingValue.(resource.Resource); ok {
		_ = tracing.TraceError(ctx, r.tracer, operation, func(ctx context.Context) error {
			ctx = injectTraceContext(ctx, v, operation)
			v.Schema(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				return ErrSchema
			}
			return nil
		})
	}
}

// Update implements resource.Resource.
func (r *ResourceTracer) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx = logging.ContextWithLogger(ctx, r.logger)
	operation := "Update"
	if v, ok := r.underlyingValue.(resource.Resource); ok {
		_ = tracing.TraceError(ctx, r.tracer, operation, func(ctx context.Context) error {
			ctx = injectTraceContext(ctx, v, operation)
			logging.FromContext(ctx).Debug("call")
			defer logging.FromContext(ctx).Debug("completed")
			v.Update(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				return ErrUpdate
			}
			return nil
		})
	}
}
