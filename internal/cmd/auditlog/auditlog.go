package auditlog

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"

	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func NewCmdAuditLog(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "audit-log",
		Aliases: []string{"log"},
		Short:   "View and export audit logs",
	}

	cmd.AddCommand(
		newCmdList(f),
		newCmdExport(f),
	)

	return cmd
}

type listOptions struct {
	cmdutil.ListFlags
	After    string
	Before   string
	Levels   []int
	Username string
}

func newCmdList(f *factory.Factory) *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List audit/behavior logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewCursorQuery(cmd, &opts.ListFlags)
			q.Set("start_time", cmdutil.ParseTimeFlag(opts.After))
			q.Set("end_time", cmdutil.ParseTimeFlag(opts.Before))
			if opts.Username != "" {
				q.Set("username", opts.Username)
			}
			for _, l := range opts.Levels {
				q.Add("level", strconv.Itoa(l))
			}

			body, err := client.Get("/api2/behav_log", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.RegisterPagination(cmd)
	cmd.Flags().StringVar(&opts.After, "after", "", "Start time (required, e.g. 2024-01-01, 2024-01-01T00:00:00Z)")
	cmd.Flags().StringVar(&opts.Before, "before", "", "End time (required, e.g. 2024-12-31, 2024-12-31T23:59:59Z)")
	cmd.Flags().IntSliceVar(&opts.Levels, "level", nil, "Log level filter (can specify multiple)")
	cmd.Flags().StringVar(&opts.Username, "username", "", "Filter by username")
	_ = cmd.MarkFlagRequired("after")
	_ = cmd.MarkFlagRequired("before")

	return cmd
}

type exportOptions struct {
	Output   string
	After    string
	Before   string
	Level    int
	Language int
}

func newCmdExport(f *factory.Factory) *cobra.Command {
	opts := &exportOptions{}

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export audit logs to XLS file",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if opts.After != "" {
				q.Set("start_time", cmdutil.ParseTimeFlag(opts.After))
			}
			if opts.Before != "" {
				q.Set("end_time", cmdutil.ParseTimeFlag(opts.Before))
			}
			if opts.Level > 0 {
				q.Set("level", strconv.Itoa(opts.Level))
			}
			if opts.Language > 0 {
				q.Set("language", strconv.Itoa(opts.Language))
			}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			data, err := client.Get("/api/audit/log/export", q)
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
	cmd.Flags().StringVar(&opts.After, "after", "", "Start time (e.g. 2024-01-01, 2024-01-01T00:00:00Z)")
	cmd.Flags().StringVar(&opts.Before, "before", "", "End time (e.g. 2024-12-31, 2024-12-31T23:59:59Z)")
	cmd.Flags().IntVar(&opts.Level, "level", 0, "Log level filter")
	cmd.Flags().IntVar(&opts.Language, "language", 2, "Export language (1=English, 2=Chinese)")

	return cmd
}
