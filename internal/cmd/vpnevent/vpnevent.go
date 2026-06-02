package vpnevent

import (
	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func NewCmdVpnEvent(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vpn-event",
		Aliases: []string{"vpn-events"},
		Short:   "Query VPN auth/connection event stream",
	}

	cmd.AddCommand(newCmdList(f))

	return cmd
}

type listOptions struct {
	cmdutil.ListFlags
	Type     string
	RID      string
	UID      string
	Username string
	After    string
	Before   string
}

func newCmdList(f *factory.Factory) *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List VPN events (connected/disconnected/auth_failed)",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewCursorQuery(cmd, &opts.ListFlags)
			if opts.Type != "" {
				q.Set("type", opts.Type)
			}
			if opts.RID != "" {
				q.Set("rid", opts.RID)
			}
			if opts.UID != "" {
				q.Set("uid", opts.UID)
			}
			if opts.Username != "" {
				q.Set("username", opts.Username)
			}
			if t := cmdutil.ParseTimeFlag(opts.After); t != "" {
				q.Set("after", t)
			}
			if t := cmdutil.ParseTimeFlag(opts.Before); t != "" {
				q.Set("before", t)
			}

			body, err := client.Get("/api/invpn/vpn-events", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.RegisterPagination(cmd)
	cmd.Flags().StringVar(&opts.Type, "type", "", "Filter by event type, comma-separated (e.g. auth_failed,connected)")
	cmd.Flags().StringVar(&opts.RID, "rid", "", "Filter by router ID")
	cmd.Flags().StringVar(&opts.UID, "uid", "", "Filter by user ID")
	cmd.Flags().StringVar(&opts.Username, "username", "", "Filter by username (prefix match)")
	cmd.Flags().StringVar(&opts.After, "after", "", "Only events at/after this time (e.g. 2024-01-01T00:00:00Z)")
	cmd.Flags().StringVar(&opts.Before, "before", "", "Only events before this time")

	return cmd
}
