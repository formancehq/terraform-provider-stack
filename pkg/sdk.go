package pkg

import (
	"fmt"
	"net/http"

	formance "github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/go-libs/v3/httpclient"
	"github.com/formancehq/go-libs/v3/otlp"
	"github.com/formancehq/terraform-provider-cloud/pkg"
)

func NewStackClient(url string, version string, ti *pkg.TokenInfo) (*formance.Formance, error) {
	return formance.New(
		formance.WithServerURL(url),
		formance.WithClient(
			&http.Client{
				Transport: newStackHTTPTransport(
					ti,
					map[string][]string{
						"User-Agent": {"terraform-provider-stack/" + version},
					},
				),
			},
		),
	), nil
}

type stackHttpTransport struct {
	tokenInfo           *pkg.TokenInfo
	defaultHeaders      map[string][]string
	underlyingTransport http.RoundTripper
}

func (s *stackHttpTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.tokenInfo.AccessToken))
	for key, values := range s.defaultHeaders {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
	return s.underlyingTransport.RoundTrip(request)
}

func newStackHTTPTransport(i *pkg.TokenInfo, defaultHeaders map[string][]string) *stackHttpTransport {
	return &stackHttpTransport{
		underlyingTransport: otlp.NewRoundTripper(httpclient.NewDebugHTTPTransport(http.DefaultTransport), true),
		defaultHeaders:      defaultHeaders,
		tokenInfo:           i,
	}
}

var _ http.RoundTripper = &stackHttpTransport{}
