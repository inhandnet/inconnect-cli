package api

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/inhandnet/ics-cli/internal/build"
	"github.com/inhandnet/ics-cli/internal/debug"
)

var userAgent = fmt.Sprintf("ics-cli/%s (%s/%s)", build.Version, runtime.GOOS, runtime.GOARCH)

type TokenTransport struct {
	Token        string
	RefreshToken string
	Host         string
	ClientID     string
	ClientSecret string
	OnRefresh    func(accessToken, refreshToken string, expiry time.Time)
	Base         http.RoundTripper
}

func (t *TokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", userAgent)

	if t.Token != "" && t.isSameHost(req) {
		req.Header.Set("Authorization", "Bearer "+t.Token)
	}

	resp, err := t.doRoundTrip(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode == 401 && t.RefreshToken != "" {
		resp.Body.Close()
		debug.Log("token expired, refreshing...")

		newToken, refreshErr := RefreshAccessToken(t.Host, t.ClientID, t.ClientSecret, t.RefreshToken)
		if refreshErr != nil {
			debug.Log("token refresh failed: %v", refreshErr)
			return resp, nil
		}

		t.Token = newToken.AccessToken
		if newToken.RefreshToken != "" {
			t.RefreshToken = newToken.RefreshToken
		}
		if t.OnRefresh != nil {
			t.OnRefresh(newToken.AccessToken, newToken.RefreshToken, newToken.ExpiresAt)
		}
		debug.Log("token refreshed successfully")

		req.Header.Set("Authorization", "Bearer "+t.Token)
		resp, err = t.doRoundTrip(req)
	}
	return resp, err
}

func (t *TokenTransport) doRoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	t.debugRequest(req)
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	resp, err := base.RoundTrip(req)
	if err != nil {
		debug.Log("< request error: %v", err)
		return resp, err
	}
	t.debugResponse(resp, time.Since(start))
	return resp, nil
}

func (t *TokenTransport) debugRequest(req *http.Request) {
	if !debug.Enabled {
		return
	}
	debug.Log("> %s %s", req.Method, req.URL.String())
	for _, h := range []string{"User-Agent", "Content-Type", "Authorization"} {
		if v := req.Header.Get(h); v != "" {
			if h == "Authorization" {
				v = "****"
			}
			debug.Log("> %s: %s", h, v)
		}
	}
}

func (t *TokenTransport) debugResponse(resp *http.Response, elapsed time.Duration) {
	if !debug.Enabled {
		return
	}
	debug.Log("< %d %s (%s)", resp.StatusCode, http.StatusText(resp.StatusCode), elapsed.Round(time.Millisecond))
}

func (t *TokenTransport) isSameHost(req *http.Request) bool {
	if t.Host == "" {
		return true
	}
	parsed, err := url.Parse(t.Host)
	if err != nil {
		return true
	}
	return req.URL.Host == parsed.Host
}
