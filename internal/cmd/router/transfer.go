package router

import (
	"fmt"
	"net/url"

	"github.com/tidwall/gjson"

	"github.com/inhandnet/inconnect-cli/internal/api"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
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
				// error_code 10001 (internal_error) can be returned after the router
				// has already moved to the target org (a later step failed). Re-query
				// to confirm the actual ownership before reporting a failure.
				if api.ErrorCode(err) == "10001" {
					verify := url.Values{}
					verify.Set("oid", opts.TargetOid)
					if rb, gerr := client.Get("/api/invpn/router/"+args[0], verify); gerr == nil {
						if gjson.GetBytes(rb, "result.oid").String() == opts.TargetOid {
							fmt.Fprintf(f.IO.ErrOut, "%s Router %s is now in organization %s, but a post-transfer step failed. Verify its network assignment with 'inconnect router get %s'.\n",
								iostreams.Yellow("!"), args[0], opts.TargetOid, args[0])
							return iostreams.FormatOutput(rb, f.IO, f.IO.Output)
						}
					}
				}
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
