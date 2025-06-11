package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/go-libs/v3/httpclient"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/otlp"
	cloudpkg "github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/formancehq/terraform-provider-stack/internal"
	"github.com/formancehq/terraform-provider-stack/internal/resources"
	"github.com/formancehq/terraform-provider-stack/pkg"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func (f *ProviderModelAdapter) IsOrganizationClient() bool {
	return strings.HasPrefix(f.ClientId(), "organization_")
}

func (f *ProviderModelAdapter) UserAgent() string {
	return fmt.Sprintf("terraform-provider-stack/%s", internal.Version)
}

type FormanceStackProviderModel struct {
	Cloud *FormanceCloudProviderModel `tfsdk:"cloud"`

	StackId        types.String `tfsdk:"stack_id"`
	OrganizationId types.String `tfsdk:"organization_id"`
	Uri            types.String `tfsdk:"uri"`

	ExpectedModules types.List `tfsdk:"expected_modules"`
}

type FormanceStackProvider struct {
	logger logging.Logger

	Version  string
	Endpoint string

	ClientId     string
	ClientSecret string

	SDKFactory cloudpkg.SDKFactory
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
		"expected_modules": schema.ListAttribute{
			Optional:    true,
			ElementType: types.StringType,
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
	if !creds.IsOrganizationClient() {
		resp.Diagnostics.AddError(
			"Invalid Client ID",
			"The client_id must start with 'organization_' to be used with the stack provider. "+
				"Please check your configuration and try again.",
		)
		return
	}

	if data.StackId.IsUnknown() || data.OrganizationId.IsUnknown() || data.Uri.IsUnknown() {
		return
	}

	cloud, tp := p.SDKFactory(creds)
	resp.Diagnostics.Append(p.pollStack(ctx, cloud, data.OrganizationId.ValueString(), data.StackId.ValueString())...)
	if resp.Diagnostics.HasError() {
		return
	}

	stackTpProvider := pkg.NewTokenProvider(&http.Client{
		Transport: otlp.NewRoundTripper(httpclient.NewDebugHTTPTransport(http.DefaultTransport), true),
	}, creds)
	stackToken, err := stackTpProvider.StackSecurityToken(ctx, tp, pkg.Stack{
		Id:             data.StackId.ValueString(),
		OrganizationId: data.OrganizationId.ValueString(),
		Uri:            data.Uri.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to get stack security token", err.Error())
		return
	}

	cli, err := pkg.NewStackClient(data.Uri.ValueString(), p.Version, stackToken)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create stack client", err.Error())
		return
	}

	resp.ResourceData = cli
	resp.DataSourceData = cli
}

// TODO: refactor to send continuous diagnostics
func (p *FormanceStackProvider) pollStack(ctx context.Context, cli sdk.DefaultAPI, organizationId, stackId string, expectedModules ...string) diag.Diagnostics {
	pollctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()
	diags := diag.Diagnostics{}
	for {
		select {
		case <-pollctx.Done():
			diags.Append(diag.NewErrorDiagnostic("Timeout while waiting for stack to be ready", "The stack did not become ready within the timeout period. Please check your configuration and try again."))
			return diags
		case <-time.After(2 * time.Second):
			stack, res, err := cli.GetStack(ctx, organizationId, stackId).Execute()
			if err != nil {
				cloudpkg.HandleSDKError(ctx, err, res, &diags)
				return diags
			}

			if stack.Data.State != "ACTIVE" {
				diags.AddError("Stack is not active", fmt.Sprintf("The stack is currently in state '%s'. Please wait until it is active.", stack.Data.State))
				return diags
			}

			if stack.Data.Status == "READY" {
				if len(expectedModules) > 0 {
					modules, res, err := cli.ListModules(ctx, organizationId, stackId).Execute()
					if err != nil {
						cloudpkg.HandleSDKError(ctx, err, res, &diags)
						return diags
					}

					remaining := collectionutils.Filter(modules.Data, func(module sdk.Module) bool {
						return collectionutils.Contains(expectedModules, module.Name) && module.Status != "READY"
					})

					if len(remaining) > 0 {
						continue
					}
				}
				return nil
			}

			p.logger.Debugf("Stack %s is not ready yet, current status: %s", stackId, stack.Data.Status)
		}
	}
}

// DataSources satisfies the provider.Provider interface for FormanceCloudProvider.
func (p *FormanceStackProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources satisfies the provider.Provider interface for FormanceCloudProvider.
func (p *FormanceStackProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewWebhooks(p.logger.WithField("resource", "webhooks")),
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
				"While configuring the provider, the client_secret was not found in "+
					"however the environment variable FORMANCE_CLOUD_CLIENT_SECRET was found ",
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
				"Please provide a valid url in the provider configuration block.",
		)
	}

}

func NewStackProvider(
	logger logging.Logger,
	endpoint FormanceStackEndpoint,
	clientId FormanceStackClientId,
	clientSecret FormanceStackClientSecret,
	sdkFactory cloudpkg.SDKFactory,
) func() provider.Provider {
	return func() provider.Provider {
		return &FormanceStackProvider{
			logger:       logger,
			ClientId:     string(clientId),
			ClientSecret: string(clientSecret),
			Endpoint:     string(endpoint),
			SDKFactory:   sdkFactory,
		}
	}
}

var _ provider.ProviderWithConfigValidators = &FormanceStackProvider{}
var _ provider.ProviderWithValidateConfig = &FormanceStackProvider{}
var _ provider.Provider = &FormanceStackProvider{}
