package pkg

import (
	"fmt"
	"net/http"
)

type stackHttpTransport struct {
	tokenProvider       TokenProviderImpl
	defaultHeaders      map[string][]string
	underlyingTransport http.RoundTripper
}

func (s *stackHttpTransport) RoundTrip(request *http.Request) (*http.Response, error) {
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
