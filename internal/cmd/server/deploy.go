package server

import (
	"fmt"
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdDeploy(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "deploy <id>",
		Short: "Deploy or redeploy a VPN server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body, err := client.Get("/api/invpn/server/"+args[0]+"/deploy", url.Values{})
			if err != nil {
				if body != nil {
					_ = iostreams.FormatOutput(body, f.IO, f.IO.Output)
				}
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Server %s deployment triggered\n", args[0])
			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}
