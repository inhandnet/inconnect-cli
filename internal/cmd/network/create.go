package network

import (
	"net/url"

	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type createOptions struct {
	Name          string
	Description   string
	Type          string
	Center        string
	RealIPAddress string
	SiteIsolate   bool
}

func newCmdCreate(f *factory.Factory) *cobra.Command {
	opts := &createOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a VPN network",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"name": opts.Name,
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

			respBody, err := client.Do("POST", "/api/invpn/networks", q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteCreated(f, "Network", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Network name (required)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Network description")
	cmd.Flags().StringVar(&opts.Type, "type", "", "Network type (mesh or star)")
	cmd.Flags().StringVar(&opts.Center, "center", "", "Center router ID (for star topology)")
	cmd.Flags().StringVar(&opts.RealIPAddress, "real-ip", "", "Real IP address")
	cmd.Flags().BoolVar(&opts.SiteIsolate, "site-isolate", false, "Enable site isolation")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}
