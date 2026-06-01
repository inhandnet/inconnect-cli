package billing

import (
	"fmt"

	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type updateInvoiceOptions struct {
	Action             string
	ChannelOrderNumber string
	InvoiceNumber      string
}

func newCmdUpdateInvoice(f *factory.Factory) *cobra.Command {
	opts := &updateInvoiceOptions{}

	cmd := &cobra.Command{
		Use:   "update-invoice <order-id>",
		Short: "Update invoice status for an order",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"action": opts.Action,
			}
			if opts.ChannelOrderNumber != "" {
				body["channelOrderNumber"] = opts.ChannelOrderNumber
			}
			if opts.InvoiceNumber != "" {
				body["invoiceNumber"] = opts.InvoiceNumber
			}

			respBody, err := client.Put("/api/billing/transaction/"+args[0]+"/invoice", body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Invoice status updated for order %s\n", args[0])
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.Action, "action", "", "Action: invoice or revert (required)")
	cmd.Flags().StringVar(&opts.ChannelOrderNumber, "channel-order-number", "", "Channel order number")
	cmd.Flags().StringVar(&opts.InvoiceNumber, "invoice-number", "", "Invoice number")
	_ = cmd.MarkFlagRequired("action")

	return cmd
}
