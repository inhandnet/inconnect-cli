package alert

import (
	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type listOptions struct {
	cmdutil.ListFlags
	RouterName string
	After      string
	Before     string
	RID        string
	Ack        string
}

func newCmdList(f *factory.Factory) *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List alerts",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, &opts.ListFlags)
			if opts.RouterName != "" {
				q.Set("routerName", opts.RouterName)
			}
			if opts.After != "" {
				q.Set("from", cmdutil.ParseTimeFlag(opts.After))
			}
			if opts.Before != "" {
				q.Set("to", cmdutil.ParseTimeFlag(opts.Before))
			}
			if opts.RID != "" {
				q.Set("rid", opts.RID)
			}
			if opts.Ack != "" {
				q.Set("ack", opts.Ack)
			}

			body, err := client.Get("/api/invpn/alerts", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&opts.RouterName, "router-name", "", "Filter by router name")
	cmd.Flags().StringVar(&opts.After, "after", "", "Start time (e.g. 2024-01-01, 2024-01-01T00:00:00Z)")
	cmd.Flags().StringVar(&opts.Before, "before", "", "End time (e.g. 2024-12-31, 2024-12-31T23:59:59Z)")
	cmd.Flags().StringVar(&opts.RID, "rid", "", "Filter by router ID")
	cmd.Flags().StringVar(&opts.Ack, "ack", "", "Filter by ack status (true/false)")

	return cmd
}
