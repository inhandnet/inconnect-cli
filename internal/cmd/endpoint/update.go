package endpoint

import (
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type updateOptions struct {
	RouterID string
	Name     string
	IP       string
	VIP      string
	MAC      string
}

func newCmdUpdate(f *factory.Factory) *cobra.Command {
	opts := &updateOptions{}

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an endpoint",
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
			if opts.IP != "" {
				body["ip"] = opts.IP
			}
			if opts.VIP != "" {
				body["vip"] = opts.VIP
			}
			if opts.MAC != "" {
				body["mac"] = opts.MAC
			}

			respBody, err := client.Do("PUT", "/api/invpn/router/"+opts.RouterID+"/endpoint/"+args[0], q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteUpdated(f, "Endpoint", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.RouterID, "router-id", "", "Router ID (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Endpoint name")
	cmd.Flags().StringVar(&opts.IP, "ip", "", "Real IP address")
	cmd.Flags().StringVar(&opts.VIP, "vip", "", "Virtual IP address")
	cmd.Flags().StringVar(&opts.MAC, "mac", "", "MAC address")
	_ = cmd.MarkFlagRequired("router-id")

	return cmd
}
