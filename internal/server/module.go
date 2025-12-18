package server

import (
	"context"
	"net/http"

	v3 "github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/retry"
	"github.com/formancehq/go-libs/v3/httpclient"
	"github.com/formancehq/go-libs/v3/otlp"
	"github.com/formancehq/go-libs/v3/service"
	cloudpkg "github.com/formancehq/terraform-provider-cloud/pkg"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	cloudretry "github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/retry"
	cloudspeakeasyretry "github.com/formancehq/terraform-provider-cloud/pkg/speakeasy_retry"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/formancehq/terraform-provider-stack/pkg"
	speakeasyretry "github.com/formancehq/terraform-provider-stack/pkg/speakeasy_retry"

	"github.com/spf13/pflag"
	"go.uber.org/fx"
)

const (
	FormanceStackClientSecretKey = "formance-stack-client-secret"
	FormanceStackClientIdKey     = "formance-stack-client-id"
	FormanceStackEndpointKey     = "formance-stack-api-endpoint"
)

func AddFlags(flagset *pflag.FlagSet) {
	flagset.String(FormanceStackClientIdKey, "", "Client ID for Formance Cloud (organization_%s) or Stack")
	flagset.String(FormanceStackClientSecretKey, "", "Client Secret for Formance Cloud or Stack")
	flagset.String(FormanceStackEndpointKey, "https://app.formance.cloud/api", "Endpoint for Formance Cloud Or Stack Auth module")
	cloudspeakeasyretry.AddFlags(flagset)
	speakeasyretry.AddFlags(flagset)
}

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
		fx.Provide(pkg.NewTokenProviderFn),
		cloudspeakeasyretry.NewModule(flagset),
		speakeasyretry.NewModule(flagset),
		fx.Provide(func(retry *cloudretry.Config) sdk.CloudFactory {
			opts := []membershipclient.SDKOption{}
			if retry != nil {
				opts = append(opts, membershipclient.WithRetryConfig(*retry))
			}
			return sdk.NewCloudSDK(
				opts...,
			)
		}),
		fx.Provide(func() cloudpkg.TokenProviderFactory {
			return cloudpkg.NewTokenProvider
		}), fx.Provide(func(retry *retry.Config) sdk.StackSdkFactory {
			opts := []v3.SDKOption{}
			if retry != nil {
				opts = append(opts, v3.WithRetryConfig(*retry))
			}
			return sdk.NewStackSdk(
				opts...,
			)
		}),
		fx.Provide(NewStackProvider),
		fx.Provide(NewAPI),
		fx.Invoke(func(lc fx.Lifecycle, server *API, shutdowner fx.Shutdowner) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
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
