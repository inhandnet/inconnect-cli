package billing

import (
	"fmt"
	"io"
	"os"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

type downloadReceiptOptions struct {
	Output string
	Title  string
}

func newCmdDownloadReceipt(f *factory.Factory) *cobra.Command {
	opts := &downloadReceiptOptions{}

	cmd := &cobra.Command{
		Use:   "download-receipt <payment-id>",
		Short: "Download a payment receipt PDF",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			title := opts.Title
			if title == "" {
				title = "receipt"
			}

			path := fmt.Sprintf("/api/billing/payment/%s/export/%s.pdf", args[0], title)
			data, err := client.Get(path, nil)
			if err != nil {
				return err
			}

			var w io.Writer
			if opts.Output == "" || opts.Output == "-" {
				w = os.Stdout
			} else {
				file, err := os.Create(opts.Output)
				if err != nil {
					return err
				}
				defer file.Close()
				w = file
			}

			_, err = w.Write(data)
			if err != nil {
				return err
			}

			if opts.Output != "" && opts.Output != "-" {
				fmt.Fprintf(f.IO.ErrOut, "Receipt saved to %s\n", opts.Output)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.Output, "file", "", "Output file path (default: stdout)")
	cmd.Flags().StringVar(&opts.Title, "title", "receipt", "Receipt title for filename")

	return cmd
}
