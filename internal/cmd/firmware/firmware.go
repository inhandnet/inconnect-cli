package firmware

import (
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewCmdFirmware(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "firmware",
		Aliases: []string{"fw"},
		Short:   "Manage firmware upgrades",
	}

	cmd.AddCommand(
		newCmdList(f),
		newCmdGet(f),
		newCmdCreate(f),
		newCmdUpdate(f),
		newCmdDelete(f),
		newCmdUpgrade(f),
		newCmdJobStats(f),
		newCmdDevices(f),
	)

	return cmd
}
