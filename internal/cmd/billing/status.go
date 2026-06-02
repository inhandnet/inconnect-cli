package billing

import (
	"fmt"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type updateStatusOptions struct {
	Free bool
	Type string
}

func newCmdUpdateStatus(f *factory.Factory) *cobra.Command {
	opts := &updateStatusOptions{}

	cmd := &cobra.Command{
		Use:   "update-status <org-id>",
		Short: "Update billing status for an organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"free": opts.Free,
			}
			if opts.Type != "" {
				body["type"] = opts.Type
			}

			respBody, err := client.Put("/api/billing/account/"+args[0], body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Billing status updated for org %s\n", args[0])
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().BoolVar(&opts.Free, "free", false, "Set organization as free (no billing)")
	cmd.Flags().StringVar(&opts.Type, "type", "", "Account type")

	return cmd
}
