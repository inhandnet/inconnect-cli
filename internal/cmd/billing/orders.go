package billing

import (
	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type listOrdersOptions struct {
	cmdutil.ListFlags
	Number   string
	Email    string
	After    string
	Before   string
	Type     string
	Status   string
	Invoiced string
}

func newCmdListOrders(f *factory.Factory) *cobra.Command {
	opts := &listOrdersOptions{}

	cmd := &cobra.Command{
		Use:     "list-orders",
		Aliases: []string{"orders"},
		Short:   "List billing orders/transactions",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewCursorQuery(cmd, &opts.ListFlags)
			if opts.Number != "" {
				q.Set("number", opts.Number)
			}
			if opts.Email != "" {
				q.Set("email", opts.Email)
			}
			if opts.After != "" {
				q.Set("begin", cmdutil.ParseTimeFlag(opts.After))
			}
			if opts.Before != "" {
				q.Set("end", cmdutil.ParseTimeFlag(opts.Before))
			}
			if opts.Type != "" {
				q.Set("type", opts.Type)
			}
			if opts.Status != "" {
				q.Set("status", opts.Status)
			}
			if opts.Invoiced != "" {
				q.Set("invoiced", opts.Invoiced)
			}

			body, err := client.Get("/api/billing/transaction", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.RegisterPagination(cmd)
	cmd.Flags().StringVar(&opts.Number, "number", "", "Filter by order number")
	cmd.Flags().StringVar(&opts.Email, "email", "", "Filter by email")
	cmd.Flags().StringVar(&opts.After, "after", "", "Start time (e.g. 2024-01-01, 2024-01-01T00:00:00Z)")
	cmd.Flags().StringVar(&opts.Before, "before", "", "End time (e.g. 2024-12-31, 2024-12-31T23:59:59Z)")
	cmd.Flags().StringVar(&opts.Type, "type", "", "Transaction type (default: payment)")
	cmd.Flags().StringVar(&opts.Status, "status", "", "Transaction status (default: succeeded)")
	cmd.Flags().StringVar(&opts.Invoiced, "invoiced", "", "Filter by invoice status (true/false)")

	return cmd
}
