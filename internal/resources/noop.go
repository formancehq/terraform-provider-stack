package resources

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-stack/internal"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                     = &Noop{}
	_ resource.ResourceWithConfigure        = &Noop{}
	_ resource.ResourceWithConfigValidators = &Noop{}
	_ resource.ResourceWithValidateConfig   = &Noop{}
)

type Noop struct {
	logger logging.Logger
}

type NoopModel struct {
}

func NewNoop(logger logging.Logger) func() resource.Resource {
	return func() resource.Resource {
		return &Noop{
			logger: logger,
		}
	}
}

// Schema implements resource.Resource.
func (s *Noop) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (s *Noop) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, res *resource.ValidateConfigResponse) {
	var config NoopModel
	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if res.Diagnostics.HasError() {
		return
	}

}

// ConfigValidators implements resource.ResourceWithConfigValidators.
func (s *Noop) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return nil
}

// Configure implements resource.ResourceWithConfigure.
func (s *Noop) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	_, ok := req.ProviderData.(internal.Store)
	if !ok {
		res.Diagnostics.AddError(
			"Invalid Provider Data",
			fmt.Sprintf("Expected *formance.Formance, got: %T", req.ProviderData),
		)
		return
	}
}

// Create implements resource.Resource.
func (s *Noop) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	ctx = logging.ContextWithLogger(ctx, s.logger)
	var plan NoopModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (s *Noop) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	ctx = logging.ContextWithLogger(ctx, s.logger)

	var state NoopModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

}

// Metadata implements resource.Resource.
func (s *Noop) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_noop"
}

// Read implements resource.Resource.
func (s *Noop) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state NoopModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (s *Noop) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	ctx = logging.ContextWithLogger(ctx, s.logger)

	var state NoopModel
	var plan NoopModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, &state)...)

}
