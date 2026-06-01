package mail

import (
	"fmt"

	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func NewCmdMail(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mail",
		Short: "Manage email notifications",
	}

	cmd.AddCommand(
		newCmdList(f),
		newCmdGet(f),
		newCmdCreate(f),
		newCmdVerify(f),
		newCmdCancel(f),
		newCmdRecords(f),
	)

	return cmd
}

func newCmdList(f *factory.Factory) *cobra.Command {
	lf := &cmdutil.ListFlags{}
	var title string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List email notifications",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewCursorQuery(cmd, lf)
			if title != "" {
				q.Set("title", title)
			}

			body, err := client.Get("/api/notify/mails", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	lf.RegisterPagination(cmd)
	cmd.Flags().StringVar(&title, "title", "", "Filter by title")

	return cmd
}

func newCmdGet(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get email notification details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body, err := client.Get("/api/notify/mails/"+args[0], nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}

type createMailOptions struct {
	Title     string
	Content   string
	OnlyAdmin bool
}

func newCmdCreate(f *factory.Factory) *cobra.Command {
	opts := &createMailOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create and send an email notification",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"title":     opts.Title,
				"content":   opts.Content,
				"onlyAdmin": opts.OnlyAdmin,
			}

			respBody, err := client.Post("/api/notify/mails", body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteCreated(f, "Email notification", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.Title, "title", "", "Email title (required)")
	cmd.Flags().StringVar(&opts.Content, "content", "", "Email content (HTML supported, required)")
	cmd.Flags().BoolVar(&opts.OnlyAdmin, "only-admin", false, "Send only to organization admins")
	_ = cmd.MarkFlagRequired("title")
	_ = cmd.MarkFlagRequired("content")

	return cmd
}

func newCmdVerify(f *factory.Factory) *cobra.Command {
	var title string
	var content string
	var address string
	var onlyAdmin bool

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Send a test/verification email",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"title":     title,
				"content":   content,
				"address":   address,
				"onlyAdmin": onlyAdmin,
			}

			respBody, err := client.Post("/api/notify/mails/verify", body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Verification email sent to %s\n", address)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Email title (required)")
	cmd.Flags().StringVar(&content, "content", "", "Email content (required)")
	cmd.Flags().StringVar(&address, "address", "", "Test email address (required)")
	cmd.Flags().BoolVar(&onlyAdmin, "only-admin", false, "Send only to admins")
	_ = cmd.MarkFlagRequired("title")
	_ = cmd.MarkFlagRequired("content")
	_ = cmd.MarkFlagRequired("address")

	return cmd
}

func newCmdCancel(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "cancel <id>",
		Short: "Cancel an in-progress email notification",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			respBody, err := client.Put("/api/notify/mails/"+args[0]+"/cancel", nil)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Email notification %s cancelled\n", args[0])
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}
}

func newCmdRecords(f *factory.Factory) *cobra.Command {
	lf := &cmdutil.ListFlags{}

	cmd := &cobra.Command{
		Use:   "records <id>",
		Short: "List recipients/records of an email notification",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewCursorQuery(cmd, lf)

			body, err := client.Get("/api/notify/mails/"+args[0]+"/records", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	lf.RegisterPagination(cmd)

	return cmd
}
