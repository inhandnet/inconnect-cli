package user

import (
	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type listOptions struct {
	cmdutil.ListFlags
	Email string
	Group string
	Query string
}

func newCmdList(f *factory.Factory) *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List VPN user accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, &opts.ListFlags)
			if opts.Email != "" {
				q.Set("email", opts.Email)
			}
			if opts.Group != "" {
				q.Set("group", opts.Group)
			}
			if opts.Query != "" {
				q.Set("q", opts.Query)
			}

			body, err := client.Get("/api/invpn/users", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&opts.Email, "email", "", "Filter by email")
	cmd.Flags().StringVar(&opts.Group, "group", "", "Filter by group")
	cmd.Flags().StringVar(&opts.Query, "query", "", "Search across name and email")

	return cmd
}
