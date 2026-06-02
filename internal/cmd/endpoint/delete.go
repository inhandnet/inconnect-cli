package endpoint

import (
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newCmdDelete(f *factory.Factory) *cobra.Command {
	var routerID string

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an endpoint",
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

			_, err = client.Do("DELETE", "/api/invpn/router/"+routerID+"/endpoint/"+args[0], q, nil)
			if err != nil {
				return err
			}

			cmdutil.WriteDeleted(f, "Endpoint", args[0])
			return nil
		},
	}

	cmd.Flags().StringVar(&routerID, "router-id", "", "Router ID (required)")
	_ = cmd.MarkFlagRequired("router-id")

	return cmd
}
