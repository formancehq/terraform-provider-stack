package speakeasyretry

import (
	"github.com/formancehq/formance-sdk-go/v3/pkg/retry"
	"github.com/spf13/pflag"
	"go.uber.org/fx"
)

var (
	RetryFlag                = "retry-enabled"
	RetryInitialIntervalFlag = "retry-initial-interval"
	RetryMaxIntervalFlag     = "retry-max-interval"
	RetryMaxElapsedTimeFlag  = "retry-max-elapsed-time"
	RetryExponentFlag        = "retry-exponent"
)

const (
	DefaultRetryInitialInterval = 1000
	DefaultRetryMaxInterval     = 3000
	DefaultRetryMaxElapsedTime  = 10000
	DefaultRetryExponent        = 2.0
)

func AddFlags(flags *pflag.FlagSet) {
	if flags.Lookup(RetryFlag) == nil {
		flags.Bool(RetryFlag, true, "Enable SDK retry")
	}
	if flags.Lookup(RetryInitialIntervalFlag) == nil {
		flags.Int(RetryInitialIntervalFlag, DefaultRetryInitialInterval, "Initial interval for retry backoff strategy in milliseconds")
	}
	if flags.Lookup(RetryMaxIntervalFlag) == nil {
		flags.Int(RetryMaxIntervalFlag, DefaultRetryMaxInterval, "Max interval for retry backoff strategy in milliseconds")
	}
	if flags.Lookup(RetryMaxElapsedTimeFlag) == nil {
		flags.Int(RetryMaxElapsedTimeFlag, DefaultRetryMaxElapsedTime, "Max elapsed time for retry backoff strategy in milliseconds")
	}
	if flags.Lookup(RetryExponentFlag) == nil {
		flags.Float64(RetryExponentFlag, DefaultRetryExponent, "Exponent for retry backoff strategy")
	}
}
func NewModule(flags *pflag.FlagSet) fx.Option {
	initialInterval, _ := flags.GetInt(RetryInitialIntervalFlag)
	maxInterval, _ := flags.GetInt(RetryMaxIntervalFlag)
	maxElapsedTime, _ := flags.GetInt(RetryMaxElapsedTimeFlag)
	exponent, _ := flags.GetFloat64(RetryExponentFlag)

	if enabled, _ := flags.GetBool(RetryFlag); !enabled {
		return fx.Options(
			fx.Provide(func() *retry.Config {
				return nil
			}),
		)
	}
	return fx.Options(
		fx.Supply(&retry.Config{
			Strategy: "backoff",
			Backoff: &retry.BackoffStrategy{
				InitialInterval: initialInterval,
				MaxInterval:     maxInterval,
				Exponent:        exponent,
				MaxElapsedTime:  maxElapsedTime,
			},
			RetryConnectionErrors: true,
		}),
	)
}
