package billing

import (
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewCmdBilling(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "billing",
		Short: "Manage billing, orders, and payments",
	}

	cmd.AddCommand(
		newCmdUpdateStatus(f),
		newCmdListOrders(f),
		newCmdDownloadReceipt(f),
		newCmdUpdateInvoice(f),
		newCmdGetSeller(f),
		newCmdUpdateSeller(f),
	)

	return cmd
}
