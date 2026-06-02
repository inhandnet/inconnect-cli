package router

import (
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdSubnet(f *factory.Factory) *cobra.Command {
	var networkID string

	cmd := &cobra.Command{
		Use:   "next-subnet",
		Short: "Get the next available subnet for a router",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if networkID != "" {
				q.Set("networkId", networkID)
			}

			body, err := client.Get("/api/invpn/router/subnet", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&networkID, "network-id", "", "Network ID")

	return cmd
}
