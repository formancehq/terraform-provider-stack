package resources

import (
	"context"
	"fmt"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/terraform-provider-stack/internal"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &PaymentsPool{}
	_ resource.ResourceWithConfigure = &PaymentsPool{}
)

type PaymentsPool struct {
	store *internal.ModuleStore
}

type PaymentsPoolModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	AccountsIds types.List   `tfsdk:"accounts_ids"`
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
	},
}

// Schema implements resource.Resource.
func (s *PaymentsPool) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = SchemaPaymentsPool
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

	sdkPayments := s.store.Payments()
	resp, err := sdkPayments.CreatePool(ctx, &shared.V3CreatePoolRequest{
		Name: plan.Name.String(),
		AccountIDs: collectionutils.Map(plan.AccountsIds.Elements(), func(account attr.Value) string {
			return account.(types.String).ValueString()
		}),
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
	state.AccountsIds = types.ListValueMust(
		types.StringType,
		collectionutils.Map(resp.V3GetPoolResponse.Data.PoolAccounts, func(account string) attr.Value {
			return types.StringValue(account)
		}),
	)

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

	plan.ID = state.ID

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}
