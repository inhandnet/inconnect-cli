package drc

import (
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type updateOptions struct {
	Name     string
	Desc     string
	GroupIDs []string
}

func newCmdUpdate(f *factory.Factory) *cobra.Command {
	opts := &updateOptions{}

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a config template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if opts.Name != "" {
				body["name"] = opts.Name
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

			respBody, err := client.Do("PUT", "/api/drc/"+args[0], q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteUpdated(f, "Config template", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Template name")
	cmd.Flags().StringVar(&opts.Desc, "desc", "", "Description")
	cmd.Flags().StringSliceVar(&opts.GroupIDs, "group-ids", nil, "Device group IDs")

	return cmd
}
