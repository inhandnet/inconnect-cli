package role

import (
	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func NewCmdRole(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "List user roles",
		Long:  "List user roles.\n\nUse the role ID ('_id') with 'user create --role-id'.",
	}

	cmd.AddCommand(
		newCmdList(f),
		newCmdGet(f),
	)

	return cmd
}

func newCmdList(f *factory.Factory) *cobra.Command {
	lf := &cmdutil.ListFlags{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List roles in the organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body, err := client.Get("/api/roles", cmdutil.NewQuery(cmd, lf))
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	lf.RegisterPagination(cmd)

	return cmd
}

func newCmdGet(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a role by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body, err := client.Get("/api/roles/"+args[0], nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}
