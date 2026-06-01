package org

import (
	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type listOptions struct {
	Email   string
	Creator string
	Status  string
}

func newCmdList(f *factory.Factory) *cobra.Command {
	listFlags := &cmdutil.ListFlags{}
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List organizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, listFlags)
			if opts.Email != "" {
				q.Set("email", opts.Email)
			}
			if opts.Creator != "" {
				q.Set("creator", opts.Creator)
			}
			if opts.Status != "" {
				q.Set("status", opts.Status)
			}

			body, err := client.Get("/api/invpn/org", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	listFlags.Register(cmd)
	cmd.Flags().StringVar(&opts.Email, "email", "", "Filter by email")
	cmd.Flags().StringVar(&opts.Creator, "creator", "", "Filter by creator")
	cmd.Flags().StringVar(&opts.Status, "status", "", "Filter by status")

	return cmd
}
