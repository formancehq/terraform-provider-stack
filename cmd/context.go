package cmd

import (
	"context"

	"go.uber.org/fx"
)

type fxOptions struct{}

// contextWithFxOpts returns a new context containing the provided Fx option.
func contextWithFxOpts(ctx context.Context, opts fx.Option) context.Context {
	return context.WithValue(ctx, fxOptions{}, opts)
}

// fxOptsFromContext retrieves the Fx option stored in the context, or returns an empty Fx option if none is found.
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
