package org

import (
	"crypto/md5"
	"fmt"

	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type createOptions struct {
	Email    string
	Name     string
	Password string
	Subnet   string
	Country  string
	Industry string
	Deploy   bool
}

func newCmdCreate(f *factory.Factory) *cobra.Command {
	opts := &createOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an organization (admin)",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			pwdHash := fmt.Sprintf("%X", md5.Sum([]byte(opts.Password)))

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

			path := "/api/invpn/org"
			if !opts.Deploy {
				path += "?deploy=false"
			}

			respBody, err := client.Post(path, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteCreated(f, "Organization", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.Email, "email", "", "Admin email address (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Organization name (required)")
	cmd.Flags().StringVar(&opts.Password, "password", "", "Initial login password (required)")
	cmd.Flags().StringVar(&opts.Subnet, "subnet", "10.16.0.0/12", "VPN subnet CIDR")
	cmd.Flags().StringVar(&opts.Country, "country", "", "Country or region")
	cmd.Flags().StringVar(&opts.Industry, "industry", "", "Industry")
	cmd.Flags().BoolVar(&opts.Deploy, "deploy", true, "Auto-deploy VPN server after creation")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("password")

	return cmd
}
