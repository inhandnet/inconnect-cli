package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultClientID     = "17953450251798098136"
	DefaultClientSecret = "08E9EC6793345759456CB8BAE52615F3"
)

type OAuthClient struct {
	ClientID     string
	ClientSecret string
}

func FetchOAuthClient(ctx context.Context, host string) (*OAuthClient, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, host+"/api/platform/config", http.NoBody)
	if err != nil {
		return defaultOAuthClient(), nil
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return defaultOAuthClient(), nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return defaultOAuthClient(), nil
	}
	if resp.StatusCode != http.StatusOK {
		return defaultOAuthClient(), nil
	}

	var config struct {
		Result struct {
			Auth struct {
				ClientID     json.Number `json:"clientId"`
				ClientSecret string      `json:"clientSecret"`
			} `json:"auth"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &config); err != nil {
		return defaultOAuthClient(), nil
	}

	clientID := config.Result.Auth.ClientID.String()
	if clientID == "" {
		return defaultOAuthClient(), nil
	}

	return &OAuthClient{
		ClientID:     clientID,
		ClientSecret: config.Result.Auth.ClientSecret,
	}, nil
}

func defaultOAuthClient() *OAuthClient {
	return &OAuthClient{ClientID: DefaultClientID, ClientSecret: DefaultClientSecret}
}

type OAuthToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"-"`
}

func RefreshAccessToken(host, clientID, clientSecret, refreshToken string) (*OAuthToken, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	resp, err := http.Post(host+"/oauth2/access_token", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading refresh response: %w", err)
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		Error        string `json:"error"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("parsing refresh response: %w (%s)", err, string(body))
	}
	if tokenResp.Error != "" {
		return nil, fmt.Errorf("refresh failed: %s", tokenResp.Error)
	}
	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("refresh returned empty access_token")
	}

	token := &OAuthToken{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    tokenResp.TokenType,
		ExpiresIn:    tokenResp.ExpiresIn,
	}
	if token.ExpiresIn > 0 {
		token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	}
	return token, nil
}

func ExchangeCodeForToken(host, code, clientID, clientSecret, redirectURI string) (*OAuthToken, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"redirect_uri":  {redirectURI},
	}

	resp, err := http.PostForm(host+"/oauth2/access_token", data)
	if err != nil {
		return nil, fmt.Errorf("token exchange request failed: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		Error        string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}
	if tokenResp.Error != "" {
		return nil, fmt.Errorf("token exchange failed: %s", tokenResp.Error)
	}
	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("token exchange returned empty access_token")
	}

	token := &OAuthToken{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    tokenResp.TokenType,
		ExpiresIn:    tokenResp.ExpiresIn,
	}
	if token.ExpiresIn > 0 {
		token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	}

	return token, nil
}
