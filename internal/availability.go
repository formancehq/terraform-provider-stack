package internal

import (
	"context"

	formance "github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/collectionutils"

	"fmt"
)

// CheckModuleHealth verifies the health status of a specified module in the current stack.
// It returns an error if the module is not healthy or if the stack version information cannot be retrieved.
func CheckModuleHealth(ctx context.Context, sdk *formance.Formance, moduleName string) error {
	stackInfo, err := sdk.GetVersions(ctx)
	if err != nil {
		return fmt.Errorf("unable to get stack /versions: %w", err)
	}

	moduleVersions := stackInfo.GetVersionsResponse.Versions

	version := collectionutils.First(moduleVersions, func(v shared.Version) bool {
		return v.Name == moduleName
	})

	if !version.Health {
		return fmt.Errorf("%s module is not healthy: %s", moduleName, version.Version)
	}

	return nil
}
