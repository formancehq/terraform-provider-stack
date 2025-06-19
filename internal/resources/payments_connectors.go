package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"maps"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-stack/internal"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExtractKeys(m types.Map) []string {
	i := maps.Keys(m.Elements())
	next, _ := iter.Pull(i)
	keys := []string{}
	for k, ok := next(); ok; k, ok = next() {
		keys = append(keys, k)
	}
	return keys
}

func SanitizeUnknownKeys(m map[string]attr.Value, allowedKeys []string) map[string]attr.Value {
	sanitized := make(map[string]attr.Value)
	for _, key := range allowedKeys {
		if value, exists := m[key]; exists {
			sanitized[key] = value
		}
	}
	return sanitized

}

var (
	_ resource.Resource              = &PaymentsConnectors{}
	_ resource.ResourceWithConfigure = &PaymentsConnectors{}
)

type PaymentsConnectors struct {
	logger logging.Logger
	store  *internal.ModuleStore
}

type PaymentsConnectorsModel struct {
	ID          types.String `tfsdk:"id"`
	Credentials types.Map    `tfsdk:"credentials"`
	Config      types.Map    `tfsdk:"config"`
}

func (m PaymentsConnectorsModel) CreateConfig() (operations.V3InstallConnectorRequest, error) {
	var snakeAS map[string]interface{}
	if err := json.Unmarshal([]byte(m.Config.String()), &snakeAS); err != nil {
		return operations.V3InstallConnectorRequest{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := json.Unmarshal([]byte(m.Credentials.String()), &snakeAS); err != nil {
		return operations.V3InstallConnectorRequest{}, fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	data, err := json.Marshal(snakeAS)
	if err != nil {
		return operations.V3InstallConnectorRequest{}, fmt.Errorf("failed to marshal connector config: %w", err)
	}

	config := operations.V3InstallConnectorRequest{
		V3InstallConnectorRequest: &shared.V3InstallConnectorRequest{},
	}

	if err := config.V3InstallConnectorRequest.UnmarshalJSON(data); err != nil {
		return config, err
	}

	return config, nil
}

func (m PaymentsConnectorsModel) StateFromRequest(resp *shared.V3InstallConnectorRequest) (PaymentsConnectorsModel, error) {
	var plan PaymentsConnectorsModel
	data, err := json.Marshal(resp)
	if err != nil {
		return plan, fmt.Errorf("failed to marshal connector config response: %w", err)
	}

	values := make(map[string]any)
	if err := json.Unmarshal(data, &values); err != nil {
		return plan, fmt.Errorf("failed to unmarshal connector config: %w", err)
	}

	tfValues := ConvertToAttrValues(values)
	attributes := types.MapValueMust(types.DynamicType, tfValues)

	allowedCredsKeys := ExtractKeys(m.Credentials)
	allowedConfigKeys := ExtractKeys(m.Config)

	creds := SanitizeUnknownKeys(attributes.Elements(), allowedCredsKeys)
	config := SanitizeUnknownKeys(attributes.Elements(), allowedConfigKeys)
	plan.Config = types.MapValueMust(types.DynamicType, config)
	plan.Credentials = types.MapValueMust(types.DynamicType, creds)
	plan.ID = m.ID

	return plan, nil
}

func NewPaymentsConnectors(logger logging.Logger) func() resource.Resource {
	return func() resource.Resource {
		return &PaymentsConnectors{
			logger: logger,
		}
	}
}

var SchemaPaymentsConnectors = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"credentials": schema.MapAttribute{
			Sensitive:   true,
			ElementType: types.DynamicType,
			Description: "The credentials for the payment connector. This should include sensitive information like API keys, secrets, certificate, and must be handled securely.",
			Required:    true,
		},
		"config": schema.MapAttribute{
			Required:    true,
			ElementType: types.DynamicType,
			Description: "The configuration for the payment connector. It must not contain sensitive information like API keys or secrets.",
		},
	},
}

// Schema implements resource.Resource.
func (s *PaymentsConnectors) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = SchemaPaymentsConnectors
}

// Configure implements resource.ResourceWithConfigure.
func (s *PaymentsConnectors) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	store, ok := req.ProviderData.(internal.Store)
	if !ok {
		res.Diagnostics.AddError(
			"Invalid Provider Data",
			fmt.Sprintf("Expected *formance.Formance, got: %T", req.ProviderData),
		)
		return
	}

	s.store = store.NewModuleStore("payments")

}

// Create implements resource.Resource.
func (s *PaymentsConnectors) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	ctx = logging.ContextWithLogger(ctx, s.logger)

	var plan PaymentsConnectorsModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	config, err := plan.CreateConfig()
	if err != nil {
		res.Diagnostics.AddError(
			"Invalid Connector Configuration",
			fmt.Sprintf("Failed to create connector configuration: %v", err),
		)
		return
	}

	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	sdkPayments := s.store.Payments()
	resp, err := sdkPayments.CreateConnector(ctx, config)
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}

	plan.ID = types.StringValue(resp.V3InstallConnectorResponse.Data)

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (s *PaymentsConnectors) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	ctx = logging.ContextWithLogger(ctx, s.logger)

	var state PaymentsConnectorsModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	sdkPaymentsConnectors := s.store.Payments()
	_, err := sdkPaymentsConnectors.DeleteConnector(ctx, operations.V3UninstallConnectorRequest{
		ConnectorID: state.ID.ValueString(),
	})
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}
}

// Metadata implements resource.Resource.
func (s *PaymentsConnectors) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_payments_connectors"
}

// Read implements resource.Resource.
func (s *PaymentsConnectors) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	ctx = logging.ContextWithLogger(ctx, s.logger)

	var state PaymentsConnectorsModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	sdkPaymentsConnectors := s.store.Payments()
	resp, err := sdkPaymentsConnectors.GetConnector(ctx, operations.V3GetConnectorConfigRequest{
		ConnectorID: state.ID.ValueString(),
	})
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}

	newPlan, err := state.StateFromRequest(&resp.V3GetConnectorConfigResponse.Data)
	if err != nil {
		res.Diagnostics.AddError(
			"Invalid Connector Configuration",
			fmt.Sprintf("Failed to create connector configuration from response: %v", err),
		)
		return
	}
	res.Diagnostics.Append(res.State.Set(ctx, &newPlan)...)
}

// Update implements resource.Resource.
func (s *PaymentsConnectors) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	ctx = logging.ContextWithLogger(ctx, s.logger)

	var plan PaymentsConnectorsModel
	var state PaymentsConnectorsModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	sdkPayments := s.store.Payments()
	config, err := plan.CreateConfig()
	if err != nil {
		res.Diagnostics.AddError(
			"Invalid Connector Configuration",
			fmt.Sprintf("Failed to create connector configuration: %v", err),
		)
		return
	}

	_, err = sdkPayments.UpdateConnector(ctx, operations.V3UpdateConnectorConfigRequest{
		ConnectorID:               plan.ID.ValueString(),
		V3InstallConnectorRequest: config.V3InstallConnectorRequest,
	})
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}
