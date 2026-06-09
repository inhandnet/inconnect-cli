package router

import (
	"fmt"
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdTunnelClose(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tunnel-close <uuid>",
		Short: "Close an ngrok tunnel opened by `router web` or `router ssh`",
		Long: `Explicitly close an ngrok tunnel by its tunnel id (uuid), instead of
waiting for the device-side idle timeout to reclaim it. The uuid is printed
by the "router web" and "router ssh" commands when the tunnel is opened.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			uuid := args[0]

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			respBody, err := client.Do("DELETE", "/api/touch/ngrok/tunnels/"+uuid, q, nil)
			if err != nil {
				if respBody != nil {
					_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				}
				return err
			}

			fmt.Fprintf(f.IO.Out, "Tunnel %s closed.\n", uuid)
			return nil
		},
	}

	return cmd
}
