package router

import (
	"net/url"

	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newCmdDelete(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a VPN router",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			_, err = client.Do("DELETE", "/api/invpn/router/"+args[0], q, nil)
			if err != nil {
				return err
			}

			cmdutil.WriteDeleted(f, "Router", args[0])
			return nil
		},
	}
}
