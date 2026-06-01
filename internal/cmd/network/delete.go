package network

import (
	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newCmdDelete(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a VPN network",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			_, err = client.Delete("/api/invpn/networks/" + args[0])
			if err != nil {
				return err
			}

			cmdutil.WriteDeleted(f, "Network", args[0])
			return nil
		},
	}
}
