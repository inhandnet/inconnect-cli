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
	if t := os.Getenv("INCONNECT_TOKEN"); t != "" {
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

// ngrokServers maps a platform host (as stored in Context.Host, the resolved
// domain — not the cn/us/eu region alias) to its ngrok tunnel server. cn/us/eu
// follow the "ngrok.<host>" pattern; dev/beta deliberately do not, hence the
// explicit table.
var ngrokServers = map[string]string{
	"ics.inhandiot.com":            "ngrok.ics.inhandiot.com:4443",
	"ics.inhandnetworks.com":       "ngrok.ics.inhandnetworks.com:4443",
	"ics.inhandnetworks.eu":        "ngrok.ics.inhandnetworks.eu:4443",
	"dev.inconnect.inhand.design":  "ngrok.10.5.17.73.nip.io:4443",
	"beta.inconnect.inhand.design": "ngrok.10.5.17.74.nip.io:4443",
}

// NgrokServer returns the ngrok tunnel server for this context's platform host.
// Known hosts are looked up in the table; an unrecognized host falls back to
// the "ngrok.<host>:4443" convention. Returns "" only when Host is empty.
func (c *Context) NgrokServer() string {
	host := c.ngrokHostname()
	if host == "" {
		return ""
	}
	if s, ok := ngrokServers[host]; ok {
		return s
	}
	return "ngrok." + host + ":4443"
}

// ngrokHostname extracts the bare hostname from Host, tolerating an optional
// scheme, port, or trailing path.
func (c *Context) ngrokHostname() string {
	h := strings.TrimSpace(c.Host)
	h = strings.TrimRight(h, "/")
	if h == "" {
		return ""
	}
	if !strings.Contains(h, "://") {
		h = "//" + h
	}
	if u, err := url.Parse(h); err == nil && u.Hostname() != "" {
		return u.Hostname()
	}
	return strings.TrimPrefix(h, "//")
}
