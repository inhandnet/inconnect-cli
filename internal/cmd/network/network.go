package network

import (
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewCmdNetwork(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "network",
		Aliases: []string{"net"},
		Short:   "Manage VPN networks",
	}

	cmd.AddCommand(
		newCmdList(f),
		newCmdGet(f),
		newCmdCreate(f),
		newCmdUpdate(f),
		newCmdDelete(f),
		newCmdMembers(f),
		newCmdAccounts(f),
		newCmdRouters(f),
		newCmdEndpoints(f),
		newCmdCenters(f),
	)

	return cmd
}
