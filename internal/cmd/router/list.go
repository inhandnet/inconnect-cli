package router

import (
	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
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
		Long: `List VPN routers.

Each router reports two INDEPENDENT status fields — don't confuse them:

  online    (1/0)         Device-management channel (MQTT): whether the
                          platform can reach the device. Config push, remote
                          control, ngrok and reboot all depend on this.
  connected (true/false)  VPN tunnel (OpenVPN): whether the device has an
                          active VPN tunnel to its server and can carry VPN
                          traffic.

A device is commonly online=1 while connected=false (manageable, but no VPN
tunnel — e.g. its server isn't running). Filter with --online / --connected.`,
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
