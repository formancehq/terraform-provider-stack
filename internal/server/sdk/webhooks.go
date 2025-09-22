package sdk

import (
	"context"

	formance "github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
)

//go:generate mockgen -typed -destination=webhooks_generated.go -package=sdk . WebhooksSdkImpl
type WebhooksSdkImpl interface {
	InsertConfig(ctx context.Context, request shared.ConfigUser) (*operations.InsertConfigResponse, error)
	GetManyConfigs(ctx context.Context, request operations.GetManyConfigsRequest) (*operations.GetManyConfigsResponse, error)
	DeleteConfig(ctx context.Context, request operations.DeleteConfigRequest, opts ...operations.Option) (*operations.DeleteConfigResponse, error)
	UpdateConfig(ctx context.Context, request operations.UpdateConfigRequest, opts ...operations.Option) (*operations.UpdateConfigResponse, error)
}

var _ WebhooksSdkImpl = &defaultWebhooksSdk{}

type defaultWebhooksSdk struct {
	*formance.Webhooks
}

func (s *defaultWebhooksSdk) InsertConfig(ctx context.Context, request shared.ConfigUser) (*operations.InsertConfigResponse, error) {
	return s.V1.InsertConfig(ctx, request)
}
func (s *defaultWebhooksSdk) GetManyConfigs(ctx context.Context, request operations.GetManyConfigsRequest) (*operations.GetManyConfigsResponse, error) {
	return s.V1.GetManyConfigs(ctx, request)
}
func (s *defaultWebhooksSdk) DeleteConfig(ctx context.Context, request operations.DeleteConfigRequest, opts ...operations.Option) (*operations.DeleteConfigResponse, error) {
	return s.V1.DeleteConfig(ctx, request, opts...)
}
func (s *defaultWebhooksSdk) UpdateConfig(ctx context.Context, request operations.UpdateConfigRequest, opts ...operations.Option) (*operations.UpdateConfigResponse, error) {
	return s.V1.UpdateConfig(ctx, request, opts...)
}

func newWebhooksSdk(webhooks *formance.Webhooks) WebhooksSdkImpl {
	return &defaultWebhooksSdk{
		Webhooks: webhooks,
	}
}
