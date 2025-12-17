package sdk

import (
	"context"

	formance "github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
)

//go:generate mockgen -typed -destination=stack_generated.go -package=sdk . StackSdkImpl
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

type StackSdkFactory func(opts ...formance.SDKOption) StackSdkImpl

func NewStackSdk() StackSdkFactory {
	return func(opts ...formance.SDKOption) StackSdkImpl {
		c := formance.New(opts...)
		return &defaultStackSdk{
			Formance:              c,
			LedgerSdkImpl:         newLedgerSdk(c.Ledger),
			PaymentsSdkImpl:       newPaymentsSdk(c.Payments),
			WebhooksSdkImpl:       newWebhooksSdk(c.Webhooks),
			ReconciliationSdkImpl: newReconciliationSdk(c.Reconciliation),
		}
	}
}
