package server

import (
	"net/url"
	"os"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdList(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List VPN servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			} else if oid := os.Getenv("INCONNECT_OID"); oid != "" {
				q.Set("oid", oid)
			}

			path := "/api/invpn/servers"
			if oid := q.Get("oid"); oid != "" {
				path = "/api/invpn/org/" + oid + "/servers"
				q.Del("oid")
			}

			body, err := client.Get(path, q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(redactBody(cmd, body), f.IO, f.IO.Output)
		},
	}
}
