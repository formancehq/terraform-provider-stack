package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/formancehq/go-libs/v3/logging"
	cloudpkg "github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-stack/internal"
	"github.com/formancehq/terraform-provider-stack/internal/resources"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/formancehq/terraform-provider-stack/pkg"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FormanceCloudProviderModel struct {
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Endpoint     types.String `tfsdk:"endpoint"`
}

type ProviderModelAdapter struct {
	m *FormanceCloudProviderModel
}

func NewProviderModelAdapter(m *FormanceCloudProviderModel) *ProviderModelAdapter {
	return &ProviderModelAdapter{
		m: m,
	}
}

func (f *ProviderModelAdapter) ClientId() string {
	return f.m.ClientId.ValueString()
}
func (f *ProviderModelAdapter) ClientSecret() string {
	return f.m.ClientSecret.ValueString()
}
func (f *ProviderModelAdapter) Endpoint() string {
	return f.m.Endpoint.ValueString()
}

func IsOrganizationClient(clientId string) bool {
	return strings.HasPrefix(clientId, "organization_")
}

func (f *ProviderModelAdapter) UserAgent() string {
	return fmt.Sprintf("terraform-provider-stack/%s", internal.Version)
}

type FormanceStackProviderModel struct {
	Cloud *FormanceCloudProviderModel `tfsdk:"cloud"`

	StackId        types.String `tfsdk:"stack_id"`
	OrganizationId types.String `tfsdk:"organization_id"`
	Uri            types.String `tfsdk:"uri"`

	WaitModule types.String `tfsdk:"wait_module_duration"`
}

type FormanceStackProvider struct {
	logger            logging.Logger
	transport         http.RoundTripper
	cloudFactory      sdk.CloudFactory
	cloudtokenFactory cloudpkg.TokenProviderFactory
	stackTokenFactory pkg.TokenProviderFactory
	stackSdkFactory   sdk.StackSdkFactory

	Version  string
	Endpoint string

	ClientId     string
	ClientSecret string
}

var SchemaStack = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"stack_id": schema.StringAttribute{
			Required: true,
		},
		"organization_id": schema.StringAttribute{
			Required: true,
		},
		"uri": schema.StringAttribute{
			Required: true,
		},
		"cloud": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"client_secret": schema.StringAttribute{
					Optional:  true,
					Sensitive: true,
				},
				"client_id": schema.StringAttribute{
					Optional: true,
				},
				"endpoint": schema.StringAttribute{
					Optional: true,
				},
			},
		},
		"wait_module_duration": schema.StringAttribute{
			Optional:    true,
			Description: "The duration to wait for the module to be ready before proceeding.",
		},
	},
}

// Metadata satisfies the provider.Provider interface for FormanceCloudProvider
func (p *FormanceStackProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "formancestack"
	resp.Version = internal.Version
}

// Schema satisfies the provider.Provider interface for FormanceCloudProvider.
func (p *FormanceStackProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = SchemaStack
}

// Configure satisfies the provider.Provider interface for FormanceCloudProvider.
func (p *FormanceStackProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	p.logger.Debugf("Configuring stack provider version %s", p.Version)
	defer p.logger.Debugf("Stack provider configured with client_id: %s, endpoint: %s", p.ClientId, p.Endpoint)
	ctx = logging.ContextWithLogger(ctx, p.logger)
	var data FormanceStackProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Cloud == nil {
		data.Cloud = &FormanceCloudProviderModel{
			ClientId:     types.StringValue(p.ClientId),
			ClientSecret: types.StringValue(p.ClientSecret),
			Endpoint:     types.StringValue(p.Endpoint),
		}
	}

	if data.Cloud.ClientId.ValueString() == "" {
		if p.ClientId != "" {
			data.Cloud.ClientId = types.StringValue(p.ClientId)
		}
	}

	if data.Cloud.ClientSecret.ValueString() == "" {
		if p.ClientSecret != "" {
			data.Cloud.ClientSecret = types.StringValue(p.ClientSecret)
		}
	}

	if data.Cloud.Endpoint.ValueString() == "" {
		data.Cloud.Endpoint = types.StringValue(p.Endpoint)
	}
	p.logger.Debugf("Stack provider configuration: %+v", data)

	tfcreds := &FormanceCloudProviderModel{
		ClientId:     data.Cloud.ClientId,
		ClientSecret: data.Cloud.ClientSecret,
		Endpoint:     data.Cloud.Endpoint,
	}

	creds := NewProviderModelAdapter(tfcreds)
	if !IsOrganizationClient(creds.ClientId()) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid client_id: %s", creds.ClientId()),
			"The client_id must start with 'organization_' to be used with the stack provider. "+
				"Please check your configuration and try again.",
		)
		return
	}

	if data.StackId.IsUnknown() || data.OrganizationId.IsUnknown() || data.Uri.IsUnknown() {
		return
	}

	cloudtp := p.cloudtokenFactory(p.transport, creds)
	cloudSDK := p.cloudFactory(creds, cloudpkg.NewTransport(p.transport, cloudtp))
	stack := pkg.Stack{
		Id:             data.StackId.ValueString(),
		OrganizationId: data.OrganizationId.ValueString(),
		Uri:            data.Uri.ValueString(),
	}
	stackTpProvider := p.stackTokenFactory(p.transport, creds, cloudtp, stack)
	cli, err := p.stackSdkFactory(data.Uri.ValueString(), p.Version, p.transport, stackTpProvider)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create stack client", err.Error())
		return
	}

	store := internal.Store{
		Stack:             stack,
		StackSdkImpl:      cli,
		CloudSDK:          cloudSDK,
		WaitModuleTimeout: 2 * time.Minute,
	}

	if data.WaitModule.ValueString() != "" {
		duration, err := time.ParseDuration(data.WaitModule.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid wait_module_duration",
				fmt.Sprintf("The provided wait_module_duration '%s' is not a valid duration: %s", data.WaitModule.ValueString(), err.Error()),
			)
			return
		}
		store.WaitModuleTimeout = duration
	}

	resp.ResourceData = store
	resp.DataSourceData = store
}

