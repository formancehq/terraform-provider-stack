package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-stack/internal"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                     = &LedgerSchema{}
	_ resource.ResourceWithConfigure        = &LedgerSchema{}
	_ resource.ResourceWithConfigValidators = &LedgerSchema{}
	_ resource.ResourceWithValidateConfig   = &LedgerSchema{}
)

type LedgerSchema struct {
	store *internal.ModuleStore
}

// ConfigValidators implements resource.ResourceWithConfigValidators.
func (s *LedgerSchema) ConfigValidators(context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("chart"),
			path.MatchRoot("transactions"),
		),
	}
}

type LedgerSchemaModel struct {
	Version        types.String  `tfsdk:"version"`
	Ledger         types.String  `tfsdk:"ledger"`
	Chart          types.Dynamic `tfsdk:"chart"`
	Transactions   types.Dynamic `tfsdk:"transactions"`
	IdempotencyKey types.String  `tfsdk:"idempotency_key"`
}

func NewLedgerSchema() func() resource.Resource {
	return func() resource.Resource {
		return &LedgerSchema{}
	}
}

var SchemaLedgerSchema = schema.Schema{
	Description: "Resource for managing a Formance Ledger Schema. For advanced usage and configuration, see the [Ledger documentation](https://docs.formance.com/ledger/).",
	Attributes: map[string]schema.Attribute{
		"version": schema.StringAttribute{
			Required:    true,
			Description: "The version of the schema.",
		},
		"ledger": schema.StringAttribute{
			Required:    true,
			Description: "The name of the ledger.",
		},
		"chart": schema.DynamicAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The chart of account definition in JSON format.",
			PlanModifiers: []planmodifier.Dynamic{
				dynamicplanmodifier.RequiresReplace(),
			},
			Default: dynamicdefault.StaticValue(types.DynamicValue(types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{}))),
		},
		"idempotency_key": schema.StringAttribute{
			Optional:    true,
			Description: "The idempotency key of the schema.",
		},
		"transactions": schema.DynamicAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The transaction templates defined in the schema.",
			PlanModifiers: []planmodifier.Dynamic{
				dynamicplanmodifier.RequiresReplace(),
			},
			Default: dynamicdefault.StaticValue(types.DynamicValue(types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{}))),
		},
	},
}

// Schema implements resource.Resource.
func (s *LedgerSchema) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = SchemaLedgerSchema
}

// Configure implements resource.ResourceWithConfigure.
func (s *LedgerSchema) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

	s.store = store.NewModuleStore("ledger")
}

func (s *LedgerSchemaModel) parseSchema() (map[string]shared.V2ChartSegment, error) {
	_, ok := s.Chart.UnderlyingValue().(types.Object)
	if !ok {
		return nil, fmt.Errorf("schema is not a valid JSON object")
	}
	v2Schema := map[string]shared.V2ChartSegment{}
	if err := json.Unmarshal([]byte(s.Chart.String()), &v2Schema); err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}

	return v2Schema, nil
}

func (s *LedgerSchemaModel) parseTransaction() (map[string]shared.V2TransactionTemplate, error) {
	_, ok := s.Transactions.UnderlyingValue().(types.Object)
	if !ok {
		return nil, fmt.Errorf("transactions is not a valid JSON object")
	}
	v2Schema := map[string]shared.V2TransactionTemplate{}
	if err := json.Unmarshal([]byte(s.Transactions.String()), &v2Schema); err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}

	return v2Schema, nil
}

func (s *LedgerSchema) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, res *resource.ValidateConfigResponse) {
	var conf LedgerSchemaModel
	res.Diagnostics.Append(req.Config.Get(ctx, &conf)...)
	if res.Diagnostics.HasError() {
		return
	}

	if !conf.Chart.IsNull() {
		if _, ok := conf.Chart.UnderlyingValue().(types.Object); !ok {
			res.Diagnostics.AddError("Invalid Ledger Query", "The ledger_query must be a valid JSON object.")
		} else {
			logging.FromContext(ctx).Debug("Ledger query is valid")
			_, err := conf.parseSchema()
			if err != nil {
				res.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("Failed to create configuration for reconciliation policy: %s", err))
			}
		}
	}

	if !conf.Transactions.IsNull() {
		if _, ok := conf.Transactions.UnderlyingValue().(types.Object); !ok {
			res.Diagnostics.AddError("Invalid Ledger Query", "The ledger_query must be a valid JSON object.")
		} else {
			logging.FromContext(ctx).Debug("Ledger query is valid")
			_, err := conf.parseTransaction()
			if err != nil {
				res.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("Failed to create configuration for reconciliation policy: %s", err))
			}
		}
	}

}

