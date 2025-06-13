package sdk

import (
	"net/http"

	formance "github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/terraform-provider-stack/pkg"
)

//go:generate mockgen -destination=stack_generated.go -package=sdk . StackSdkImpl
type StackSdkImpl interface{}

var _ StackSdkImpl = &defaultStackSdk{}

type defaultStackSdk struct {
	*formance.Formance
}

type StackSdkFactory func(url string, version string, transport http.RoundTripper, tp pkg.TokenProviderImpl) (StackSdkImpl, error)

func NewStackSdk() StackSdkFactory {
	return func(url, version string, transport http.RoundTripper, tp pkg.TokenProviderImpl) (StackSdkImpl, error) {
		client, err := pkg.NewStackClient(url, version, transport, tp)
		if err != nil {
			return nil, err
		}

		return &defaultStackSdk{
			Formance: client,
		}, nil
	}
}
