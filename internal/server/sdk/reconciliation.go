package sdk

import (
	"context"

	formance "github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
)

//go:generate mockgen -destination=reconciliation_generated.go -package=sdk . ReconciliationSdkImpl
type ReconciliationSdkImpl interface {
	CreatePolicy(ctx context.Context, request shared.PolicyRequest, opts ...operations.Option) (*operations.CreatePolicyResponse, error)
	GetPolicy(ctx context.Context, request operations.GetPolicyRequest, opts ...operations.Option) (*operations.GetPolicyResponse, error)
	DeletePolicy(ctx context.Context, request operations.DeletePolicyRequest, opts ...operations.Option) (*operations.DeletePolicyResponse, error)
}

var _ ReconciliationSdkImpl = &defaultReconciliationSdk{}

type defaultReconciliationSdk struct {
	*formance.Reconciliation
}

func (s *defaultReconciliationSdk) CreatePolicy(ctx context.Context, request shared.PolicyRequest, opts ...operations.Option) (*operations.CreatePolicyResponse, error) {
	return s.V1.CreatePolicy(ctx, request, opts...)
}

func (s *defaultReconciliationSdk) GetPolicy(ctx context.Context, request operations.GetPolicyRequest, opts ...operations.Option) (*operations.GetPolicyResponse, error) {
	return s.V1.GetPolicy(ctx, request, opts...)
}
func (s *defaultReconciliationSdk) DeletePolicy(ctx context.Context, request operations.DeletePolicyRequest, opts ...operations.Option) (*operations.DeletePolicyResponse, error) {
	return s.V1.DeletePolicy(ctx, request, opts...)
}

func newReconciliationSdk(reconciliation *formance.Reconciliation) ReconciliationSdkImpl {
	return &defaultReconciliationSdk{
		Reconciliation: reconciliation,
	}
}
