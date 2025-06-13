package cmd

import (
	"io"
	"os"
	"path"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/otlp"
	"github.com/formancehq/go-libs/v3/otlp/otlptraces"
	"github.com/formancehq/go-libs/v3/service"
	"github.com/formancehq/terraform-provider-stack/internal"
	"github.com/formancehq/terraform-provider-stack/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/fx"
)

func init() {
	cobra.EnableTraverseRunHooks = true
}

func Execute() {
	app := &App{}
	service.Execute(app.CobraCommand())
}

type App struct{}

func (a *App) Flags(pflags *pflag.FlagSet) {
	pflags.Bool(service.DebugFlag, false, "Debug mode")
	pflags.Bool(logging.JsonFormattingLoggerFlag, true, "Format logs as json")
	pflags.Duration(service.GracePeriodFlag, 0, "Grace period for shutdown")
	otlp.AddFlags(pflags)
	otlptraces.AddFlags(pflags)
}

func (a *App) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		NewVersion(),
	}
}

func (a *App) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                internal.ServiceName,
		Short:              "Formance cloud terraform provider",
		Long:               "Formance cloud terraform provider server and CLI",
		PersistentPreRunE:  a.preRunE,
		RunE:               runServe,
		PersistentPostRunE: a.postRunE,
	}

	a.Flags(cmd.PersistentFlags())
	server.AddFlags(cmd.Flags())
	cmd.PersistentFlags().Bool(LogDebugFlag, false, "Enable debug logging")

	cmd.AddCommand(a.SubCommands()...)

	return cmd
}

func logToFile() (io.Writer, error) {
	formanceDir := path.Join(os.Getenv("HOME"), ".formance")
	if err := os.MkdirAll(formanceDir, 0755); err != nil {
		return nil, err
	}

	logFile := path.Join(formanceDir, "tf-cloud-provider.log")
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	if err := file.Chmod(0644); err != nil {
		return nil, err
	}

	return file, nil
}

const (
	LogDebugFlag = "log-debug"
)

func (app *App) preRunE(cmd *cobra.Command, args []string) error {
	file, err := logToFile()
	if err != nil {
		return err
	}

	json, _ := cmd.Flags().GetBool(logging.JsonFormattingLoggerFlag)
	otelTraces, _ := cmd.Flags().GetString(otlptraces.OtelTracesExporterFlag)
	logDebug, _ := cmd.Flags().GetBool(LogDebugFlag)
	logger := logging.NewDefaultLogger(file, logDebug, json, otelTraces != "")
	cmd.SetContext(logging.ContextWithLogger(cmd.Context(), logger))
	logging.FromContext(cmd.Context()).Debugf("PreRunE %s", internal.ServiceName)

	options := fx.Options(
		fxOptsFromContext(cmd.Context()),
		fx.Supply(internal.AppInfo{
			Name:                internal.ServiceName,
			Version:             internal.Version,
			TerraformRepository: internal.TerraformRepository,
			BuildDate:           internal.BuildDate,
			Commit:              internal.Commit,
		}),
		otlp.FXModuleFromFlags(cmd, otlp.WithServiceVersion(internal.Version)),
		otlptraces.FXModuleFromFlags(cmd),
	)

	if !service.IsDebug(cmd) {
		options = fx.Options(
			options,
			fx.NopLogger,
		)
	}

	cmd.SetContext(contextWithFxOpts(cmd.Context(), options))
	return nil
}

func (app *App) postRunE(cmd *cobra.Command, args []string) error {
	logging.FromContext(cmd.Context()).Debugf("PostRunE %s", internal.ServiceName)
	return service.NewWithLogger(
		logging.FromContext(cmd.Context()),
		fxOptsFromContext(cmd.Context()),
	).Run(cmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	cmd.SetContext(contextWithFxOpts(cmd.Context(), fx.Options(
		fxOptsFromContext(cmd.Context()),
		server.NewModule(cmd.Context(), cmd.Flags()),
	)))
	return nil
}
