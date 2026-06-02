package alert

import (
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewCmdAlert(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alert",
		Short: "Manage VPN alerts and alert rules",
	}

	cmd.AddCommand(
		newCmdList(f),
		newCmdGet(f),
		newCmdAck(f),
		newCmdStats(f),
		newCmdRule(f),
	)

	return cmd
}
