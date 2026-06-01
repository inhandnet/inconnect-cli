package router

import (
	"fmt"
	"net/url"

	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdOvpn(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "ovpn <id>",
		Short: "Download router OpenVPN config",
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

			body, err := client.Get("/api/invpn/router/"+args[0]+"/router.ovpn", q)
			if err != nil {
				return err
			}

			fmt.Fprint(f.IO.Out, string(body))
			return nil
		},
	}
}

func newCmdClientOvpn(f *factory.Factory) *cobra.Command {
	var uid, compLzo string

	cmd := &cobra.Command{
		Use:   "client-ovpn",
		Short: "Download OpenVPN client config",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}
			if uid != "" {
				q.Set("uid", uid)
			}
			if compLzo != "" {
				q.Set("comp-lzo", compLzo)
			}

			body, err := client.Get("/api/invpn/client.ovpn", q)
			if err != nil {
				if body != nil {
					_ = iostreams.FormatOutput(body, f.IO, f.IO.Output)
				}
				return err
			}

			fmt.Fprint(f.IO.Out, string(body))
			return nil
		},
	}

	cmd.Flags().StringVar(&uid, "uid", "", "User ID")
	cmd.Flags().StringVar(&compLzo, "comp-lzo", "", "LZO compression (yes/no)")

	return cmd
}
