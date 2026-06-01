package registerlog

import (
	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func NewCmdRegisterLog(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "register-log",
		Aliases: []string{"reglog"},
		Short:   "View device registration logs",
	}

	cmd.AddCommand(
		newCmdList(f),
	)

	return cmd
}

func newCmdList(f *factory.Factory) *cobra.Command {
	lf := &cmdutil.ListFlags{}

	cmd := &cobra.Command{
		Use:   "list <serial-number>",
		Short: "List registration events for a device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewCursorQuery(cmd, lf)

			body, err := client.Get("/api/devices/"+args[0]+"/register-events", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	lf.RegisterPagination(cmd)

	return cmd
}
