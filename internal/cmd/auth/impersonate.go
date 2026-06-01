package auth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/inhandnet/ics-cli/internal/api"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
)

func newCmdImpersonate(f *factory.Factory) *cobra.Command {
	var (
		userID string
		orgID  string
		stop   bool
	)

	cmd := &cobra.Command{
		Use:   "impersonate",
		Short: "Impersonate another user (requires ROOT privilege)",
		Example: `  # Impersonate by org ID (auto-resolves org admin)
  ics auth impersonate --org 5e0956c46aa6d10001e931e6

  # Impersonate a specific user in an org
  ics auth impersonate --org <oid> --user <uid>

  # Stop impersonation and restore admin identity
  ics auth impersonate --stop`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			name := cfg.ActiveContextName()
			ctx, ok := cfg.Contexts[name]
			if !ok {
				return fmt.Errorf("no active context; run 'ics auth login' first")
			}

			if stop {
				if !ctx.IsImpersonating() {
					return fmt.Errorf("not currently impersonating")
				}
				ctx.Token = ctx.AdminToken
				ctx.RefreshToken = ctx.AdminRefreshToken
				ctx.ExpiresAt = ctx.AdminExpiresAt
				ctx.AdminToken = ""
				ctx.AdminRefreshToken = ""
				ctx.AdminExpiresAt = time.Time{}
				ctx.OrgID = ""
				if err := f.SaveConfig(); err != nil {
					return err
				}
				fmt.Fprintf(f.IO.Out, "%s Impersonation stopped, admin identity restored\n", iostreams.Green("✓"))
				return nil
			}

			if orgID == "" {
				return fmt.Errorf("--org is required (or use --stop)")
			}

			if ctx.IsImpersonating() {
				return fmt.Errorf("already impersonating; run 'ics auth impersonate --stop' first")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			realOrgID, err := resolveOrgID(client, orgID)
			if err != nil {
				return fmt.Errorf("resolving org: %w", err)
			}

			if userID == "" {
				resolved, err := resolveUserForOrg(client, realOrgID)
				if err != nil {
					return fmt.Errorf("resolving user for org: %w", err)
				}
				userID = resolved
			}

			// Refresh token before impersonation — the access_token is passed as a query param
			// which TokenTransport cannot update on 401 retry.
			if !ctx.ExpiresAt.IsZero() && time.Now().After(ctx.ExpiresAt) && ctx.RefreshToken != "" {
				newToken, refreshErr := api.RefreshAccessToken(
					ctx.APIURL(), ctx.ClientID, ctx.ClientSecret, ctx.RefreshToken)
				if refreshErr != nil {
					return fmt.Errorf("token expired and refresh failed: %w\nHint: run 'ics auth login' to re-authenticate", refreshErr)
				}
				ctx.Token = newToken.AccessToken
				if newToken.RefreshToken != "" {
					ctx.RefreshToken = newToken.RefreshToken
				}
				ctx.ExpiresAt = newToken.ExpiresAt
				_ = f.SaveConfig()
			}

			token := ctx.EffectiveToken()
			q := url.Values{}
			q.Set("oid", orgID)
			q.Set("uid", userID)
			q.Set("access_token", token)
			q.Set("verbose", "100")

			body, err := client.Post("/api/token/impersonate?"+q.Encode(), nil)
			if err != nil {
				return fmt.Errorf("impersonate failed: %w", err)
			}

			var tokenResp struct {
				AccessToken  string `json:"access_token"`
				RefreshToken string `json:"refresh_token"`
				ExpiresIn    int64  `json:"expires_in"`
				Error        string `json:"error"`
				ErrorCode    int    `json:"error_code"`
			}
			if err := json.Unmarshal(body, &tokenResp); err != nil {
				return fmt.Errorf("parsing token response: %w", err)
			}
			if tokenResp.Error != "" {
				return fmt.Errorf("impersonate failed: %s (code: %d)", tokenResp.Error, tokenResp.ErrorCode)
			}
			if tokenResp.AccessToken == "" {
				return fmt.Errorf("impersonate failed: no access_token in response")
			}

			ctx.AdminToken = ctx.Token
			ctx.AdminRefreshToken = ctx.RefreshToken
			ctx.AdminExpiresAt = ctx.ExpiresAt

			ctx.Token = tokenResp.AccessToken
			ctx.RefreshToken = tokenResp.RefreshToken
			ctx.OrgID = realOrgID

			if err := f.SaveConfig(); err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "%s Impersonating user %s in org %s\n", iostreams.Green("✓"), userID, orgID)
			fmt.Fprintf(f.IO.Out, "Run 'ics auth impersonate --stop' to restore admin identity\n")
			return nil
		},
	}

	cmd.Flags().StringVar(&userID, "user", "", "User ID to impersonate")
	cmd.Flags().StringVar(&orgID, "org", "", "Organization ID")
	cmd.Flags().BoolVar(&stop, "stop", false, "Stop impersonation and restore admin identity")

	return cmd
}

func resolveOrgID(client *api.APIClient, oid string) (string, error) {
	if len(oid) == 24 {
		return oid, nil
	}

	q := url.Values{}
	q.Set("name", oid)
	body, err := client.Get("/api2/organizations", q)
	if err != nil {
		return "", err
	}

	var matches []string
	for _, org := range gjson.GetBytes(body, "result").Array() {
		if org.Get("name").String() == oid {
			if id := org.Get("_id").String(); id != "" {
				matches = append(matches, id)
			}
		}
	}

	switch len(matches) {
	case 0:
		return "", fmt.Errorf("organization %q not found", oid)
	case 1:
		return matches[0], nil
	default:
		return "", fmt.Errorf("multiple organizations named %q found: %v\nHint: use the full org ID instead (e.g. --org %s)", oid, matches, matches[0])
	}
}

func resolveUserForOrg(client *api.APIClient, oid string) (string, error) {
	q := url.Values{}
	q.Set("oid", oid)
	q.Set("limit", "0")

	body, err := client.Get("/api2/users", q)
	if err != nil {
		return "", err
	}

	results := gjson.GetBytes(body, "result")
	if !results.Exists() || len(results.Array()) == 0 {
		return "", fmt.Errorf("no users found for org %s", oid)
	}

	var candidates []gjson.Result
	for _, user := range results.Array() {
		if user.Get("oid").String() == oid {
			candidates = append(candidates, user)
		}
	}
	if len(candidates) == 0 {
		candidates = results.Array()
	}

	for _, user := range candidates {
		if user.Get("roleName").String() == "admin" {
			uid := user.Get("_id").String()
			if uid != "" {
				return uid, nil
			}
		}
	}

	uid := candidates[0].Get("_id").String()
	if uid == "" {
		return "", fmt.Errorf("no user ID found for org %s", oid)
	}
	return uid, nil
}
