package endpoint

import (
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type createOptions struct {
	RouterID string
	Name     string
	IP       string
	VIP      string
	MAC      string
}

func newCmdCreate(f *factory.Factory) *cobra.Command {
	opts := &createOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an endpoint on a router",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			body := map[string]any{
				"ip":   opts.IP,
				"name": opts.Name,
			}
			if opts.VIP != "" {
				body["vip"] = opts.VIP
			}
			if opts.MAC != "" {
				body["mac"] = opts.MAC
			}

			respBody, err := client.Do("POST", "/api/invpn/router/"+opts.RouterID+"/endpoint", q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteCreated(f, "Endpoint", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.RouterID, "router-id", "", "Router ID (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Endpoint name (required)")
	cmd.Flags().StringVar(&opts.IP, "ip", "", "Real IP address (required)")
	cmd.Flags().StringVar(&opts.VIP, "vip", "", "Virtual IP address")
	cmd.Flags().StringVar(&opts.MAC, "mac", "", "MAC address")
	_ = cmd.MarkFlagRequired("router-id")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("ip")

	return cmd
}
