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
		proto      string
		port       int
		server     string
		timeoutSec int
	)

	cmd := &cobra.Command{
		Use:   "web <id>",
		Short: "Open device web management via ngrok tunnel",
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

			fmt.Fprintf(f.IO.ErrOut, "Establishing ngrok tunnel to device %s...\n", args[0])

			// Same task as the web UI: type=23 (ngrok connect) via /api2/tasks/run.
			// The tunnel URL is returned in result.data.response.
			body := map[string]any{
				"name":     "ngrok connect",
				"type":     23,
				"objectId": args[0],
				"priority": 30,
				"timeout":  timeoutSec * 1000,
				"data": map[string]any{
					"server": server,
					"proto":  proto,
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

			if !noBrowser {
				browser.Open(tunnelURL)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&noBrowser, "no-browser", false, "Print URL only, don't open browser")
	cmd.Flags().StringVar(&proto, "proto", "http", "Tunnel protocol")
	cmd.Flags().IntVar(&port, "port", 80, "Device-side port to expose")
	cmd.Flags().StringVar(&server, "server", "ngrok.j3r0lin.com:4443", "Ngrok server address")
	cmd.Flags().IntVar(&timeoutSec, "timeout", 60, "Task timeout in seconds")

	return cmd
}
