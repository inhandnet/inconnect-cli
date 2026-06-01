package org

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/spf13/cobra"
)

type exportOptions struct {
	Output   string
	Email    string
	Name     string
	Creator  string
	Status   string
	Language int
}

func newCmdExport(f *factory.Factory) *cobra.Command {
	opts := &exportOptions{}

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export organizations to XLSX file",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if opts.Email != "" {
				q.Set("email", opts.Email)
			}
			if opts.Name != "" {
				q.Set("name", opts.Name)
			}
			if opts.Creator != "" {
				q.Set("creator", opts.Creator)
			}
			if opts.Status != "" {
				q.Set("status", opts.Status)
			}
			if opts.Language > 0 {
				q.Set("language", fmt.Sprintf("%d", opts.Language))
			}

			data, err := client.Get("/api/invpn/org/export", q)
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
				fmt.Fprintf(f.IO.ErrOut, "Exported to %s\n", opts.Output)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.Output, "file", "", "Output file path (default: stdout)")
	cmd.Flags().StringVar(&opts.Email, "email", "", "Filter by email")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Filter by name")
	cmd.Flags().StringVar(&opts.Creator, "creator", "", "Filter by creator")
	cmd.Flags().StringVar(&opts.Status, "status", "", "Filter by status")
	cmd.Flags().IntVar(&opts.Language, "language", 2, "Export language (1=English, 2=Chinese)")

	return cmd
}
