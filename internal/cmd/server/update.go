package server

import (
	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type updateOptions struct {
	Subnet          string
	SecondIncrement int
	ThirdIncrement  int
}

func newCmdUpdate(f *factory.Factory) *cobra.Command {
	opts := &updateOptions{}

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a VPN server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if opts.Subnet != "" {
				body["subnet"] = opts.Subnet
			}
			if cmd.Flags().Changed("second-increment") {
				body["secondIncrement"] = opts.SecondIncrement
			}
			if cmd.Flags().Changed("third-increment") {
				body["thirdIncrement"] = opts.ThirdIncrement
			}

			respBody, err := client.Put("/api/invpn/server/"+args[0], body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteUpdated(f, "Server", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.Subnet, "subnet", "", "Server subnet (e.g. 10.8.0.0/16)")
	cmd.Flags().IntVar(&opts.SecondIncrement, "second-increment", 0, "Second octet increment for IP allocation")
	cmd.Flags().IntVar(&opts.ThirdIncrement, "third-increment", 0, "Third octet increment for IP allocation")

	return cmd
}
