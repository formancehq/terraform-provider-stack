package sdk

import (
	"context"
	"net/http"

	"github.com/formancehq/terraform-provider-cloud/pkg"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/operations"
)

//go:generate mockgen -typed -destination=cloud_generated.go -package=sdk . CloudSDK
type CloudSDK interface {
	GetStack(ctx context.Context, organizationID, stackID string) (*operations.GetStackResponse, error)
	ListModules(ctx context.Context, organizationID, stackID string) (*operations.ListModulesResponse, error)
}

var _ CloudSDK = &sdkImpl{}

type sdkImpl struct {
	sdk *membershipclient.FormanceCloud
}

// GetStack implements CloudSDK.
func (s *sdkImpl) GetStack(ctx context.Context, organizationID string, stackID string) (*operations.GetStackResponse, error) {
	return s.sdk.GetStack(ctx, organizationID, stackID)
}

// ListModules implements CloudSDK.
func (s *sdkImpl) ListModules(ctx context.Context, organizationID string, stackID string) (*operations.ListModulesResponse, error) {
	return s.sdk.ListModules(ctx, organizationID, stackID)
}

type CloudFactory func(creds pkg.Creds, transport http.RoundTripper) CloudSDK

func NewCloudSDK(opts ...membershipclient.SDKOption) CloudFactory {
	return func(creds pkg.Creds, transport http.RoundTripper) CloudSDK {
		tp := pkg.NewTokenProvider(transport, creds)
		return &sdkImpl{sdk: pkg.NewSDK(creds.Endpoint(), transport, tp)}
	}
}