// DataSources satisfies the provider.Provider interface for FormanceCloudProvider.
func (p *FormanceStackProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources satisfies the provider.Provider interface for FormanceCloudProvider.
func (p *FormanceStackProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewWebhooks(p.logger.WithField("resource", "webhooks")),
		// resources.NewLedger(p.logger.WithField("resource", "ledger")),
		resources.NewNoop(p.logger.WithField("resource", "noop")),
	}
}

func (p FormanceStackProvider) ConfigValidators(ctx context.Context) []provider.ConfigValidator {
	return []provider.ConfigValidator{}
}

func (p FormanceStackProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	var data FormanceStackProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if data.Cloud == nil {
		data.Cloud = &FormanceCloudProviderModel{}
	}

	if data.Cloud.ClientId.ValueString() == "" {
		if p.ClientId != "" {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("cloud.client_id"),
				"Missing client_id Configuration",
				"While configuring the provider, the client_id was not found "+
					"However the FORMANCE_CLOUD_CLIENT_ID environment variable was set ",
			)
		} else {
			resp.Diagnostics.AddAttributeError(
				path.Root("cloud.client_id"),
				"Missing Client ID Configuration",
				"While configuring the provider, the client id was not found. "+
					"the FORMANCE_CLOUD_CLIENT_ID environment variable or provider "+
					"configuration block client_id attribute.",
			)
		}
	}

	if data.Cloud.ClientSecret.ValueString() == "" {
		if p.ClientSecret != "" {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("cloud.client_secret"),
				"Missing client_secret Configuration",
				"While configuring the provider, the client_secret was not found "+
					"in the configuration. However, the FORMANCE_CLOUD_CLIENT_SECRET environment variable was found.",
			)
		} else {
			resp.Diagnostics.AddAttributeError(
				path.Root("cloud.client_secret"),
				"Missing API Token Configuration",
				"While configuring the provider, the API token was not found in "+
					"the FORMANCE_CLOUD_CLIENT_SECRET environment variable or provider "+
					"configuration block api_token attribute.",
			)
		}
	}

	if data.Cloud.Endpoint.ValueString() == "" {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("cloud.endpoint"),
			fmt.Sprintf("Missing Endpoint Configuration use %s", p.Endpoint),
			"While configuring the provider, the endpoint was not found "+
				"However the FORMANCE_CLOUD_API_ENDPOINT environment variable was set",
		)
	}

	if data.StackId.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("stack_id"),
			"Missing Stack ID Configuration",
			"While configuring the provider, the stack_id was not found. "+
				"Please provide a valid stack_id in the provider configuration block.",
		)
	}

	if data.OrganizationId.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("organization_id"),
			"Missing Organization ID Configuration",
			"While configuring the provider, the organization_id was not found. "+
				"Please provide a valid organization_id in the provider configuration block.",
		)
	}

	if data.Uri.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("uri"),
			"Missing URI Configuration",
			"While configuring the provider, the url was not found. "+
				"Please provide a valid uri in the provider configuration block.",
		)
	}

}

func NewStackProvider(
	logger logging.Logger,
	endpoint FormanceStackEndpoint,
	clientId FormanceStackClientId,
	clientSecret FormanceStackClientSecret,
	transport http.RoundTripper,
	cloudSdkFactory sdk.CloudFactory,
	cloudTokenFactory cloudpkg.TokenProviderFactory,
	stackTokenFactory pkg.TokenProviderFactory,
	stackSdkFactory sdk.StackSdkFactory,
) ProviderFactory {
	return func() provider.Provider {
		return &FormanceStackProvider{
			logger:            logger,
			ClientId:          string(clientId),
			ClientSecret:      string(clientSecret),
			Endpoint:          string(endpoint),
			transport:         transport,
			cloudFactory:      cloudSdkFactory,
			cloudtokenFactory: cloudTokenFactory,
			stackTokenFactory: stackTokenFactory,
			stackSdkFactory:   stackSdkFactory,
		}
	}
}

var _ provider.ProviderWithConfigValidators = &FormanceStackProvider{}
var _ provider.ProviderWithValidateConfig = &FormanceStackProvider{}
var _ provider.Provider = &FormanceStackProvider{}
