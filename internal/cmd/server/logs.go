package server

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type logsOptions struct {
	Tail   int
	Since  string
	Format string
}

func newCmdLogs(f *factory.Factory) *cobra.Command {
	opts := &logsOptions{}

	cmd := &cobra.Command{
		Use:   "logs <server-id>",
		Short: "Stream the org's OpenVPN server Pod logs (current Pod lifecycle only)",
		Long: `Read the OpenVPN server Pod logs in real time (not persisted).

--tail and --since are two mutually exclusive query modes: --tail returns the
last N lines, --since returns logs from a time offset (server caps the volume).
Without flags the server defaults to the last 200 lines.

Only logs from the current Pod lifecycle are available; a Pod restart loses
previous output. Logs may contain client real IP, so this requires admin.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			oid, _ := cmd.Flags().GetString("oid")
			if oid == "" {
				return fmt.Errorf("--oid is required (or set an active org via 'auth switch-org')")
			}

			q := url.Values{}
			q.Set("oid", oid)
			if cmd.Flags().Changed("tail") {
				q.Set("tail", strconv.Itoa(opts.Tail))
			}
			if opts.Since != "" {
				q.Set("since", opts.Since)
			}
			if opts.Format != "" {
				q.Set("format", opts.Format)
			}

			body, err := client.Get("/api/invpn/servers/"+args[0]+"/logs", q)
			if err != nil {
				return err
			}

			if opts.Format == "json" {
				return iostreams.FormatOutput(body, f.IO, f.IO.Output)
			}
			_, err = f.IO.Out.Write(body)
			return err
		},
	}

	cmd.Flags().IntVar(&opts.Tail, "tail", 200, "Return the last N lines (max 2000); mutually exclusive with --since")
	cmd.Flags().StringVar(&opts.Since, "since", "", "Logs since a time offset (e.g. 10m, 1h); mutually exclusive with --tail")
	cmd.Flags().StringVar(&opts.Format, "format", "text", "Output format: text (raw) or json (line-wrapped)")
	cmd.MarkFlagsMutuallyExclusive("tail", "since")

	return cmd
}
