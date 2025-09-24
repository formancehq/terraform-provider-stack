package sdk

import (
	"context"

	formance "github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
)

//go:generate mockgen -typed -destination=payments_generated.go -package=sdk . PaymentsSdkImpl
type PaymentsSdkImpl interface {
	CreatePool(ctx context.Context, request *shared.V3CreatePoolRequest) (*operations.V3CreatePoolResponse, error)
	GetPool(ctx context.Context, request operations.V3GetPoolRequest) (*operations.V3GetPoolResponse, error)
	DeletePool(ctx context.Context, request operations.V3DeletePoolRequest) (*operations.V3DeletePoolResponse, error)

	AddAccountToPool(ctx context.Context, request operations.V3AddAccountToPoolRequest) (*operations.V3AddAccountToPoolResponse, error)
	RemoveAccountFromPool(ctx context.Context, request operations.V3RemoveAccountFromPoolRequest) (*operations.V3RemoveAccountFromPoolResponse, error)

	CreateConnector(ctx context.Context, request operations.V3InstallConnectorRequest) (*operations.V3InstallConnectorResponse, error)
	GetConnector(ctx context.Context, request operations.V3GetConnectorConfigRequest) (*operations.V3GetConnectorConfigResponse, error)
	DeleteConnector(ctx context.Context, request operations.V3UninstallConnectorRequest) (*operations.V3UninstallConnectorResponse, error)
	UpdateConnector(ctx context.Context, request operations.V3UpdateConnectorConfigRequest) (*operations.V3UpdateConnectorConfigResponse, error)
}

var _ PaymentsSdkImpl = &defaultPaymentsSdk{}

type defaultPaymentsSdk struct {
	*formance.Payments
}

func (s *defaultPaymentsSdk) CreatePool(ctx context.Context, request *shared.V3CreatePoolRequest) (*operations.V3CreatePoolResponse, error) {
	return s.V3.CreatePool(ctx, request)
}

func (s *defaultPaymentsSdk) GetPool(ctx context.Context, request operations.V3GetPoolRequest) (*operations.V3GetPoolResponse, error) {
	return s.V3.GetPool(ctx, request)
}

func (s *defaultPaymentsSdk) DeletePool(ctx context.Context, request operations.V3DeletePoolRequest) (*operations.V3DeletePoolResponse, error) {
	return s.V3.DeletePool(ctx, request)
}

func (s *defaultPaymentsSdk) AddAccountToPool(ctx context.Context, request operations.V3AddAccountToPoolRequest) (*operations.V3AddAccountToPoolResponse, error) {
	return s.V3.AddAccountToPool(ctx, request)
}
func (s *defaultPaymentsSdk) RemoveAccountFromPool(ctx context.Context, request operations.V3RemoveAccountFromPoolRequest) (*operations.V3RemoveAccountFromPoolResponse, error) {
	return s.V3.RemoveAccountFromPool(ctx, request)
}

func (s *defaultPaymentsSdk) CreateConnector(ctx context.Context, request operations.V3InstallConnectorRequest) (*operations.V3InstallConnectorResponse, error) {
	return s.V3.InstallConnector(ctx, request)
}

func (s *defaultPaymentsSdk) GetConnector(ctx context.Context, request operations.V3GetConnectorConfigRequest) (*operations.V3GetConnectorConfigResponse, error) {
	return s.V3.GetConnectorConfig(ctx, request)
}

func (s *defaultPaymentsSdk) DeleteConnector(ctx context.Context, request operations.V3UninstallConnectorRequest) (*operations.V3UninstallConnectorResponse, error) {
	return s.V3.UninstallConnector(ctx, request)
}

func (s *defaultPaymentsSdk) UpdateConnector(ctx context.Context, request operations.V3UpdateConnectorConfigRequest) (*operations.V3UpdateConnectorConfigResponse, error) {
	return s.V3.V3UpdateConnectorConfig(ctx, request)
}

func newPaymentsSdk(payments *formance.Payments) PaymentsSdkImpl {
	return &defaultPaymentsSdk{
		Payments: payments,
	}
}
