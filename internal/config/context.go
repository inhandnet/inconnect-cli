package config

import (
	"net/url"
	"os"
	"strings"
	"time"
)

type Context struct {
	Host         string    `yaml:"host"`
	Token        string    `yaml:"token,omitempty"`
	RefreshToken string    `yaml:"refresh_token,omitempty"`
	User         string    `yaml:"user,omitempty"`
	Authority    string    `yaml:"authority,omitempty"`
	ExpiresAt    time.Time `yaml:"expires_at,omitempty"`
	ClientID     string    `yaml:"client_id,omitempty"`
	ClientSecret string    `yaml:"client_secret,omitempty"`

	AdminToken        string    `yaml:"admin_token,omitempty"`
	AdminRefreshToken string    `yaml:"admin_refresh_token,omitempty"`
	AdminExpiresAt    time.Time `yaml:"admin_expires_at,omitempty"`
	OrgID             string    `yaml:"org_id,omitempty"`
}

func (c *Context) IsImpersonating() bool {
	return c.AdminToken != ""
}

func (c *Context) EffectiveToken() string {
	if t := os.Getenv("ICS_TOKEN"); t != "" {
		return t
	}
	return c.Token
}

func (c *Context) APIURL() string {
	return "https://" + c.Host
}

func (c *Context) BaseURL() string {
	host := c.normalizedHost()
	if !strings.HasPrefix(host, "http") {
		host = "https://" + host
	}
	return strings.TrimRight(host, "/")
}

func (c *Context) normalizedHost() string {
	h := c.Host
	h = strings.TrimSpace(h)
	h = strings.TrimRight(h, "/")
	if u, err := url.Parse(h); err == nil && u.Host != "" {
		return u.Scheme + "://" + u.Host
	}
	return h
}

func ResolveBaseURL(host string) string {
	ctx := &Context{Host: host}
	return ctx.BaseURL()
}
