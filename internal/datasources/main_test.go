package datasources_test

import (
	"context"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
)

func test(t *testing.T, fn func(ctx context.Context)) {
	t.Parallel()

	ctx := logging.TestingContext()

	fn(ctx)
}
