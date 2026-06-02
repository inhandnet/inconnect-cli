package drc

import (
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewCmdDRC(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "drc",
		Aliases: []string{"config-template"},
		Short:   "Manage device running config templates",
	}

	cmd.AddCommand(
		newCmdList(f),
		newCmdGet(f),
		newCmdCreate(f),
		newCmdUpdate(f),
		newCmdDelete(f),
		newCmdDevices(f),
	)

	return cmd
}
