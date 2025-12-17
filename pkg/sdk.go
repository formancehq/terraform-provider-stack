package pkg

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
)

type stackHttpTransport struct {
	tokenProvider       TokenProviderImpl
	defaultHeaders      map[string][]string
	underlyingTransport http.RoundTripper
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

	return s.underlyingTransport.RoundTrip(request)
}

func NewStackHTTPTransport(tp TokenProviderImpl, transport http.RoundTripper, defaultHeaders map[string][]string) *stackHttpTransport {
	return &stackHttpTransport{
		underlyingTransport: transport,
		defaultHeaders:      defaultHeaders,
		tokenProvider:       tp,
	}
}

var _ http.RoundTripper = &stackHttpTransport{}