// Create implements resource.Resource.
func (s *LedgerSchema) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan LedgerSchemaModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		chartData       map[string]shared.V2ChartSegment
		transactionData map[string]shared.V2TransactionTemplate
		err             error
	)
	if !plan.Chart.IsNull() {
		chartData, err = plan.parseSchema()
		if err != nil {
			res.Diagnostics.AddError("Invalid Schema", fmt.Sprintf("Failed to parse schema: %s", err))
			return
		}
	}
	if !plan.Transactions.IsNull() {
		transactionData, err = plan.parseTransaction()
		if err != nil {
			res.Diagnostics.AddError("Invalid Transactions", fmt.Sprintf("Failed to parse transactions: %s", err))
			return
		}
	}

	var IK *string
	if !plan.IdempotencyKey.IsNull() {
		ik := plan.IdempotencyKey.ValueString()
		IK = &ik
	}

	config := operations.V2InsertSchemaRequest{
		Ledger:         plan.Ledger.ValueString(),
		Version:        plan.Version.ValueString(),
		IdempotencyKey: IK,
		V2SchemaData: shared.V2SchemaData{
			Chart:        chartData,
			Transactions: transactionData,
		},
	}
	ledgerSdk := s.store.Ledger()
	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}
	_, err = ledgerSdk.InsertSchema(ctx, config)
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (s *LedgerSchema) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	res.Diagnostics.AddWarning("Delete not implemented", "The Delete method for LedgerSchema is not implemented.")

}

// Metadata implements resource.Resource.
func (s *LedgerSchema) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ledger_schema"
}

// Read implements resource.Resource.
func (s *LedgerSchema) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state LedgerSchemaModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	ledgerSdk := s.store.Ledger()
	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}
	readSchemaResponse, err := ledgerSdk.GetSchema(ctx, operations.V2GetSchemaRequest{
		Ledger:  state.Ledger.ValueString(),
		Version: state.Version.ValueString(),
	})
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}

	chart := readSchemaResponse.V2SchemaResponse.Data.Chart
	if len(chart) >= 0 {
		data, err := json.Marshal(chart)
		if err != nil {
			res.Diagnostics.AddError("Schema Marshalling Error", fmt.Sprintf("Failed to marshal schema: %s", err))
			return
		}
		var m = make(map[string]any)
		if err := json.Unmarshal(data, &m); err != nil {
			res.Diagnostics.AddError("Schema Unmarshalling Error", fmt.Sprintf("Failed to unmarshal schema: %s", err))
			return
		}
		tfValues := ConvertToAttrValues(m)
		state.Chart = types.DynamicValue(NewDynamicObjectValue(tfValues).Value())
	}
	transactions := readSchemaResponse.V2SchemaResponse.Data.Transactions
	if len(transactions) >= 0 {
		data, err := json.Marshal(transactions)
		if err != nil {
			res.Diagnostics.AddError("Transactions Marshalling Error", fmt.Sprintf("Failed to marshal transactions: %s", err))
			return
		}
		var m = make(map[string]any)
		if err := json.Unmarshal(data, &m); err != nil {
			res.Diagnostics.AddError("Transactions Unmarshalling Error", fmt.Sprintf("Failed to unmarshal transactions: %s", err))
			return
		}
		tfValues := ConvertToAttrValues(m)
		state.Transactions = types.DynamicValue(NewDynamicObjectValue(tfValues).Value())
	}

	state.Version = types.StringValue(readSchemaResponse.V2SchemaResponse.Data.Version)
	state.Ledger = types.StringValue(state.Ledger.ValueString())
	res.Diagnostics.Append(res.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (s *LedgerSchema) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	res.Diagnostics.AddWarning("Update not implemented", "The Update method for LedgerSchema is not implemented. Recreating the resource. Make sure to inscrease the version")
}
