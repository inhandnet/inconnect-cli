package datausage

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func NewCmdDataUsage(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "data-usage",
		Aliases: []string{"usage"},
		Short:   "View VPN data usage statistics",
	}

	cmd.AddCommand(
		newCmdSummary(f),
		newCmdRouter(f),
		newCmdAccount(f),
		newCmdRouterMonth(f),
		newCmdAccountMonth(f),
		newCmdRouterExport(f),
		newCmdAccountExport(f),
	)

	return cmd
}

func newCmdSummary(f *factory.Factory) *cobra.Command {
	var month, date string

	cmd := &cobra.Command{
		Use:   "summary",
		Short: "Get organization-level data usage summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}
			if month != "" {
				q.Set("month", month)
			}
			if date != "" {
				q.Set("date", date)
			}

			body, err := client.Get("/api/invpn/data-usage", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&month, "month", "", "Month (YYYY-MM)")
	cmd.Flags().StringVar(&date, "date", "", "Date (YYYY-MM-DD)")

	return cmd
}

func newCmdRouter(f *factory.Factory) *cobra.Command {
	var date, id string

	cmd := &cobra.Command{
		Use:   "router",
		Short: "Get daily router data usage details",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}
			q.Set("date", date)
			if id != "" {
				q.Set("id", id)
			}

			body, err := client.Get("/api/invpn/data-usage/router", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&date, "date", "", "Date (required, YYYY-MM-DD)")
	cmd.Flags().StringVar(&id, "id", "", "Router ID filter")
	_ = cmd.MarkFlagRequired("date")

	return cmd
}

func newCmdAccount(f *factory.Factory) *cobra.Command {
	var date, id string

	cmd := &cobra.Command{
		Use:   "account",
		Short: "Get daily account data usage details",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}
			q.Set("date", date)
			if id != "" {
				q.Set("id", id)
			}

			body, err := client.Get("/api/invpn/data-usage/account", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&date, "date", "", "Date (required, YYYY-MM-DD)")
	cmd.Flags().StringVar(&id, "id", "", "Account ID filter")
	_ = cmd.MarkFlagRequired("date")

	return cmd
}

func newCmdRouterMonth(f *factory.Factory) *cobra.Command {
	var month string

	cmd := &cobra.Command{
		Use:   "router-month",
		Short: "Get monthly router data usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}
			q.Set("month", month)

			body, err := client.Get("/api/invpn/data-usage/router/month", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&month, "month", "", "Month (required, YYYY-MM)")
	_ = cmd.MarkFlagRequired("month")

	return cmd
}

func newCmdAccountMonth(f *factory.Factory) *cobra.Command {
	var month string

	cmd := &cobra.Command{
		Use:   "account-month",
		Short: "Get monthly account data usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}
			q.Set("month", month)

			body, err := client.Get("/api/invpn/data-usage/account/month", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&month, "month", "", "Month (required, YYYY-MM)")
	_ = cmd.MarkFlagRequired("month")

	return cmd
}

func newCmdRouterExport(f *factory.Factory) *cobra.Command {
	var id, date, output string
	var language int

	cmd := &cobra.Command{
		Use:   "router-export",
		Short: "Export router data usage to Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}
			if id != "" {
				q.Set("id", id)
			}
			if date != "" {
				q.Set("date", date)
			}
			q.Set("language", strconv.Itoa(language))

			data, err := client.Get("/api/invpn/data-usage/router/export", q)
			if err != nil {
				return err
			}

			return writeExport(f, data, output)
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "Router ID filter")
	cmd.Flags().StringVar(&date, "date", "", "Date filter (YYYY-MM-DD)")
	cmd.Flags().StringVar(&output, "file", "", "Output file path (default: stdout)")
	cmd.Flags().IntVar(&language, "language", 2, "Export language (1=English, 2=Chinese)")

	return cmd
}

func newCmdAccountExport(f *factory.Factory) *cobra.Command {
	var id, date, output string
	var language int

	cmd := &cobra.Command{
		Use:   "account-export",
		Short: "Export account data usage to Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}
			if id != "" {
				q.Set("id", id)
			}
			if date != "" {
				q.Set("date", date)
			}
			q.Set("language", strconv.Itoa(language))

			data, err := client.Get("/api/invpn/data-usage/account/export", q)
			if err != nil {
				return err
			}

			return writeExport(f, data, output)
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "Account ID filter")
	cmd.Flags().StringVar(&date, "date", "", "Date filter (YYYY-MM-DD)")
	cmd.Flags().StringVar(&output, "file", "", "Output file path (default: stdout)")
	cmd.Flags().IntVar(&language, "language", 2, "Export language (1=English, 2=Chinese)")

	return cmd
}

func writeExport(f *factory.Factory, data []byte, output string) error {
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

	_, err := w.Write(data)
	if err != nil {
		return err
	}

	if output != "" && output != "-" {
		fmt.Fprintf(f.IO.ErrOut, "Exported to %s\n", output)
	}
	return nil
}
