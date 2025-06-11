package resources

import (
	"context"
	"fmt"

	formance "github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-stack/internal"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                     = &Webhooks{}
	_ resource.ResourceWithConfigure        = &Webhooks{}
	_ resource.ResourceWithConfigValidators = &Webhooks{}
	_ resource.ResourceWithValidateConfig   = &Webhooks{}
)

type Webhooks struct {
	logger logging.Logger
	sdk    *formance.Formance
}

type WebhooksModel struct {
	ID         types.String `tfsdk:"id"`
	Endpoint   types.String `tfsdk:"endpoint"`
	EventTypes types.List   `tfsdk:"event_types"`
	Name       types.String `tfsdk:"name"`
	Secret     types.String `tfsdk:"secret"`
}

func NewWebhooks(logger logging.Logger) func() resource.Resource {
	return func() resource.Resource {
		return &Webhooks{
			logger: logger,
		}
	}
}

var SchemaWebhooks = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"endpoint": schema.StringAttribute{
			Required:    true,
			Description: "The endpoint to which webhooks will be sent.",
		},
		"event_types": schema.ListAttribute{
			Required:    true,
			ElementType: types.StringType,
			Description: "The types of events that will trigger webhooks.",
		},
		"name": schema.StringAttribute{
			Optional:    true,
			Description: "A name for the webhook configuration.",
		},
		"secret": schema.StringAttribute{
			Optional:  true,
			Sensitive: true,
		},
	},
}

// Schema implements resource.Resource.
func (s *Webhooks) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = SchemaWebhooks
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (s *Webhooks) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, res *resource.ValidateConfigResponse) {
	var config WebhooksModel
	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if res.Diagnostics.HasError() {
		return
	}

}

// ConfigValidators implements resource.ResourceWithConfigValidators.
func (s *Webhooks) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return nil
}

// Configure implements resource.ResourceWithConfigure.
func (s *Webhooks) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	sdk, ok := req.ProviderData.(*formance.Formance)
	if !ok {
		res.Diagnostics.AddError(
			"Invalid Provider Data",
			fmt.Sprintf("Expected *formance.Formance, got: %T", req.ProviderData),
		)
		return
	}

	s.sdk = sdk

	if err := internal.CheckModuleHealth(ctx, s.sdk, "webhooks"); err != nil {
		res.Diagnostics.AddError(
			"Webhooks Module Not Healthy",
			fmt.Sprintf("The webhooks module is not healthy: %s", err.Error()),
		)
		return
	}

}

// Create implements resource.Resource.
func (s *Webhooks) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan WebhooksModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	config := shared.ConfigUser{}
	if plan.Name.ValueString() != "" {
		config.Name = pointer.For(plan.Name.ValueString())
	}
	if plan.Secret.ValueString() != "" {
		config.Secret = pointer.For(plan.Secret.ValueString())
	}
	config.Endpoint = plan.Endpoint.ValueString()
	config.EventTypes = collectionutils.Map(plan.EventTypes.Elements(), func(v attr.Value) string {
		return v.(types.String).ValueString()
	})

	resp, err := s.sdk.Webhooks.V1.InsertConfig(ctx, config)
	if err != nil {
		res.Diagnostics.AddError(
			"Error Creating Webhook Configuration",
			fmt.Sprintf("Unable to create webhook configuration: %s", err.Error()),
		)
		return
	}
	data := resp.ConfigResponse.Data
	plan.ID = types.StringValue(data.ID)
	plan.Endpoint = types.StringValue(data.Endpoint)
	plan.EventTypes = types.ListValueMust(types.StringType, collectionutils.Map(data.EventTypes, func(s string) attr.Value {
		return types.StringValue(s)
	}))

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...) // Save the plan as state
}

// Delete implements resource.Resource.
func (s *Webhooks) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state WebhooksModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	_, err := s.sdk.Webhooks.V1.DeleteConfig(ctx, operations.DeleteConfigRequest{
		ID: state.ID.ValueString(),
	})
	if err != nil {
		res.Diagnostics.AddError(
			"Error Creating Webhook Configuration",
			fmt.Sprintf("Unable to create webhook configuration: %s", err.Error()),
		)
		return
	}
}

// Metadata implements resource.Resource.
func (s *Webhooks) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhooks"
}

// Read implements resource.Resource.
func (s *Webhooks) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state WebhooksModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	resp, err := s.sdk.Webhooks.V1.GetManyConfigs(ctx, operations.GetManyConfigsRequest{
		ID: pointer.For(state.ID.ValueString()),
	})
	if err != nil {
		res.Diagnostics.AddError(
			"Error Creating Webhook Configuration",
			fmt.Sprintf("Unable to create webhook configuration: %s", err.Error()),
		)
		return
	}

	data := resp.ConfigsResponse.Cursor.Data
	if len(data) == 0 {
		res.Diagnostics.AddError(
			"Webhook Configuration Not Found",
			fmt.Sprintf("Webhook configuration with ID %s not found", state.ID.ValueString()),
		)
		return
	}
	config := data[0]
	state.ID = types.StringValue(config.ID)
	state.Endpoint = types.StringValue(config.Endpoint)
	state.EventTypes = types.ListValueMust(types.StringType, collectionutils.Map(config.EventTypes, func(s string) attr.Value {
		return types.StringValue(s)
	}))
	state.Secret = types.StringValue(config.Secret)

	res.Diagnostics.Append(res.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (s *Webhooks) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var state WebhooksModel
	// var plan WebhooksModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	res.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, &state)...)

}
