package resources

import (
	"context"
	"fmt"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-stack/internal"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &ReconciliationPolicy{}
	_ resource.ResourceWithConfigure = &ReconciliationPolicy{}
)

type ReconciliationPolicy struct {
	logger logging.Logger
	store  *internal.ModuleStore
}

type ReconciliationPolicyModel struct {
	ID             types.String `tfsdk:"id"`
	LedgerName     types.String `tfsdk:"ledger_name"`
	Name           types.String `tfsdk:"name"`
	PaymentsPoolID types.String `tfsdk:"payments_pool_id"`
	LedgerQuery    types.Map    `tfsdk:"ledger_query"`
	CreatedAt      types.String `tfsdk:"created_at"`
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
		"ledger_query": schema.MapAttribute{
			Description: "The ledger query used to filter transactions for the reconciliation policy.",
			Optional:    true,
			ElementType: types.DynamicType,
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

	s.store = store.NewModuleStore("payments")

}

// Create implements resource.Resource.
func (s *ReconciliationPolicy) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	ctx = logging.ContextWithLogger(ctx, s.logger)

	var plan ReconciliationPolicyModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	ledgerQuery := map[string]any{}
	// tfValues := ConvertToAttrValues(values)
	// attributes := types.MapValueMust(types.DynamicType, tfValues)
	// plan.LedgerQuery.ElementsAs(ctx, ledgerQuery, false)
	sdkPayments := s.store.Reconciliation()
	resp, err := sdkPayments.CreatePolicy(ctx, shared.PolicyRequest{
		LedgerName:     plan.LedgerName.ValueString(),
		Name:           plan.Name.ValueString(),
		PaymentsPoolID: plan.PaymentsPoolID.ValueString(),
		LedgerQuery:    ledgerQuery,
	})
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}

	plan.CreatedAt = types.StringValue(resp.PolicyResponse.Data.CreatedAt.String())

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (s *ReconciliationPolicy) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	ctx = logging.ContextWithLogger(ctx, s.logger)

	var state ReconciliationPolicyModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	sdkReconciliationPolicy := s.store.Reconciliation()
	_, err := sdkReconciliationPolicy.DeletePolicy(ctx, operations.DeletePolicyRequest{
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
	ctx = logging.ContextWithLogger(ctx, s.logger)

	var state ReconciliationPolicyModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	sdkReconciliationPolicy := s.store.Reconciliation()
	resp, err := sdkReconciliationPolicy.GetPolicy(ctx, operations.GetPolicyRequest{
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
	// state.LedgerQuery = types.MapValueMust(
	// 	types.StringType,
	// 	collectionutils.ConvertMap(resp.PolicyResponse.Data.LedgerQuery, func(value any) attr.Value {
	// 		return types.StringValue(value)
	// 	}),
	// )

	res.Diagnostics.Append(res.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (s *ReconciliationPolicy) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	res.Diagnostics.AddWarning("Update Not Implemented", "The Update method for ReconciliationPolicy is not implemented. Please use Create or Delete to manage reconciliation policies.")
}
