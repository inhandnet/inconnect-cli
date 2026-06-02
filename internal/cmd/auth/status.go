package auth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/inhandnet/inconnect-cli/internal/api"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdStatus(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			name := cfg.ActiveContextName()
			if name == "" {
				fmt.Fprintln(f.IO.Out, "No active context")
				return nil
			}

			ctx, ok := cfg.Contexts[name]
			if !ok {
				return fmt.Errorf("context %q not found", name)
			}

			out := f.IO.Out
			fmt.Fprintf(out, "Context:  %s\n", iostreams.Bold(name))
			fmt.Fprintf(out, "API:      %s\n", ctx.BaseURL())
			if ctx.User != "" {
				fmt.Fprintf(out, "User:     %s\n", ctx.User)
			}

			tokenExpired := !ctx.ExpiresAt.IsZero() && ctx.ExpiresAt.Before(time.Now())

			if tokenExpired && ctx.RefreshToken != "" {
				clientID, clientSecret := ctx.ClientID, ctx.ClientSecret
				if clientID == "" || clientSecret == "" {
					if c, err := api.FetchOAuthClient(cmd.Context(), ctx.BaseURL()); err == nil {
						clientID, clientSecret = c.ClientID, c.ClientSecret
					}
				}
				newToken, err := api.RefreshAccessToken(ctx.BaseURL(), clientID, clientSecret, ctx.RefreshToken)
				if err == nil {
					ctx.Token = newToken.AccessToken
					if newToken.RefreshToken != "" {
						ctx.RefreshToken = newToken.RefreshToken
					}
					if !newToken.ExpiresAt.IsZero() {
						ctx.ExpiresAt = newToken.ExpiresAt
					}
					_ = f.SaveConfig()
					tokenExpired = false
					fmt.Fprintf(out, "Status:   %s\n", iostreams.Green("logged in (token refreshed)"))
				} else {
					fmt.Fprintf(out, "Status:   %s\n", iostreams.Red("token expired, refresh failed — please login again"))
				}
			} else {
				switch {
				case ctx.Token == "":
					fmt.Fprintf(out, "Status:   %s\n", iostreams.Red("not logged in"))
				case tokenExpired:
					fmt.Fprintf(out, "Status:   %s\n", iostreams.Red("token expired, please login again"))
				case ctx.IsImpersonating():
					fmt.Fprintf(out, "Status:   %s\n", iostreams.Yellow("logged in (impersonating)"))
				default:
					fmt.Fprintf(out, "Status:   %s\n", iostreams.Green("logged in"))
				}
			}

			if !ctx.ExpiresAt.IsZero() {
				fmt.Fprintf(out, "Expires:  %s\n", ctx.ExpiresAt.Local().Format("2006-01-02 15:04:05"))
			}

			if ctx.EffectiveToken() != "" && !tokenExpired {
				client, err := f.APIClient()
				if err != nil {
					return nil
				}

				body, err := client.Get("/api/users/this", nil)
				if err == nil {
					var resp struct {
						Result struct {
							Name     string `json:"name"`
							Email    string `json:"email"`
							Phone    string `json:"phone"`
							RoleName string `json:"roleName"`
							Oid      string `json:"oid"`
						} `json:"result"`
					}
					if json.Unmarshal(body, &resp) == nil {
						me := resp.Result
						account := me.Name
						if account == "" {
							account = me.Email
						}
						if account == "" {
							account = me.Phone
						}
						if account != "" {
							if me.Email != "" && account != me.Email {
								account += fmt.Sprintf(" (%s)", me.Email)
							}
							fmt.Fprintf(out, "Account:  %s\n", account)
						}
						if me.RoleName != "" {
							fmt.Fprintf(out, "Role:     %s\n", me.RoleName)
						}
					}
				}

				body, err = client.Get("/api/orgs/this", url.Values{})
				if err == nil {
					var resp struct {
						Result struct {
							ID   string `json:"_id"`
							Name string `json:"name"`
						} `json:"result"`
					}
					if json.Unmarshal(body, &resp) == nil && resp.Result.Name != "" {
						fmt.Fprintf(out, "Org:      %s (%s)\n", resp.Result.Name, resp.Result.ID)
					}
				}
			}

			return nil
		},
	}
}
