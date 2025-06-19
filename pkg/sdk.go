package pkg

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	formance "github.com/formancehq/formance-sdk-go/v3"
)

func NewStackClient(opts ...formance.SDKOption) (*formance.Formance, error) {
	return formance.New(opts...), nil
}

type stackHttpTransport struct {
	tokenProvider       TokenProviderImpl
	defaultHeaders      map[string][]string
	underlyingTransport http.RoundTripper
}

type responseContextKey struct{}

func ResponseFromContext(ctx context.Context) *http.Response {
	v := ctx.Value(responseContextKey{})
	if v == nil {
		return nil
	}

	resp, ok := v.(*http.Response)
	if !ok {
		return nil
	}

	return resp
}

func generateTraceID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Sprintf("failed to generate trace ID: %v", err))
	}
	return hex.EncodeToString(b)
}

func generateSpanID() string {
	b := make([]byte, 8) 
	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Sprintf("failed to generate span ID: %v", err))
	}
	return hex.EncodeToString(b)
}

func injectTraceparentHeader(request *http.Request) {
	traceID := generateTraceID()
	spanID := generateSpanID()
	traceFlags := "01" // sampled

	traceparent := "00-" + traceID + "-" + spanID + "-" + traceFlags
	request.Header.Set("Traceparent", traceparent)
}

func (s *stackHttpTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	injectTraceparentHeader(request)
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

	resp, err := s.underlyingTransport.RoundTrip(request)
	ctx := context.WithValue(request.Context(), responseContextKey{}, resp)
	_ = request.WithContext(ctx)
	if err != nil {
		return resp, fmt.Errorf("error during round trip: %w", err)
	}

	return resp, nil
}

func NewStackHTTPTransport(tp TokenProviderImpl, transport http.RoundTripper, defaultHeaders map[string][]string) *stackHttpTransport {
	return &stackHttpTransport{
		underlyingTransport: transport,
		defaultHeaders:      defaultHeaders,
		tokenProvider:       tp,
	}
}

var _ http.RoundTripper = &stackHttpTransport{}
