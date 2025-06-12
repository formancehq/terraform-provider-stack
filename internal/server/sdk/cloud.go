package sdk

import (
	"context"
	"net/http"

	"github.com/formancehq/terraform-provider-cloud/pkg"
	cloud "github.com/formancehq/terraform-provider-cloud/sdk"
)

//go:generate mockgen -destination=cloud_generated.go -package=sdk . CloudSDK
type CloudSDK interface {
	GetStack(ctx context.Context, organizationID, stackID string) (*cloud.CreateStackResponse, *http.Response, error)
	ListModules(ctx context.Context, organizationID, stackID string) (*cloud.ListModulesResponse, *http.Response, error)
}

var _ CloudSDK = &sdkImpl{}

type sdkImpl struct {
	sdk cloud.DefaultAPI
}

// GetStack implements CloudSDK.
func (s *sdkImpl) GetStack(ctx context.Context, organizationID string, stackID string) (*cloud.CreateStackResponse, *http.Response, error) {
	return s.sdk.GetStack(ctx, organizationID, stackID).Execute()
}

// ListModules implements CloudSDK.
func (s *sdkImpl) ListModules(ctx context.Context, organizationID string, stackID string) (*cloud.ListModulesResponse, *http.Response, error) {
	return s.sdk.ListModules(ctx, organizationID, stackID).Execute()
}

type CloudFactory func(creds pkg.Creds, transport http.RoundTripper) CloudSDK

// NewCloudSDK returns a factory function that creates a CloudSDK instance using the provided credentials and HTTP transport.
func NewCloudSDK() CloudFactory {
	return func(creds pkg.Creds, transport http.RoundTripper) CloudSDK {
		return &sdkImpl{sdk: pkg.NewSDK(creds, transport)}
	}
}
