package server

import (
	"context"

	"github.com/formancehq/go-libs/v3/service"
	cloudpkg "github.com/formancehq/terraform-provider-cloud/pkg"

	"github.com/spf13/pflag"
	"go.uber.org/fx"
)

const (
	FormanceStackClientSecretKey = "formance-stack-client-secret"
	FormanceStackClientIdKey     = "formance-stack-client-id"
	FormanceStackEndpointKey     = "formance-stack-api-endpoint"
)

func AddFlags(flagset *pflag.FlagSet) {
	flagset.String(FormanceStackClientIdKey, "", "Client ID for Formance Cloud (user_%s) or Stack")
	flagset.String(FormanceStackClientSecretKey, "", "Client Secret for Formance Cloud or Stack")
	flagset.String(FormanceStackEndpointKey, "https://app.formance.cloud/api", "Endpoint for Formance Cloud")
}

func NewModule(ctx context.Context, flagset *pflag.FlagSet) fx.Option {
	clientId, _ := flagset.GetString(FormanceStackClientIdKey)
	clientSecret, _ := flagset.GetString(FormanceStackClientSecretKey)
	endpoint, _ := flagset.GetString(FormanceStackEndpointKey)
	return fx.Options(
		fx.Supply(FormanceStackClientId(clientId)),
		fx.Supply(FormanceStackClientSecret(clientSecret)),
		fx.Supply(FormanceStackEndpoint(endpoint)),
		fx.Provide(
			func() cloudpkg.SDKFactory {
				return cloudpkg.NewSDK
			},
		),
		fx.Provide(NewStackProvider),
		fx.Provide(
			NewAPI,
		),
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
