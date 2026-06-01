package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/inhandnet/ics-cli/internal/api"
	"github.com/inhandnet/ics-cli/internal/browser"
	"github.com/inhandnet/ics-cli/internal/config"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

const defaultCallbackPort = 18920

type loginOptions struct {
	ContextName string
	Host        string
	Port        int
	Timeout     time.Duration
}

func newCmdLogin(f *factory.Factory) *cobra.Command {
	opts := &loginOptions{}

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login via browser",
		Example: `  # Login to China region (default)
  ics auth login

  # Login to dev environment
  ics auth login --host dev

  # Login to a custom domain
  ics auth login --host ics.example.com`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBrowserLogin(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.ContextName, "context", "default", "Context name to create/update")
	cmd.Flags().StringVar(&opts.Host, "host", "cn", `Platform region: "cn", "us", "eu", "dev", "beta", or a custom domain`)
	cmd.Flags().IntVar(&opts.Port, "port", defaultCallbackPort, "Local callback server port")
	cmd.Flags().DurationVar(&opts.Timeout, "timeout", 3*time.Minute, "Timeout waiting for browser login")

	return cmd
}

var regionHosts = map[string]string{
	"cn":   "ics.inhandiot.com",
	"us":   "ics.inhandnetworks.com",
	"eu":   "ics.inhandnetworks.eu",
	"dev":  "dev.inconnect.inhand.design",
	"beta": "beta.inconnect.inhand.design",
}

func resolveHost(host string) (string, error) {
	if host == "" {
		return "", fmt.Errorf("host is required")
	}
	if domain, ok := regionHosts[strings.ToLower(host)]; ok {
		return domain, nil
	}
	return host, nil
}

func runBrowserLogin(f *factory.Factory, opts *loginOptions) error {
	out := f.IO.Out

	host, err := resolveHost(opts.Host)
	if err != nil {
		return err
	}

	var apiURL string
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		apiURL = host
	} else {
		apiURL = "https://" + host
	}

	oauthClient, err := api.FetchOAuthClient(context.Background(), apiURL)
	if err != nil {
		return fmt.Errorf("fetching OAuth config from %s: %w", apiURL, err)
	}

	state := fmt.Sprintf("ics-cli-%d", time.Now().UnixNano())
	redirectURI := fmt.Sprintf("http://localhost:%d/callback", opts.Port)

	loginURL := fmt.Sprintf("%s/user/login?redirect_uri=%s&state=%s",
		apiURL, url.QueryEscape(redirectURI), state)

	fmt.Fprintln(out, "Opening browser for authentication...")
	fmt.Fprintln(out, iostreams.Gray("If the browser doesn't open, visit:"))
	fmt.Fprintln(out, iostreams.Gray(loginURL))
	fmt.Fprintln(out)

	browser.Open(loginURL)

	fmt.Fprintln(out, "Waiting for login...")

	result, err := api.WaitForCallback(opts.Port, opts.Timeout)
	if err != nil {
		return err
	}

	if result.State != state {
		return fmt.Errorf("state mismatch: expected %s, got %s", state, result.State)
	}

	fmt.Fprintln(out, "Exchanging authorization code...")

	token, err := api.ExchangeCodeForToken(apiURL, result.Code, oauthClient.ClientID, oauthClient.ClientSecret, redirectURI)
	if err != nil {
		return err
	}

	return saveLogin(f, opts, host, oauthClient.ClientID, oauthClient.ClientSecret, token)
}

func saveLogin(f *factory.Factory, opts *loginOptions, host, clientID, clientSecret string, token *api.OAuthToken) error {
	cfg, err := f.Config()
	if err != nil {
		return err
	}

	ctx := &config.Context{
		Host:         host,
		Token:        token.AccessToken,
		RefreshToken: token.RefreshToken,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	if !token.ExpiresAt.IsZero() {
		ctx.ExpiresAt = token.ExpiresAt
	}

	if user, authority := fetchCurrentUser(ctx); user != "" {
		ctx.User = user
		ctx.Authority = authority
	}

	cfg.SetContext(opts.ContextName, ctx)
	cfg.CurrentContext = opts.ContextName

	if err := f.SaveConfig(); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	apiURL := "https://" + host
	fmt.Fprintf(f.IO.Out, "%s Logged in to %s (context: %s) as %s\n",
		iostreams.Green("✓"), apiURL, opts.ContextName, iostreams.Bold(ctx.User))
	return nil
}

func fetchCurrentUser(ctx *config.Context) (string, string) {
	transport := &api.TokenTransport{
		Token: ctx.Token,
		Base:  http.DefaultTransport,
	}
	client := api.NewAPIClient(ctx.BaseURL(), transport, 0)
	body, err := client.Get("/api/users/this", url.Values{})
	if err != nil {
		return "", ""
	}
	var resp struct {
		Result struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			RoleName string `json:"roleName"`
			IsRoot   bool   `json:"isRoot"`
		} `json:"result"`
	}
	if json.Unmarshal(body, &resp) != nil {
		return "", ""
	}
	name := resp.Result.Name
	if name == "" {
		name = resp.Result.Email
	}
	if name == "" {
		name = resp.Result.Phone
	}
	if resp.Result.Email != "" && name != resp.Result.Email {
		name = fmt.Sprintf("%s (%s)", name, resp.Result.Email)
	}

	authority := resp.Result.RoleName
	if resp.Result.IsRoot {
		authority = "root"
	}

	return name, authority
}

