package router

import (
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type updateOptions struct {
	Name         string
	LanInterface string
	MobileNumber string
	Address      string
	Subnet       string
	CustomFields []string
}

func newCmdUpdate(f *factory.Factory) *cobra.Command {
	opts := &updateOptions{}

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a VPN router",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			body := map[string]any{}
			if opts.Name != "" {
				body["name"] = opts.Name
			}
			if opts.LanInterface != "" {
				body["lanInterface"] = opts.LanInterface
			}
			if opts.MobileNumber != "" {
				body["mobileNumber"] = opts.MobileNumber
			}
			if opts.Address != "" {
				body["address"] = opts.Address
			}
			if opts.Subnet != "" {
				body["subnet"] = opts.Subnet
			}
			if len(opts.CustomFields) > 0 {
				cf := parseKeyValues(opts.CustomFields)
				if len(cf) > 0 {
					body["customFields"] = cf
				}
			}

			respBody, err := client.Do("PUT", "/api/invpn/router/"+args[0], q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteUpdated(f, "Router", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Router name")
	cmd.Flags().StringVar(&opts.LanInterface, "lan-interface", "", "LAN interface")
	cmd.Flags().StringVar(&opts.MobileNumber, "mobile-number", "", "Mobile number")
	cmd.Flags().StringVar(&opts.Address, "address", "", "Physical address")
	cmd.Flags().StringVar(&opts.Subnet, "subnet", "", "Router subnet")
	cmd.Flags().StringSliceVar(&opts.CustomFields, "custom-field", nil, "Custom field (key=value, can specify multiple)")

	return cmd
}
