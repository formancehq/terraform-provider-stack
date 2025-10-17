package resources

import (
	"context"
	"fmt"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-stack/internal"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                     = &Ledger{}
	_ resource.ResourceWithConfigure        = &Ledger{}
	_ resource.ResourceWithConfigValidators = &Ledger{}
	_ resource.ResourceWithValidateConfig   = &Ledger{}
)

type Ledger struct {
	store *internal.ModuleStore
}

type LedgerModel struct {
	Name   types.String `tfsdk:"name"`
	Bucket types.String `tfsdk:"bucket"`
	// TODO: Handle features in the SDK
	// Features types.Map    `tfsdk:"features"`
	Metadata types.Map `tfsdk:"metadata"`
}

func NewLedger() func() resource.Resource {
	return func() resource.Resource {
		return &Ledger{}
	}
}

var SchemaLedger = schema.Schema{
	Description: "Resource for managing a Formance Ledger. For advanced usage and configuration, see the [Ledger documentation](https://docs.formance.com/ledger/).",
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required:    true,
			Description: "The name of the ledger.",
		},
		"bucket": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The bucket where the ledger data will be stored. If not provided, a default bucket will be used.",
		},
		// TODO: Handle features in the SDK
		// "features": schema.MapAttribute{
		// 	Optional:    true,
		// 	ElementType: types.StringType,
		// },
		"metadata": schema.MapAttribute{
			Optional:    true,
			Computed:    true,
			Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
			ElementType: types.StringType,
			Description: "Metadata associated with the ledger, stored as key-value pairs. Advanced usage: See [Ledger Advanced Filtering](https://docs.formance.com/ledger/advanced/filtering) and [Ledger documentation](https://docs.formance.com/ledger/) for more information.",
		},
	},
}

// Schema implements resource.Resource.
func (s *Ledger) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = SchemaLedger
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (s *Ledger) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, res *resource.ValidateConfigResponse) {
	var config LedgerModel
	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if res.Diagnostics.HasError() {
		return
	}

	if config.Name.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Invalid Ledger Name",
			"The ledger name cannot be null. Please provide a valid name.",
		)
		return
	}

}

// ConfigValidators implements resource.ResourceWithConfigValidators.
func (s *Ledger) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return nil
}

// Configure implements resource.ResourceWithConfigure.
func (s *Ledger) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

// Create implements resource.Resource.
func (s *Ledger) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan LedgerModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	config := operations.V2CreateLedgerRequest{}
	config.Ledger = plan.Name.ValueString()
	if !plan.Bucket.IsNull() {
		config.V2CreateLedgerRequest.Bucket = pointer.For(plan.Bucket.ValueString())
	}

	// TODO: Handle features in the SDK
	// if !plan.Features.IsNull() {
	// 	config.V2CreateLedgerRequest.Features = collectionutils.ConvertMap(plan.Features.Elements(), func(v attr.Value) string {
	// 		return v.String()
	// 	})
	// }
	if !plan.Metadata.IsNull() {
		config.V2CreateLedgerRequest.Metadata = collectionutils.ConvertMap(plan.Metadata.Elements(), func(v attr.Value) string {
			return v.(types.String).ValueString()
		})
	}

	ledgerSdk := s.store.Ledger()
	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}
	_, err := ledgerSdk.CreateLedger(ctx, config)
	if err != nil {

		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}

	l, err := ledgerSdk.GetLedger(ctx, operations.V2GetLedgerRequest{
		Ledger: plan.Name.ValueString(),
	})
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}

	data := l.V2GetLedgerResponse.Data
	plan.Name = types.StringValue(data.Name)
	plan.Bucket = types.StringValue(data.Bucket)
	plan.Metadata = types.MapValueMust(types.StringType,
		collectionutils.ConvertMap(data.Metadata, func(v string) attr.Value {
			return types.StringValue(v)
		}),
	)

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (s *Ledger) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state LedgerModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}
}

// Metadata implements resource.Resource.
func (s *Ledger) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ledger"
}

// Read implements resource.Resource.
func (s *Ledger) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state LedgerModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	ledgerSdk := s.store.Ledger()
	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}
	ledger, err := ledgerSdk.GetLedger(ctx, operations.V2GetLedgerRequest{
		Ledger: state.Name.ValueString(),
	})
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}

	data := ledger.V2GetLedgerResponse.Data
	state.Name = types.StringValue(data.Name)
	state.Bucket = types.StringValue(data.Bucket)
	state.Metadata = types.MapValueMust(types.StringType,
		collectionutils.ConvertMap(data.Metadata, func(v string) attr.Value {
			return types.StringValue(v)
		}),
	)
	// TODO: Handle features, update the sdk
	// state.Features = types.MapValueMust(types.StringType, collectionutils.ConvertMap(data.Features, func(v string) attr.Value {
	// 	return types.StringValue(v)
	// }))

	res.Diagnostics.Append(res.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (s *Ledger) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var state LedgerModel
	var plan LedgerModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	ledgerSdk := s.store.Ledger()

	s.store.CheckModuleHealth(ctx, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}
	_, err := ledgerSdk.UpdateLedgerMetadata(ctx, operations.V2UpdateLedgerMetadataRequest{
		Ledger: state.Name.ValueString(),
		RequestBody: collectionutils.ConvertMap(plan.Metadata.Elements(), func(v attr.Value) string {
			return v.(types.String).ValueString()
		}),
	})
	if err != nil {
		sdk.HandleStackError(ctx, err, &res.Diagnostics)
		return
	}

	state.Metadata = plan.Metadata

	res.Diagnostics.Append(res.State.Set(ctx, &state)...)

}
