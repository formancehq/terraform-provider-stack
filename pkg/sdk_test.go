package pkg

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	pkg "github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestStackHttpTranport(t *testing.T) {
	t.Parallel()
	ctx := logging.TestingContext()

	ctrl := gomock.NewController(t)
	mock := NewMockTokenProviderImpl(ctrl)
	transport := newStackHTTPTransport(mock, http.DefaultTransport, nil)

	mock.EXPECT().StackSecurityToken(gomock.Any()).Return(&pkg.TokenInfo{
		AccessToken: "test-token",
	}, nil).Times(1)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Contains(t, r.Header, "Authorization")
		require.Contains(t, r.Header, "Traceparent")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer s.Close()

	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.URL+"/", nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

}
