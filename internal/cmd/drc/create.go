package drc

import (
	"net/url"
	"os"

	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type createOptions struct {
	Name        string
	Model       string
	Content     string
	ContentFile string
	ContentType int
	Desc        string
	GroupIDs    []string
}

func newCmdCreate(f *factory.Factory) *cobra.Command {
	opts := &createOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a config template",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			content := opts.Content
			if opts.ContentFile != "" {
				data, err := os.ReadFile(opts.ContentFile)
				if err != nil {
					return err
				}
				content = string(data)
			}

			body := map[string]any{
				"name":    opts.Name,
				"model":   opts.Model,
				"content": content,
			}
			if opts.ContentType > 0 {
				body["contentType"] = opts.ContentType
			}
			if opts.Desc != "" {
				body["desc"] = opts.Desc
			}
			if len(opts.GroupIDs) > 0 {
				body["groupIds"] = opts.GroupIDs
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			respBody, err := client.Do("POST", "/api/drc", q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteCreated(f, "Config template", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Template name (required)")
	cmd.Flags().StringVar(&opts.Model, "model", "", "Device model (required)")
	cmd.Flags().StringVar(&opts.Content, "content", "", "Config content")
	cmd.Flags().StringVar(&opts.ContentFile, "content-file", "", "Read config content from file")
	cmd.Flags().IntVar(&opts.ContentType, "content-type", 0, "Content type (0=default)")
	cmd.Flags().StringVar(&opts.Desc, "desc", "", "Description")
	cmd.Flags().StringSliceVar(&opts.GroupIDs, "group-ids", nil, "Device group IDs")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("model")

	return cmd
}
