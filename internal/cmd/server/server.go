package server

import (
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewCmdServer(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server",
		Aliases: []string{"srv"},
		Short:   "Manage VPN servers",
	}

	cmd.PersistentFlags().Bool("show-secrets", false, "Reveal private keys in output (redacted by default)")

	cmd.AddCommand(
		newCmdList(f),
		newCmdGet(f),
		newCmdCreate(f),
		newCmdUpdate(f),
		newCmdDelete(f),
		newCmdDeploy(f),
		newCmdCommand(f),
		newCmdNetworks(f),
		newCmdStop(f),
		newCmdRecover(f),
		newCmdIssueKeypair(f),
	)

	return cmd
}
