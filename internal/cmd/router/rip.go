package router

import (
	"fmt"
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdNextVip(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "next-vip <router-id>",
		Short: "Get next available endpoint VIP for a router",
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

			body, err := client.Get("/api/invpn/router/"+args[0]+"/endpoint/vip", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}

type setRipOptions struct {
	RealIPAddress string
	Enable        bool
	Center        bool
}

func newCmdSetRip(f *factory.Factory) *cobra.Command {
	opts := &setRipOptions{}

	cmd := &cobra.Command{
		Use:   "set-rip <router-id>",
		Short: "Set real IP address for a router",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if opts.RealIPAddress != "" {
				body["realIPAddress"] = opts.RealIPAddress
			}
			if cmd.Flags().Changed("enable") {
				body["enable"] = opts.Enable
			}
			if cmd.Flags().Changed("center") {
				body["center"] = opts.Center
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			respBody, err := client.Do("PUT", "/api/invpn/routers/"+args[0]+"/rip", q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Real IP updated for router %s\n", args[0])
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.RealIPAddress, "ip", "", "Real IP address")
	cmd.Flags().BoolVar(&opts.Enable, "enable", true, "Enable real IP")
	cmd.Flags().BoolVar(&opts.Center, "center", false, "Set as center router")

	return cmd
}

func newCmdNatConf(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "nat-conf <router-id>",
		Short: "Download NAT config for a router",
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

			body, err := client.Get("/api/invpn/router/"+args[0]+"/nat.conf", q)
			if err != nil {
				return err
			}

			_, err = f.IO.Out.Write(body)
			return err
		},
	}
}

func newCmdLocations(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "locations",
		Short: "Get router locations",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			body, err := client.Get("/api/invpn/routers/locations", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}
