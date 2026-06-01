package router

import (
	"fmt"
	"net/url"

	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type transferOptions struct {
	TargetOid string
	NetworkID string
}

func newCmdTransfer(f *factory.Factory) *cobra.Command {
	opts := &transferOptions{}

	cmd := &cobra.Command{
		Use:   "transfer <router-id>",
		Short: "Transfer a router to another organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			oid, _ := cmd.Flags().GetString("oid")
			if oid == "" {
				return fmt.Errorf("--oid is required (source organization ID)")
			}

			q := url.Values{}
			q.Set("to", opts.TargetOid)
			if opts.NetworkID != "" {
				q.Set("networkId", opts.NetworkID)
			}

			path := "/api/invpn/orgs/" + oid + "/router/" + args[0] + "/transfer"
			respBody, err := client.Do("PUT", path, q, nil)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Router %s transferred to organization %s\n", args[0], opts.TargetOid)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.TargetOid, "to", "", "Target organization ID (required)")
	cmd.Flags().StringVar(&opts.NetworkID, "network-id", "", "Target network ID")
	_ = cmd.MarkFlagRequired("to")

	return cmd
}
