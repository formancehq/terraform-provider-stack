package sdk

import (
	"context"

	formance "github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
)

//go:generate mockgen -destination=ledger_generated.go -package=sdk . LedgerSdkImpl
type LedgerSdkImpl interface {
	CreateLedger(ctx context.Context, request operations.V2CreateLedgerRequest) (*operations.V2CreateLedgerResponse, error)
	GetLedger(ctx context.Context, request operations.V2GetLedgerRequest) (*operations.V2GetLedgerResponse, error)
	DeleteLedger(ctx context.Context, name string) error
	UpdateLedgerMetadata(ctx context.Context, request operations.V2UpdateLedgerMetadataRequest, opts ...operations.Option) (*operations.V2UpdateLedgerMetadataResponse, error)
}

var _ LedgerSdkImpl = &defaultLedger{}

type defaultLedger struct {
	*formance.Ledger
}

func (s *defaultLedger) CreateLedger(ctx context.Context, request operations.V2CreateLedgerRequest) (*operations.V2CreateLedgerResponse, error) {
	return s.V2.CreateLedger(ctx, request)
}

func (s *defaultLedger) GetLedger(ctx context.Context, request operations.V2GetLedgerRequest) (*operations.V2GetLedgerResponse, error) {
	return s.V2.GetLedger(ctx, request)
}

func (s *defaultLedger) DeleteLedger(ctx context.Context, name string) error {
	return nil
}

func (s *defaultLedger) UpdateLedgerMetadata(ctx context.Context, request operations.V2UpdateLedgerMetadataRequest, opts ...operations.Option) (*operations.V2UpdateLedgerMetadataResponse, error) {
	return s.V2.UpdateLedgerMetadata(ctx, request, opts...)
}

func newLedgerSdk(ledger *formance.Ledger) LedgerSdkImpl {
	return &defaultLedger{
		Ledger: ledger,
	}
}
