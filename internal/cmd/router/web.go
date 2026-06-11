package router

import (
	"fmt"
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/browser"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

func newCmdWeb(f *factory.Factory) *cobra.Command {
	var (
		noBrowser  bool
		port       int
		server     string
		timeoutSec int
	)

	cmd := &cobra.Command{
		Use:   "web <id>",
		Short: "Open device web management UI via ngrok tunnel",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			server, err = resolveNgrokServer(f, server)
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			fmt.Fprintf(f.IO.ErrOut, "Establishing ngrok tunnel to device %s...\n", args[0])

			body := map[string]any{
				"name":     "ngrok connect",
				"type":     23,
				"objectId": args[0],
				"priority": 30,
				"timeout":  timeoutSec * 1000,
				"data": map[string]any{
					"server": server,
					"proto":  "http",
					"port":   port,
				},
			}

			respBody, err := client.Do("POST", "/api2/tasks/run", q, body)
			if err != nil {
				if respBody != nil {
					_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				}
				return err
			}

			if taskErr := gjson.GetBytes(respBody, "result.error").String(); taskErr != "" {
				return fmt.Errorf("%s", taskErr)
			}

			tunnelURL := gjson.GetBytes(respBody, "result.data.response").String()
			if tunnelURL == "" {
				return fmt.Errorf("no tunnel URL returned")
			}

			fmt.Fprintln(f.IO.Out, tunnelURL)

			// The ngrok server encodes stateless web tunnels with the tunnel
			// uuid as the leading DNS label (https://<uuid>.<domain>), so the
			// uuid can be recovered here and used to close the tunnel.
			if u, perr := url.Parse(tunnelURL); perr == nil {
				if leaf, _ := splitFirstLabel(u.Hostname()); isHexUUID(leaf) {
					fmt.Fprintf(f.IO.Out, "Tunnel ID: %s   (close with: inconnect router tunnel-close %s)\n", leaf, leaf)
				}
			}

			if !noBrowser {
				browser.Open(tunnelURL)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&noBrowser, "no-browser", false, "Print URL only, don't open browser")
	cmd.Flags().IntVar(&port, "port", 80, "Device-side web port to expose")
	cmd.Flags().StringVar(&server, "server", "", "Ngrok server address (default: auto-detected from the active context's host)")
	cmd.Flags().IntVar(&timeoutSec, "timeout", 60, "Task timeout in seconds")

	return cmd
}
