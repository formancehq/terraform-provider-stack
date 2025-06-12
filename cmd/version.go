package cmd

import (
	"context"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-stack/internal"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// NewVersion creates a Cobra command that logs and displays the Terraform provider's version information, then shuts down the application.
func NewVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  `All software has versions. This is Terraform provider's`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SetContext(contextWithFxOpts(cmd.Context(), fx.Options(
				fxOptsFromContext(cmd.Context()),
				fx.Invoke(func(lc fx.Lifecycle, appInfo internal.AppInfo, shutdowner fx.Shutdowner) {
					lc.Append(fx.StartHook(func(ctx context.Context) error {
						logging.FromContext(ctx).Infof(appInfo.String())
						return shutdowner.Shutdown()
					}))
				}),
			)))
			return nil
		},
	}
}
