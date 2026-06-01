package router

import (
	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type listOptions struct {
	cmdutil.ListFlags
	SerialNumber string
	VIP          string
	RIP          string
	Connected    string
	Online       string
	Query        string
}

func newCmdList(f *factory.Factory) *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List VPN routers",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, &opts.ListFlags)
			if opts.SerialNumber != "" {
				q.Set("serialNumber", opts.SerialNumber)
			}
			if opts.VIP != "" {
				q.Set("vip", opts.VIP)
			}
			if opts.RIP != "" {
				q.Set("rip", opts.RIP)
			}
			if opts.Connected != "" {
				q.Set("connected", opts.Connected)
			}
			if opts.Online != "" {
				q.Set("online", opts.Online)
			}
			if opts.Query != "" {
				q.Set("q", opts.Query)
			}

			body, err := client.Get("/api/invpn/routers", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&opts.SerialNumber, "serial", "", "Filter by serial number")
	cmd.Flags().StringVar(&opts.VIP, "vip", "", "Filter by virtual IP")
	cmd.Flags().StringVar(&opts.RIP, "rip", "", "Filter by real IP address")
	cmd.Flags().StringVar(&opts.Connected, "connected", "", "Filter by connection status (true/false)")
	cmd.Flags().StringVar(&opts.Online, "online", "", "Filter by online status")
	cmd.Flags().StringVar(&opts.Query, "query", "", "Search across name and serial number")

	return cmd
}
