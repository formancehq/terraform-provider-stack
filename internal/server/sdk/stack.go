package sdk

import (
	"context"
	"net/http"

	formance "github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/terraform-provider-stack/pkg"
)

//go:generate mockgen -destination=stack_generated.go -package=sdk . StackSdkImpl
type StackSdkImpl interface {
	GetVersions(ctx context.Context) (*operations.GetVersionsResponse, error)
	Ledger() LedgerSdkImpl
	Payments() PaymentsSdkImpl
	Webhooks() WebhooksSdkImpl
	Reconciliation() ReconciliationSdkImpl
}

var _ StackSdkImpl = &defaultStackSdk{}

type defaultStackSdk struct {
	*formance.Formance
	LedgerSdkImpl
	PaymentsSdkImpl
	WebhooksSdkImpl
	ReconciliationSdkImpl
}

func (s *defaultStackSdk) GetVersions(ctx context.Context) (*operations.GetVersionsResponse, error) {
	return s.Formance.GetVersions(ctx)
}

func (s *defaultStackSdk) Ledger() LedgerSdkImpl {
	return s.LedgerSdkImpl
}
func (s *defaultStackSdk) Payments() PaymentsSdkImpl {
	return s.PaymentsSdkImpl
}
func (s *defaultStackSdk) Webhooks() WebhooksSdkImpl {
	return s.WebhooksSdkImpl
}
func (s *defaultStackSdk) Reconciliation() ReconciliationSdkImpl {
	return s.ReconciliationSdkImpl
}

type StackSdkFactory func(url string, version string, transport http.RoundTripper, tp pkg.TokenProviderImpl) (StackSdkImpl, error)

func NewStackSdk() StackSdkFactory {
	return func(url, version string, transport http.RoundTripper, tp pkg.TokenProviderImpl) (StackSdkImpl, error) {
		client, err := pkg.NewStackClient(url, version, transport, tp)
		if err != nil {
			return nil, err
		}

		return &defaultStackSdk{
			Formance:              client,
			LedgerSdkImpl:         newLedgerSdk(client.Ledger),
			PaymentsSdkImpl:       newPaymentsSdk(client.Payments),
			WebhooksSdkImpl:       newWebhooksSdk(client.Webhooks),
			ReconciliationSdkImpl: newReconciliationSdk(client.Reconciliation),
		}, nil
	}
}
