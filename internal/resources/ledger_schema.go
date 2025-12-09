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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                     = &Ledger{}
	_ resource.ResourceWithConfigure        = &Ledger{}
	_ resource.ResourceWithConfigValidators = &Ledger{}
	_ resource.ResourceWithValidateConfig   = &Ledger{}
)

type LedgerSchema struct {
	store *internal.ModuleStore
}

type LedgerSchemaModel struct {
	Version types.String  `tfsdk:"version"`
	Ledger  types.String  `tfsdk:"ledger"`
	Schema  types.Dynamic `tfsdk:"schema"`
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
		"schema": schema.DynamicAttribute{
			Required:    true,
			Description: "The schema definition in JSON format.",
			PlanModifiers: []planmodifier.Dynamic{
				dynamicplanmodifier.RequiresReplace(),
			},
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
	_, ok := s.Schema.UnderlyingValue().(types.Object)
	if !ok {
		return nil, fmt.Errorf("schema is not a valid JSON object")
	}
	v2Schema := map[string]shared.V2ChartSegment{}
	if err := json.Unmarshal([]byte(s.Schema.String()), &v2Schema); err != nil {
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

	if _, ok := conf.Schema.UnderlyingValue().(types.Object); !ok {
		res.Diagnostics.AddError("Invalid Ledger Query", "The ledger_query must be a valid JSON object.")
	} else {
		logging.FromContext(ctx).Debug("Ledger query is valid")
		_, err := conf.parseSchema()
		if err != nil {
			res.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("Failed to create configuration for reconciliation policy: %s", err))
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

	schemaData, err := plan.parseSchema()
	if err != nil {
		res.Diagnostics.AddError("Invalid Schema", fmt.Sprintf("Failed to parse schema: %s", err))
		return
	}
	config := operations.V2InsertSchemaRequest{
		Ledger:  plan.Ledger.ValueString(),
		Version: plan.Version.ValueString(),
		V2SchemaData: shared.V2SchemaData{
			Chart: schemaData,
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

	schema := readSchemaResponse.V2SchemaResponse.Data.Chart
	if len(schema) >= 0 {
		data, err := json.Marshal(schema)
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
		state.Schema = types.DynamicValue(NewDynamicObjectValue(tfValues).Value())
	}
	state.Version = types.StringValue(readSchemaResponse.V2SchemaResponse.Data.Version)
	state.Ledger = types.StringValue(state.Ledger.ValueString())
	res.Diagnostics.Append(res.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (s *LedgerSchema) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	res.Diagnostics.AddWarning("Update not implemented", "The Update method for LedgerSchema is not implemented. Recreating the resource. Make sure to inscrease the version")
}
