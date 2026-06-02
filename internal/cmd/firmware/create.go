package firmware

import (
	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type createOptions struct {
	FID        string
	Name       string
	Version    string
	Model      string
	Desc       string
	JobTimeout int64
}

func newCmdCreate(f *factory.Factory) *cobra.Command {
	opts := &createOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a firmware package",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"fid":     opts.FID,
				"name":    opts.Name,
				"version": opts.Version,
				"model":   opts.Model,
			}
			if opts.Desc != "" {
				body["desc"] = opts.Desc
			}
			if opts.JobTimeout > 0 {
				body["jobTimeout"] = opts.JobTimeout
			}

			respBody, err := client.Post("/api/firmware", body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteCreated(f, "Firmware", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.FID, "fid", "", "File ID from upload (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Firmware name (required)")
	cmd.Flags().StringVar(&opts.Version, "version", "", "Firmware version (required)")
	cmd.Flags().StringVar(&opts.Model, "model", "", "Device model (required)")
	cmd.Flags().StringVar(&opts.Desc, "desc", "", "Description")
	cmd.Flags().Int64Var(&opts.JobTimeout, "job-timeout", 0, "Job timeout in minutes (1-60)")
	_ = cmd.MarkFlagRequired("fid")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("version")
	_ = cmd.MarkFlagRequired("model")

	return cmd
}
