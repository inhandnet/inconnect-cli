package firmware

import (
	"net/url"

	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdGet(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a firmware package by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body, err := client.Get("/api/firmware/"+args[0], url.Values{})
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}
