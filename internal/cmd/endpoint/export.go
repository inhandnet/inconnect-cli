package endpoint

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newCmdExport(f *factory.Factory) *cobra.Command {
	var name, routerID, output string
	var language int

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export endpoints to Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}
			if name != "" {
				q.Set("name", name)
			}
			if routerID != "" {
				q.Set("rid", routerID)
			}
			q.Set("language", strconv.Itoa(language))

			data, err := client.Get("/api/invpn/endpoint/export", q)
			if err != nil {
				return err
			}

			var w io.Writer
			if output == "" || output == "-" {
				w = os.Stdout
			} else {
				file, err := os.Create(output)
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

			if output != "" && output != "-" {
				fmt.Fprintf(f.IO.ErrOut, "Exported to %s\n", output)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Filter by name")
	cmd.Flags().StringVar(&routerID, "router-id", "", "Filter by router ID")
	cmd.Flags().StringVar(&output, "file", "", "Output file path (default: stdout)")
	cmd.Flags().IntVar(&language, "language", 2, "Export language (1=English, 2=Chinese)")

	return cmd
}

func newCmdBatchDelete(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "batch-delete <endpoint-id> [endpoint-id...]",
		Short: "Batch delete endpoints",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"resourceIds": args,
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			_, err = client.Do("POST", "/api/invpn/router/endpoints/batch-delete", q, body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "%d endpoints deleted\n", len(args))
			return nil
		},
	}
}
