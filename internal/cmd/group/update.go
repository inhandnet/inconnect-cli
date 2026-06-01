package group

import (
	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdUpdate(f *factory.Factory) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a permission group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if name != "" {
				body["name"] = name
			}

			respBody, err := client.Put("/api/invpn/groups/"+args[0], body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteUpdated(f, "Group", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Group name")

	return cmd
}
