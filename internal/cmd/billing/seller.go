package billing

import (
	"fmt"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdGetSeller(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get-seller",
		Short: "Get order notification email settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body, err := client.Get("/api/charge/seller", nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}

type updateSellerOptions struct {
	Email string
}

func newCmdUpdateSeller(f *factory.Factory) *cobra.Command {
	opts := &updateSellerOptions{}

	cmd := &cobra.Command{
		Use:   "update-seller",
		Short: "Update order notification email settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if opts.Email != "" {
				body["email"] = opts.Email
			}

			respBody, err := client.Put("/api/charge/seller", body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Seller notification settings updated\n")
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.Email, "email", "", "Notification email address")

	return cmd
}
