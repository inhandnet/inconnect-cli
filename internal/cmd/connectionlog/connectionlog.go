package connectionlog

import (
	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func NewCmdConnectionLog(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "connection-log",
		Aliases: []string{"conn-log"},
		Short:   "Query VPN connection session logs",
	}

	cmd.AddCommand(newCmdList(f))

	return cmd
}

type listOptions struct {
	cmdutil.ListFlags
	RID      string
	UID      string
	Username string
	Type     string
	Status   string
	After    string
	Before   string
}

func newCmdList(f *factory.Factory) *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List VPN connection session logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewCursorQuery(cmd, &opts.ListFlags)
			if opts.RID != "" {
				q.Set("rid", opts.RID)
			}
			if opts.UID != "" {
				q.Set("uid", opts.UID)
			}
			if opts.Username != "" {
				q.Set("username", opts.Username)
			}
			if opts.Type != "" {
				q.Set("type", opts.Type)
			}
			if opts.Status != "" {
				q.Set("status", opts.Status)
			}
			if t := cmdutil.ParseTimeFlag(opts.After); t != "" {
				q.Set("after", t)
			}
			if t := cmdutil.ParseTimeFlag(opts.Before); t != "" {
				q.Set("before", t)
			}

			body, err := client.Get("/api/invpn/connection-logs", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.RegisterPagination(cmd)
	cmd.Flags().StringVar(&opts.RID, "rid", "", "Filter by router ID")
	cmd.Flags().StringVar(&opts.UID, "uid", "", "Filter by user ID")
	cmd.Flags().StringVar(&opts.Username, "username", "", "Filter by account username (prefix match)")
	cmd.Flags().StringVar(&opts.Type, "type", "", "Filter by account type: user or router")
	cmd.Flags().StringVar(&opts.Status, "status", "", "Filter by session status: active or closed")
	cmd.Flags().StringVar(&opts.After, "after", "", "Only sessions started at/after this time (e.g. 2024-01-01, 2024-01-01T00:00:00Z)")
	cmd.Flags().StringVar(&opts.Before, "before", "", "Only sessions started before this time")

	return cmd
}
