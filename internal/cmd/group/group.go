package group

import (
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewCmdGroup(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "Manage VPN permission groups",
	}

	cmd.AddCommand(
		newCmdList(f),
		newCmdGet(f),
		newCmdCreate(f),
		newCmdUpdate(f),
		newCmdDelete(f),
		newCmdAccounts(f),
		newCmdRouters(f),
	)

	return cmd
}
