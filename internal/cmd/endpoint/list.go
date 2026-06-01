package endpoint

import (
	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type listOptions struct {
	cmdutil.ListFlags
	RouterID string
	VIP      string
	RIP      string
}

func newCmdList(f *factory.Factory) *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List endpoints",
		Long:    "List endpoints for a specific router (--router-id) or all endpoints.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, &opts.ListFlags)
			if opts.VIP != "" {
				q.Set("vip", opts.VIP)
			}
			if opts.RIP != "" {
				q.Set("rip", opts.RIP)
			}

			path := "/api/invpn/routers/all/endpoints"
			if opts.RouterID != "" {
				path = "/api/invpn/routers/" + opts.RouterID + "/endpoints"
			}

			body, err := client.Get(path, q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&opts.RouterID, "router-id", "", "Filter by router ID")
	cmd.Flags().StringVar(&opts.VIP, "vip", "", "Filter by virtual IP")
	cmd.Flags().StringVar(&opts.RIP, "rip", "", "Filter by real IP")

	return cmd
}
