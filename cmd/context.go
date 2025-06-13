package cmd

import (
	"context"

	"go.uber.org/fx"
)

type fxOptions struct{}

func contextWithFxOpts(ctx context.Context, opts fx.Option) context.Context {
	return context.WithValue(ctx, fxOptions{}, opts)
}

func fxOptsFromContext(ctx context.Context) fx.Option {
	opts := ctx.Value(fxOptions{})
	if opts == nil {
		return fx.Options()
	}

	if opts, ok := opts.(fx.Option); ok {
		return opts
	}
	return fx.Options()
}
