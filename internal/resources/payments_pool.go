package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/go-libs/v3/query"
	"github.com/formancehq/terraform-provider-stack/internal"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/compare"
)

var (
	_ resource.Resource                     = &PaymentsPool{}
	_ resource.ResourceWithConfigure        = &PaymentsPool{}
	_ resource.ResourceWithValidateConfig   = &PaymentsPool{}
	_ resource.ResourceWithConfigValidators = &PaymentsPool{}
)

type PaymentsPool struct {
	store *internal.ModuleStore
}

type PaymentsPoolModel struct {
	ID          types.String  `tfsdk:"id"`
	Name        types.String  `tfsdk:"name"`
	AccountsIds types.List    `tfsdk:"accounts_ids"`
	Query       types.Dynamic `tfsdk:"query"`
}

func NewPaymentsPool() func() resource.Resource {
	return func() resource.Resource {
		return &PaymentsPool{}
	}
}

var SchemaPaymentsPool = schema.Schema{
	Description: "Resource for managing a Formance Payments Pool. For advanced usage and configuration, see the [Payments documentation](https://docs.formance.com/payments/).",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "The unique identifier of the payments pool.",
		},
		"name": schema.StringAttribute{
			Description: "The name of the pool.",
			Required:    true,
		},
		"accounts_ids": schema.ListAttribute{
			Description: "The list of accounts IDs associated with the pool. For more information, see the [Payments documentation](https://docs.formance.com/payments/).",
			ElementType: types.StringType,
			Optional:    true,
		},
		"query": schema.DynamicAttribute{
			Description: "The query to filter payments associated with the pool. For more information, see the [Payments documentation](https://docs.formance.com/payments/).",
			Optional:    true,
		},
	},
}

// ConfigValidators implements resource.ResourceWithConfigValidators.
func (s *PaymentsPool) ConfigValidators(context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("accounts_ids"),
			path.MatchRoot("query"),
		),
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("accounts_ids"),
			path.MatchRoot("query"),
		),
	}
}

// Schema implements resource.Resource.
func (s *PaymentsPool) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = SchemaPaymentsPool
}

func (m *PaymentsPoolModel) ParseQuery() (map[string]any, error) {
	object, ok := m.Query.UnderlyingValue().(types.Object)
	if !ok {
		return nil, nil
	}
	qb, err := query.ParseJSON(object.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}
	if qb == nil {
		return nil, nil
	}
	var query map[string]any
	if err := json.Unmarshal([]byte(m.Query.String()), &query); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ledger query: %w", err)
	}
	return query, nil
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (s *PaymentsPool) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, res *resource.ValidateConfigResponse) {
	var conf PaymentsPoolModel
	res.Diagnostics.Append(req.Config.Get(ctx, &conf)...)
	if res.Diagnostics.HasError() {
		return
	}

	if !conf.Query.IsNull() {
		if _, ok := conf.Query.UnderlyingValue().(types.Object); !ok {
			res.Diagnostics.AddError("Invalid Ledger Query", "The ledger_query must be a valid JSON object.")
		} else {
			_, err := conf.ParseQuery()
			if err != nil {
				res.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("Failed to read payments pool policy query: %s", err))
			}
		}
	}

}

// Configure implements resource.ResourceWithConfigure.
func (s *PaymentsPool) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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
func (s *PaymentsPool) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan PaymentsPoolModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	query, err := plan.ParseQuery()
	if err != nil {
		res.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("Failed to read payments pool policy query: %s", err))
		return
	}

	sdkPayments := s.store.Payments()
	resp, err := sdkPayments.CreatePool(ctx, &shared.V3CreatePoolRequest{
		Name: plan.Name.ValueString(),
		AccountIDs: collectionutils.Map(plan.AccountsIds.Elements(), func(account attr.Value) string {
			return account.(types.String).ValueString()
		}),
		Query: query,
	})
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}

	plan.ID = types.StringValue(resp.V3CreatePoolResponse.Data)

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (s *PaymentsPool) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state PaymentsPoolModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	sdkPaymentsPool := s.store.Payments()
	_, err := sdkPaymentsPool.DeletePool(ctx, operations.V3DeletePoolRequest{
		PoolID: state.ID.ValueString(),
	})
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}
}

