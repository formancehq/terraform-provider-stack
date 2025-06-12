package server

import (
	"context"
	"fmt"

	"github.com/formancehq/terraform-provider-stack/internal"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

type FormanceStackEndpoint string
type FormanceStackClientSecret string
type FormanceStackClientId string
type ProviderFactory func() provider.Provider

type API struct {
	factory ProviderFactory
}

func (a *API) Run(ctx context.Context, debug bool) error {
	opts := providerserver.ServeOpts{
		Address: fmt.Sprintf("%s/%s", "registry.terraform.io", internal.TerraformRepository),
		Debug:   debug,
	}

	err := providerserver.Serve(ctx, a.factory, opts)
	if err != nil {
		return err
	}

	return nil
}
// NewAPI creates a new API instance with the provided provider factory.
func NewAPI(factory ProviderFactory) *API {
	return &API{
		factory: factory,
	}
}
