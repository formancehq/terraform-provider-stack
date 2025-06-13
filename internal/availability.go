package internal

import (
	"context"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"

	"fmt"
)

func CheckModuleHealth(ctx context.Context, sdk sdk.StackSdkImpl, moduleName string) error {
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
