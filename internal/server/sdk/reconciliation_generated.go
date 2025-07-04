// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/formancehq/terraform-provider-stack/internal/server/sdk (interfaces: ReconciliationSdkImpl)
//
// Generated by this command:
//
//	mockgen -destination=reconciliation_generated.go -package=sdk . ReconciliationSdkImpl
//

// Package sdk is a generated GoMock package.
package sdk

import (
	context "context"
	reflect "reflect"

	operations "github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	shared "github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	gomock "go.uber.org/mock/gomock"
)

// MockReconciliationSdkImpl is a mock of ReconciliationSdkImpl interface.
type MockReconciliationSdkImpl struct {
	ctrl     *gomock.Controller
	recorder *MockReconciliationSdkImplMockRecorder
	isgomock struct{}
}

// MockReconciliationSdkImplMockRecorder is the mock recorder for MockReconciliationSdkImpl.
type MockReconciliationSdkImplMockRecorder struct {
	mock *MockReconciliationSdkImpl
}

// NewMockReconciliationSdkImpl creates a new mock instance.
func NewMockReconciliationSdkImpl(ctrl *gomock.Controller) *MockReconciliationSdkImpl {
	mock := &MockReconciliationSdkImpl{ctrl: ctrl}
	mock.recorder = &MockReconciliationSdkImplMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReconciliationSdkImpl) EXPECT() *MockReconciliationSdkImplMockRecorder {
	return m.recorder
}

// CreatePolicy mocks base method.
func (m *MockReconciliationSdkImpl) CreatePolicy(ctx context.Context, request shared.PolicyRequest, opts ...operations.Option) (*operations.CreatePolicyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, request}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreatePolicy", varargs...)
	ret0, _ := ret[0].(*operations.CreatePolicyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePolicy indicates an expected call of CreatePolicy.
func (mr *MockReconciliationSdkImplMockRecorder) CreatePolicy(ctx, request any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, request}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePolicy", reflect.TypeOf((*MockReconciliationSdkImpl)(nil).CreatePolicy), varargs...)
}

// DeletePolicy mocks base method.
func (m *MockReconciliationSdkImpl) DeletePolicy(ctx context.Context, request operations.DeletePolicyRequest, opts ...operations.Option) (*operations.DeletePolicyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, request}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeletePolicy", varargs...)
	ret0, _ := ret[0].(*operations.DeletePolicyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeletePolicy indicates an expected call of DeletePolicy.
func (mr *MockReconciliationSdkImplMockRecorder) DeletePolicy(ctx, request any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, request}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePolicy", reflect.TypeOf((*MockReconciliationSdkImpl)(nil).DeletePolicy), varargs...)
}

// GetPolicy mocks base method.
func (m *MockReconciliationSdkImpl) GetPolicy(ctx context.Context, request operations.GetPolicyRequest, opts ...operations.Option) (*operations.GetPolicyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, request}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetPolicy", varargs...)
	ret0, _ := ret[0].(*operations.GetPolicyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPolicy indicates an expected call of GetPolicy.
func (mr *MockReconciliationSdkImplMockRecorder) GetPolicy(ctx, request any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, request}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPolicy", reflect.TypeOf((*MockReconciliationSdkImpl)(nil).GetPolicy), varargs...)
}
