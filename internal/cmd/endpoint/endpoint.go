package endpoint

import (
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewCmdEndpoint(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "endpoint",
		Aliases: []string{"ep"},
		Short:   "Manage VPN endpoints on routers",
	}

	cmd.AddCommand(
		newCmdList(f),
		newCmdCreate(f),
		newCmdUpdate(f),
		newCmdDelete(f),
		newCmdExport(f),
		newCmdBatchDelete(f),
	)

	return cmd
}