// Metadata implements resource.Resource.
func (s *PaymentsPool) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_payments_pool"
}

// Read implements resource.Resource.
func (s *PaymentsPool) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state PaymentsPoolModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	sdkPaymentsPool := s.store.Payments()
	resp, err := sdkPaymentsPool.GetPool(ctx, operations.V3GetPoolRequest{
		PoolID: state.ID.ValueString(),
	})
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}

	state.ID = types.StringValue(resp.V3GetPoolResponse.Data.ID)
	state.Name = types.StringValue(resp.V3GetPoolResponse.Data.Name)
	if len(resp.V3GetPoolResponse.Data.PoolAccounts) > 0 {
		state.AccountsIds = types.ListValueMust(
			types.StringType,
			collectionutils.Map(resp.V3GetPoolResponse.Data.PoolAccounts, func(account string) attr.Value {
				return types.StringValue(account)
			}),
		)
	}

	query := resp.V3GetPoolResponse.Data.Query
	if len(query) > 0 {
		tfValues := ConvertToAttrValues(query)
		state.Query = types.DynamicValue(NewDynamicObjectValue(tfValues).Value())
	}

	res.Diagnostics.Append(res.State.Set(ctx, &state)...)
}

func diff(a, b []string) []string {
	temp := map[string]int{}
	for _, s := range a {
		temp[s]++
	}
	for _, s := range b {
		temp[s]--
	}

	var result []string
	for s, v := range temp {
		if v > 0 {
			result = append(result, s)
		}
	}
	return result
}

// Update implements resource.Resource.
func (s *PaymentsPool) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var plan PaymentsPoolModel
	var state PaymentsPoolModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	planAccountIds := collectionutils.Map(plan.AccountsIds.Elements(), func(account attr.Value) string {
		return account.(types.String).ValueString()
	})
	stateAccountIds := collectionutils.Map(state.AccountsIds.Elements(), func(account attr.Value) string {
		return account.(types.String).ValueString()
	})

	accountToAdd := diff(planAccountIds, stateAccountIds)
	accountToRemove := diff(stateAccountIds, planAccountIds)

	for _, accountID := range accountToAdd {
		_, err := s.store.Payments().AddAccountToPool(ctx, operations.V3AddAccountToPoolRequest{
			PoolID:    state.ID.ValueString(),
			AccountID: accountID,
		})
		if err != nil {
			sdk.HandleStackError(ctx, err, &res.Diagnostics)
			return
		}
	}

	for _, accountID := range accountToRemove {
		_, err := s.store.Payments().RemoveAccountFromPool(ctx, operations.V3RemoveAccountFromPoolRequest{
			PoolID:    state.ID.ValueString(),
			AccountID: accountID,
		})
		if err != nil {
			sdk.HandleStackError(ctx, err, &res.Diagnostics)
			return
		}
	}

	planQuery, err := plan.ParseQuery()
	if err != nil {
		res.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("Failed to read payments pool policy query: %s", err))
		return
	}
	_, err = state.ParseQuery()
	if err != nil {
		res.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("Failed to read payments pool policy query: %s", err))
		return
	}
	if err := compare.ValuesDiffer().CompareValues(plan.Query, state.Query); err == nil {
		_, err := s.store.Payments().UpdatePool(ctx, operations.V3UpdatePoolQueryRequest{
			V3UpdatePoolQueryRequest: &shared.V3UpdatePoolQueryRequest{
				Query: planQuery,
			},
			PoolID: state.ID.String(),
		})
		if err != nil {
			sdk.HandleStackError(ctx, err, &res.Diagnostics)
			return
		}
	}

	plan.ID = state.ID

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}
