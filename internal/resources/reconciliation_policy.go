package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/query"
	"github.com/formancehq/terraform-provider-stack/internal"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                   = &ReconciliationPolicy{}
	_ resource.ResourceWithConfigure      = &ReconciliationPolicy{}
	_ resource.ResourceWithValidateConfig = &ReconciliationPolicy{}
)

type ReconciliationPolicy struct {
	logger logging.Logger
	store  *internal.ModuleStore
}

type ReconciliationPolicyModel struct {
	ID             types.String  `tfsdk:"id"`
	LedgerName     types.String  `tfsdk:"ledger_name"`
	Name           types.String  `tfsdk:"name"`
	PaymentsPoolID types.String  `tfsdk:"payments_pool_id"`
	LedgerQuery    types.Dynamic `tfsdk:"ledger_query"`
	CreatedAt      types.String  `tfsdk:"created_at"`
}

var (
	ErrParseLedgerQuery = fmt.Errorf("failed to parse ledger query")
)

func (m *ReconciliationPolicyModel) CreateConfig() (shared.PolicyRequest, error) {
	var ledgerQuery map[string]any

	if object, ok := m.LedgerQuery.UnderlyingValue().(types.Object); ok {
		fmt.Println(object.String())
		qb, err := query.ParseJSON(object.String())
		if err != nil {
			return shared.PolicyRequest{}, fmt.Errorf("%w: %w", ErrParseLedgerQuery, err)
		}

		if qb != nil {
			if err := json.Unmarshal([]byte(m.LedgerQuery.String()), &ledgerQuery); err != nil {
				return shared.PolicyRequest{}, fmt.Errorf("failed to unmarshal ledger query: %w", err)
			}
		}
	}

	return shared.PolicyRequest{
		LedgerName:     m.LedgerName.ValueString(),
		Name:           m.Name.ValueString(),
		PaymentsPoolID: m.PaymentsPoolID.ValueString(),
		LedgerQuery:    ledgerQuery,
	}, nil
}

func NewReconciliationPolicy(logger logging.Logger) func() resource.Resource {
	return func() resource.Resource {
		return &ReconciliationPolicy{
			logger: logger,
		}
	}
}

var SchemaReconciliationPolicy = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "The unique identifier of the reconciliation policy.",
		},
		"ledger_name": schema.StringAttribute{
			Required: true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the pool.",
			Required:    true,
		},
		"payments_pool_id": schema.StringAttribute{
			Description: "The ID of the payments pool associated with the reconciliation policy.",
			Required:    true,
		},
		"ledger_query": schema.DynamicAttribute{
			Description: "The ledger query used to filter transactions for reconciliation. It must be a valid JSON object representing a query Builder.",
			Optional:    true,
		},
		"created_at": schema.StringAttribute{
			Computed:    true,
			Description: "The timestamp when the reconciliation policy was created.",
		},
	},
}

// Schema implements resource.Resource.
func (s *ReconciliationPolicy) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = SchemaReconciliationPolicy
}

// Configure implements resource.ResourceWithConfigure.
func (s *ReconciliationPolicy) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

	s.store = store.NewModuleStore("reconciliation")

}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (s *ReconciliationPolicy) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, res *resource.ValidateConfigResponse) {
	logger := s.logger.WithField("func", "ValidateConfig")
	logger.Debug("Validating reconciliation policy configuration")
	defer logger.Debug("Finished validating reconciliation policy configuration")
	ctx = logging.ContextWithLogger(ctx, logger)

	var conf ReconciliationPolicyModel
	res.Diagnostics.Append(req.Config.Get(ctx, &conf)...)
	if res.Diagnostics.HasError() {
		return
	}

	if _, ok := conf.LedgerQuery.UnderlyingValue().(types.Object); !ok {
		res.Diagnostics.AddError("Invalid Ledger Query", "The ledger_query must be a valid JSON object.")
	} else {
		logger.Debug("Ledger query is valid")
		_, err := conf.CreateConfig()
		if err != nil {
			res.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("Failed to create configuration for reconciliation policy: %s", err))
		}
	}

}

// Create implements resource.Resource.
func (s *ReconciliationPolicy) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	logger := s.logger.WithField("func", "Create")
	logger.Debug("Creating reconciliation policy")
	defer logger.Debug("Finished creating reconciliation policy")
	ctx = logging.ContextWithLogger(ctx, logger)

	var plan ReconciliationPolicyModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	config, err := plan.CreateConfig()
	if err != nil {
		res.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("Failed to create configuration for reconciliation policy: %s", err))
		return
	}

	sdkReconciliation := s.store.Reconciliation()
	resp, err := sdkReconciliation.CreatePolicy(ctx, config)
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}

	plan.ID = types.StringValue(resp.PolicyResponse.Data.ID)
	plan.CreatedAt = types.StringValue(resp.PolicyResponse.Data.CreatedAt.String())
	fmt.Println(plan)
	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (s *ReconciliationPolicy) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	logger := s.logger.WithField("func", "Delete")
	logger.Debug("Deleting reconciliation policy")
	defer logger.Debug("Finished deleting reconciliation policy")
	ctx = logging.ContextWithLogger(ctx, logger)

	var state ReconciliationPolicyModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	sdkReconciliation := s.store.Reconciliation()
	_, err := sdkReconciliation.DeletePolicy(ctx, operations.DeletePolicyRequest{
		PolicyID: state.ID.ValueString(),
	})
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}
}

// Metadata implements resource.Resource.
func (s *ReconciliationPolicy) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reconciliation_policy"
}

// Read implements resource.Resource.
func (s *ReconciliationPolicy) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	logger := s.logger.WithField("func", "Read")
	logger.Debug("Reading reconciliation policy")
	defer logger.Debug("Finished reading reconciliation policy")
	ctx = logging.ContextWithLogger(ctx, logger)

	var state ReconciliationPolicyModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	sdkReconciliation := s.store.Reconciliation()
	resp, err := sdkReconciliation.GetPolicy(ctx, operations.GetPolicyRequest{
		PolicyID: state.ID.ValueString(),
	})
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}

	state.ID = types.StringValue(resp.PolicyResponse.Data.ID)
	state.LedgerName = types.StringValue(resp.PolicyResponse.Data.LedgerName)
	state.Name = types.StringValue(resp.PolicyResponse.Data.Name)
	state.PaymentsPoolID = types.StringValue(resp.PolicyResponse.Data.PaymentsPoolID)
	state.CreatedAt = types.StringValue(resp.PolicyResponse.Data.CreatedAt.String())

	ledgerQuery := resp.PolicyResponse.Data.LedgerQuery
	if len(ledgerQuery) > 0 {
		tfValues := ConvertToAttrValues(ledgerQuery)
		state.LedgerQuery = types.DynamicValue(NewDynamicObjectValue(tfValues).Value())
	}

	res.Diagnostics.Append(res.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (s *ReconciliationPolicy) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	res.Diagnostics.AddWarning("Update Not Implemented", "The Update method for ReconciliationPolicy is not implemented. Please use Create or Delete to manage reconciliation policies.")
}
