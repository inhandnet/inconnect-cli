package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/inhandnet/inconnect-cli/internal/api"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
)

func newCmdSwitchOrg(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch-org <org-id>",
		Short: "Switch to a different organization",
		Args:  cobra.ExactArgs(1),
		Example: `  # List your organizations first
  inconnect auth orgs

  # Switch to another org
  inconnect auth switch-org 5e0956c46aa6d10001e931e6`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			name := cfg.ActiveContextName()
			ctx, ok := cfg.Contexts[name]
			if !ok {
				return fmt.Errorf("no active context")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			userBody, err := client.Get("/api2/users/this", nil)
			if err != nil {
				return fmt.Errorf("getting current user: %w", err)
			}
			uid := gjson.GetBytes(userBody, "result._id").String()
			if uid == "" {
				return fmt.Errorf("cannot determine current user ID")
			}

			clientID, clientSecret := ctx.ClientID, ctx.ClientSecret
			if clientID == "" || clientSecret == "" {
				oauthClient, err := api.FetchOAuthClient(cmd.Context(), ctx.APIURL())
				if err != nil {
					return fmt.Errorf("fetching OAuth client: %w", err)
				}
				clientID, clientSecret = oauthClient.ClientID, oauthClient.ClientSecret
			}

			oid := args[0]
			body := map[string]string{
				"clientId":     clientID,
				"clientSecret": clientSecret,
			}

			resp, err := client.Post(fmt.Sprintf("/api/users/%s/orgs/%s/switch", uid, oid), body)
			if err != nil {
				return fmt.Errorf("switch org failed: %w", err)
			}

			var tokenResp struct {
				AccessToken  string `json:"access_token"`
				RefreshToken string `json:"refresh_token"`
				ExpiresIn    int    `json:"expires_in"`
			}
			if err := json.Unmarshal(resp, &tokenResp); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}
			if tokenResp.AccessToken == "" {
				return fmt.Errorf("switch org failed: no access_token in response")
			}

			ctx.Token = tokenResp.AccessToken
			if tokenResp.RefreshToken != "" {
				ctx.RefreshToken = tokenResp.RefreshToken
			}
			if tokenResp.ExpiresIn > 0 {
				ctx.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
			}
			ctx.OrgID = oid

			if err := f.SaveConfig(); err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "%s Switched to organization %s\n", iostreams.Green("✓"), oid)
			return nil
		},
	}

	return cmd
}

func newCmdOrgs(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "orgs",
		Short: "List organizations you belong to",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			userBody, err := client.Get("/api2/users/this", nil)
			if err != nil {
				return fmt.Errorf("getting current user: %w", err)
			}
			uid := gjson.GetBytes(userBody, "result._id").String()
			if uid == "" {
				return fmt.Errorf("cannot determine current user ID")
			}

			body, err := client.Get(fmt.Sprintf("/api/users/%s/orgs", uid), nil)
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")
			return iostreams.FormatOutput(body, f.IO, output)
		},
	}
}
