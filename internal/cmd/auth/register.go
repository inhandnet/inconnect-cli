package auth

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/inhandnet/ics-cli/internal/api"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

type registerOptions struct {
	Host     string
	Email    string
	Name     string
	Password string
	Subnet   string
	Country  string
	Industry string
	Deploy   bool
}

func newCmdRegister(f *factory.Factory) *cobra.Command {
	opts := &registerOptions{}

	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new organization account",
		Example: `  # Register on China region (default)
  ics auth register --email user@example.com --name "My Org" --password "P@ssw0rd"

  # Register on dev environment
  ics auth register --host dev --email user@example.com --name "My Org" --password "P@ssw0rd"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			host, err := resolveHost(opts.Host)
			if err != nil {
				return err
			}
			apiURL := "https://" + host

			client := api.NewAPIClient(apiURL, &api.TokenTransport{Base: http.DefaultTransport}, 0)

			pwdHash := fmt.Sprintf("%x", md5.Sum([]byte(opts.Password)))

			body := map[string]any{
				"email":    opts.Email,
				"name":     opts.Name,
				"password": pwdHash,
				"subnet":   opts.Subnet,
			}

			metadata := map[string]string{}
			if opts.Country != "" {
				metadata["country"] = opts.Country
			}
			if opts.Industry != "" {
				metadata["industry"] = opts.Industry
			}
			if len(metadata) > 0 {
				body["metadata"] = metadata
			}

			path := "/api/invpn/register"
			if !opts.Deploy {
				path += "?deploy=false"
			}

			respBody, err := client.Post(path, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			if gjson.GetBytes(respBody, "error").Exists() {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return fmt.Errorf("registration failed: %s", gjson.GetBytes(respBody, "error").String())
			}

			fmt.Fprintf(f.IO.ErrOut, "Registration successful\n")
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.Host, "host", "cn", `Platform region: "cn", "us", "eu", "dev", "beta", or a custom domain`)
	cmd.Flags().StringVar(&opts.Email, "email", "", "Email address (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Organization name (required)")
	cmd.Flags().StringVar(&opts.Password, "password", "", "Login password (required)")
	cmd.Flags().StringVar(&opts.Subnet, "subnet", "10.16.0.0/12", "VPN subnet CIDR")
	cmd.Flags().StringVar(&opts.Country, "country", "", "Country or region")
	cmd.Flags().StringVar(&opts.Industry, "industry", "", "Industry")
	cmd.Flags().BoolVar(&opts.Deploy, "deploy", true, "Auto-deploy VPN server after registration")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("password")

	return cmd
}
