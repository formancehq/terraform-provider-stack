package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	cloudpkg "github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
)

//go:generate mockgen -typed -destination=tokenprovider_generated.go -package=pkg . TokenProviderImpl
type TokenProviderImpl interface {
	StackSecurityToken(ctx context.Context) (*cloudpkg.TokenInfo, error)
}

type TokenProviderFactory func(transport http.RoundTripper, creds cloudpkg.Creds, tokenProvider cloudpkg.TokenProviderImpl, stack Stack) TokenProviderImpl

type TokenProvider struct {
	client        *http.Client
	tokenProvider cloudpkg.TokenProviderImpl
	creds         cloudpkg.Creds

	stack Stack
	token *cloudpkg.TokenInfo
}

func NewTokenProviderFn() TokenProviderFactory {
	return func(transport http.RoundTripper, creds cloudpkg.Creds, tokenProvider cloudpkg.TokenProviderImpl, stack Stack) TokenProviderImpl {
		return NewTokenProvider(transport, creds, tokenProvider, stack)
	}
}

func NewTokenProvider(
	transport http.RoundTripper,
	creds cloudpkg.Creds,
	tokenProvider cloudpkg.TokenProviderImpl,
	stack Stack,
) TokenProvider {
	return TokenProvider{
		client: &http.Client{
			Transport: transport,
		},
		creds:         creds,
		tokenProvider: tokenProvider,
		stack:         stack,
		token:         &cloudpkg.TokenInfo{},
	}
}

type Stack struct {
	Id             string
	OrganizationId string
	Uri            string
}

func (p TokenProvider) StackSecurityToken(ctx context.Context) (*cloudpkg.TokenInfo, error) {
	token, err := p.tokenProvider.RefreshToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to refresh token: %w", err)
	}

	form := url.Values{
		"grant_type":         []string{string(oidc.GrantTypeTokenExchange)},
		"audience":           []string{fmt.Sprintf("stack://%s/%s", p.stack.OrganizationId, p.stack.Id)},
		"subject_token":      []string{token.AccessToken},
		"subject_token_type": []string{"urn:ietf:params:oauth:token-type:access_token"},
	}

	membershipDiscoveryConfiguration, err := client.Discover(ctx, p.creds.Endpoint(), p.client)
	if err != nil {
		return nil, fmt.Errorf("unable to discover membership configuration: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, membershipDiscoveryConfiguration.TokenEndpoint,
		bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(p.creds.ClientId(), p.creds.ClientSecret())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ret, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to exchange token: %w", err)
	}

	defer func() {
		if err := ret.Body.Close(); err != nil {
			panic(err)
		}
	}()
	if ret.StatusCode != http.StatusOK {
		data, err := io.ReadAll(ret.Body)
		if err != nil {
			panic(err)
		}
		return nil, errors.New(string(data))
	}

	securityToken := oauth2.Token{}
	if err := json.NewDecoder(ret.Body).Decode(&securityToken); err != nil {
		return nil, err
	}

	baseUri, err := url.Parse(p.stack.Uri)
	if err != nil {
		return nil, fmt.Errorf("invalid stack Uri %s: %w", p.stack.Uri, err)
	}
	apiUrl, err := url.JoinPath(baseUri.String(), "api", "auth")
	if err != nil {
		return nil, fmt.Errorf("invalid stack Uri %s: %w", p.stack.Uri, err)
	}

	form = url.Values{
		"grant_type": []string{"urn:ietf:params:oauth:grant-type:jwt-bearer"},
		"assertion":  []string{securityToken.AccessToken},
		"scope":      []string{"openid email"},
	}

	stackDiscoveryConfiguration, err := client.Discover(ctx, apiUrl, p.client)
	if err != nil {
		return nil, fmt.Errorf("unable to discover stack configuration: %w", err)
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, stackDiscoveryConfiguration.TokenEndpoint,
		bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("unable to create request for token exchange: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ret, err = p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to exchange security token: %w", err)
	}

	defer func() {
		if err := ret.Body.Close(); err != nil {
			panic(err)
		}
	}()
	if ret.StatusCode != http.StatusOK {
		data, err := io.ReadAll(ret.Body)
		if err != nil {
			panic(err)
		}
		return nil, errors.New(string(data))
	}

	stackToken := oauth2.Token{}
	if err := json.NewDecoder(ret.Body).Decode(&stackToken); err != nil {
		return nil, err
	}

	p.token.AccessToken = stackToken.AccessToken
	p.token.RefreshToken = stackToken.RefreshToken
	p.token.Expiry = stackToken.Expiry

	return p.token, nil
}
