package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/formancehq/terraform-provider-stack/pkg"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type Store struct {
	Stack pkg.Stack
	sdk.StackSdkImpl
	sdk.CloudSDK

	WaitModuleTimeout time.Duration
}

func (s *Store) NewModuleStore(module string) *ModuleStore {
	return &ModuleStore{
		module:            module,
		Stack:             s.Stack,
		StackSdkImpl:      s.StackSdkImpl,
		WaitModuleTimeout: s.WaitModuleTimeout,
	}
}

type ModuleStore struct {
	module string
	Stack  pkg.Stack
	sdk.StackSdkImpl

	WaitModuleTimeout time.Duration
}

func (ms *ModuleStore) CheckModuleHealth(ctx context.Context, diagnostics *diag.Diagnostics) {
	for {
		select {
		case <-ctx.Done():
			diagnostics.AddError("Module Health Check Cancelled",
				fmt.Sprintf("The module '%s' health check was cancelled: %s", ms.module, ctx.Err()))
			return
		case <-time.After(ms.WaitModuleTimeout):
			diagnostics.AddError("Module Health Check Timeout",
				fmt.Sprintf("The module '%s' did not become healthy within the timeout period of %f seconds.", ms.module, ms.WaitModuleTimeout.Minutes()))
			return
		case <-time.After(time.Second * 2):
			modules, err := ms.GetVersions(ctx)
			if err != nil {
				if modules == nil || modules.StatusCode > 300 {
					continue
				}
				sdk.HandleStackError(ctx, err, diagnostics)
				return
			}

			module := collectionutils.First(modules.GetVersionsResponse.Versions, func(v shared.Version) bool {
				return v.Name == ms.module
			})

			if module.Name == "" {
				continue
			}

			if !module.Health {
				continue
			}

			return
		}
	}
}
