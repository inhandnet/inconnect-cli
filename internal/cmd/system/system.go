package system

import (
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func NewCmdSystem(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "system",
		Aliases: []string{"sys"},
		Short:   "System information and management",
	}

	cmd.AddCommand(
		newCmdVersions(f),
		newCmdService(f),
	)

	return cmd
}

func newCmdVersions(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "versions",
		Short: "List all backend service versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body, err := client.Get("/api/elms/services", nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}

func newCmdService(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "service <name>",
		Short: "Get instances of a specific service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body, err := client.Get("/api/elms/services/"+args[0], nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}
