package banner

import (
	"fmt"

	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func NewCmdBanner(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "banner",
		Short: "Manage system banner messages",
	}

	cmd.AddCommand(
		newCmdList(f),
		newCmdCurrent(f),
		newCmdCreate(f),
		newCmdRevoke(f),
	)

	return cmd
}

func newCmdList(f *factory.Factory) *cobra.Command {
	lf := &cmdutil.ListFlags{}
	var keyword string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List banner messages",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewCursorQuery(cmd, lf)
			if keyword != "" {
				q.Set("keyword", keyword)
			}

			body, err := client.Get("/api/banners", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	lf.RegisterPagination(cmd)
	cmd.Flags().StringVar(&keyword, "keyword", "", "Search keyword")

	return cmd
}

func newCmdCurrent(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Get current active banner",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body, err := client.Get("/api/banners/current", nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}

type createBannerOptions struct {
	Content string
	EndTime string
}

func newCmdCreate(f *factory.Factory) *cobra.Command {
	opts := &createBannerOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a banner message",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"content": opts.Content,
			}
			if opts.EndTime != "" {
				body["endTime"] = cmdutil.ParseTimeFlag(opts.EndTime)
			}

			respBody, err := client.Post("/api/banners", body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteCreated(f, "Banner", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.Content, "content", "", "Banner content (required, max 300 chars)")
	cmd.Flags().StringVar(&opts.EndTime, "end-time", "", "Expiration time (required, e.g. 2026-12-31 or 2026-12-31T23:59:59)")
	_ = cmd.MarkFlagRequired("content")
	_ = cmd.MarkFlagRequired("end-time")

	return cmd
}

func newCmdRevoke(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "revoke <id>",
		Short: "Revoke a banner message",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			respBody, err := client.Put("/api/banners/"+args[0]+"/revoke", nil)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Banner %s revoked\n", args[0])
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}
}
