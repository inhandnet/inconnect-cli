package config

import (
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewCmdConfig(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration contexts",
	}

	cmd.AddCommand(
		newCmdCurrentContext(f),
		newCmdListContexts(f),
		newCmdUseContext(f),
		newCmdSetContext(f),
		newCmdDeleteContext(f),
	)

	return cmd
}
