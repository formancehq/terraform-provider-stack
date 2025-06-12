package server

import (
	"context"
	"net/http"

	"github.com/formancehq/go-libs/v3/httpclient"
	"github.com/formancehq/go-libs/v3/otlp"
	"github.com/formancehq/go-libs/v3/service"
	cloudpkg "github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/formancehq/terraform-provider-stack/pkg"

	"github.com/spf13/pflag"
	"go.uber.org/fx"
)

const (
	FormanceStackClientSecretKey = "formance-stack-client-secret"
	FormanceStackClientIdKey     = "formance-stack-client-id"
	FormanceStackEndpointKey     = "formance-stack-api-endpoint"
)

// AddFlags registers command-line flags for Formance Cloud or Stack client ID, client secret, and API endpoint configuration.
func AddFlags(flagset *pflag.FlagSet) {
	flagset.String(FormanceStackClientIdKey, "", "Client ID for Formance Cloud (user_%s) or Stack")
	flagset.String(FormanceStackClientSecretKey, "", "Client Secret for Formance Cloud or Stack")
	flagset.String(FormanceStackEndpointKey, "https://app.formance.cloud/api", "Endpoint for Formance Cloud")
}

// NewModule configures and returns an Fx module for integrating with Formance Cloud or Stack.
//
// It reads configuration values from the provided flag set, sets up HTTP transport with instrumentation and debugging, and supplies dependencies for the Formance SDK, token providers, stack provider, and API server. The module also registers a lifecycle hook to start the API server asynchronously and handle shutdown on failure.
func NewModule(ctx context.Context, flagset *pflag.FlagSet) fx.Option {
	clientId, _ := flagset.GetString(FormanceStackClientIdKey)
	clientSecret, _ := flagset.GetString(FormanceStackClientSecretKey)
	endpoint, _ := flagset.GetString(FormanceStackEndpointKey)
	debug, _ := flagset.GetBool(service.DebugFlag)
	transport := otlp.NewRoundTripper(http.DefaultTransport, debug)
	transport = httpclient.NewDebugHTTPTransport(transport)

	return fx.Options(
		fx.Supply(FormanceStackClientId(clientId)),
		fx.Supply(FormanceStackClientSecret(clientSecret)),
		fx.Supply(FormanceStackEndpoint(endpoint)),
		fx.Supply(fx.Annotate(transport, fx.As(new(http.RoundTripper)))),
		fx.Provide(sdk.NewCloudSDK),
		fx.Provide(pkg.NewTokenProviderFn),
		fx.Provide(cloudpkg.NewTokenProviderFactory),
		fx.Provide(sdk.NewStackSdk),
		fx.Provide(NewStackProvider),
		fx.Provide(NewAPI),
		fx.Invoke(func(lc fx.Lifecycle, server *API, shutdowner fx.Shutdowner) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					debug, _ := flagset.GetBool(service.DebugFlag)
					go func() {
						if err := server.Run(ctx, debug); err != nil {
							if err := shutdowner.Shutdown(); err != nil {
								panic(err)
							}
						}
					}()
					return nil
				},
			})
		}),
	)

}
