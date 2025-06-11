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

	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
)

type TokenProvider struct {
	client *http.Client

	creds pkg.Creds

	stack *pkg.TokenInfo
}

func NewTokenProvider(client *http.Client, creds pkg.Creds) TokenProvider {
	return TokenProvider{
		client: client,
		creds:  creds,
		stack:  &pkg.TokenInfo{},
	}
}

type Stack struct {
	Id             string
	OrganizationId string
	Uri            string
}

func (p TokenProvider) StackSecurityToken(ctx context.Context, tokenProvider pkg.TokenProviderImpl, stack Stack) (*pkg.TokenInfo, error) {
	token, err := tokenProvider.RefreshToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to refresh token: %w", err)
	}

	form := url.Values{
		"grant_type":         []string{string(oidc.GrantTypeTokenExchange)},
		"audience":           []string{fmt.Sprintf("stack://%s/%s", stack.OrganizationId, stack.Id)},
		"subject_token":      []string{token.AccessToken},
		"subject_token_type": []string{"urn:ietf:params:oauth:token-type:access_token"},
	}

	membershipDiscoveryConfiguration, err := client.Discover(ctx, p.creds.Endpoint(), p.client)
	if err != nil {
		return nil, err
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
		return nil, err
	}

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

	baseUri, err := url.Parse(stack.Uri)
	if err != nil {
		return nil, fmt.Errorf("invalid stack Uri %s: %w", stack.Uri, err)
	}
	apiUrl, err := url.JoinPath(baseUri.String(), "api", "auth")
	if err != nil {
		return nil, fmt.Errorf("invalid stack Uri %s: %w", stack.Uri, err)
	}

	form = url.Values{
		"grant_type": []string{"urn:ietf:params:oauth:grant-type:jwt-bearer"},
		"assertion":  []string{securityToken.AccessToken},
		"scope":      []string{"openid email"},
	}

	stackDiscoveryConfiguration, err := client.Discover(ctx, apiUrl, p.client)
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, stackDiscoveryConfiguration.TokenEndpoint,
		bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ret, err = p.client.Do(req)
	if err != nil {
		return nil, err
	}

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

	p.stack.AccessToken = stackToken.AccessToken
	p.stack.RefreshToken = stackToken.RefreshToken
	p.stack.Expiry = stackToken.Expiry

	return p.stack, nil
}
