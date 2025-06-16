package pkg

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	formance "github.com/formancehq/formance-sdk-go/v3"
	"go.opentelemetry.io/otel"
)

var (
	Tracer = otel.Tracer("github.com/formancehq/terraform-provider-stack")
)

func NewStackClient(url string, version string, transport http.RoundTripper, tp TokenProviderImpl) (*formance.Formance, error) {
	return formance.New(
		formance.WithServerURL(url),
		formance.WithClient(
			&http.Client{
				Transport: newStackHTTPTransport(
					tp,
					transport,
					map[string][]string{
						"User-Agent": {"terraform-provider-stack/" + version},
					},
				),
			},
		),
	), nil
}

type stackHttpTransport struct {
	tokenProvider       TokenProviderImpl
	underlyingTransport http.RoundTripper
	defaultHeaders      map[string][]string
}

func generateTraceParent() (string, string, error) {
	traceID := make([]byte, 16)
	spanID := make([]byte, 8)
	if _, err := rand.Read(traceID); err != nil {
		return "", "", err
	}
	if _, err := rand.Read(spanID); err != nil {
		return "", "", err
	}

	traceIDHex := hex.EncodeToString(traceID)
	spanIDHex := hex.EncodeToString(spanID)
	traceparent := fmt.Sprintf("00-%s-%s-01", traceIDHex, spanIDHex) // 01 = sampled

	return traceparent, traceIDHex, nil
}

func injectTraceParent(request *http.Request) error {
	traceparent, _, err := generateTraceParent()
	if err != nil {
		return fmt.Errorf("failed to generate traceparent: %w", err)
	}

	request.Header.Set("traceparent", traceparent)
	return nil
}

func (s *stackHttpTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	injectTraceParent(request)
	token, err := s.tokenProvider.StackSecurityToken(request.Context())
	if err != nil {
		return nil, fmt.Errorf("failed to get stack security token: %w", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	for key, values := range s.defaultHeaders {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	return s.underlyingTransport.RoundTrip(request.WithContext(request.Context()))
}

func newStackHTTPTransport(tp TokenProviderImpl, transport http.RoundTripper, defaultHeaders map[string][]string) *stackHttpTransport {
	return &stackHttpTransport{
		underlyingTransport: transport,
		defaultHeaders:      defaultHeaders,
		tokenProvider:       tp,
	}
}

var _ http.RoundTripper = &stackHttpTransport{}
