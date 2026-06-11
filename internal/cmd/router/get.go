package router

import (
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdGet(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a VPN router by ID",
		Long: `Get a VPN router by ID.

The router reports two INDEPENDENT status fields — don't confuse them:

  online    (1/0)         Device-management channel (MQTT): whether the
                          platform can reach the device (config push, remote
                          control, ngrok and reboot all rely on it).
  connected (true/false)  VPN tunnel (OpenVPN): whether the device has an
                          active VPN tunnel to its server and can carry VPN
                          traffic.

online=1 with connected=false is common: the device is manageable but has no
VPN tunnel (e.g. its server isn't running).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			body, err := client.Get("/api/invpn/router/"+args[0], q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}
