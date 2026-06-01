package router

import (
	"net/url"

	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdStats(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Get device statistics (online/offline counts)",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body, err := client.Get("/api/devices/stats", url.Values{})
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}

func newCmdModels(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "models",
		Short: "List supported router models",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body, err := client.Get("/api/invpn/routers/models", url.Values{})
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}
