package network

import (
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type updateOptions struct {
	Name          string
	Description   string
	Type          string
	Center        string
	RealIPAddress string
	SiteIsolate   bool
}

func newCmdUpdate(f *factory.Factory) *cobra.Command {
	opts := &updateOptions{}

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a VPN network",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if opts.Name != "" {
				body["name"] = opts.Name
			}
			if opts.Description != "" {
				body["description"] = opts.Description
			}
			if opts.Type != "" {
				body["type"] = opts.Type
			}
			if opts.Center != "" {
				body["center"] = opts.Center
			}
			if opts.RealIPAddress != "" {
				body["realIPAddress"] = opts.RealIPAddress
			}
			if cmd.Flags().Changed("site-isolate") {
				body["siteIsolate"] = opts.SiteIsolate
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			respBody, err := client.Do("PUT", "/api/invpn/networks/"+args[0], q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteUpdated(f, "Network", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Network name")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Network description")
	cmd.Flags().StringVar(&opts.Type, "type", "", "Network type (mesh or star)")
	cmd.Flags().StringVar(&opts.Center, "center", "", "Center router ID (for star topology)")
	cmd.Flags().StringVar(&opts.RealIPAddress, "real-ip", "", "Real IP address")
	cmd.Flags().BoolVar(&opts.SiteIsolate, "site-isolate", false, "Enable site isolation")

	return cmd
}
