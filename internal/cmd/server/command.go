package server

import (
	"fmt"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type commandOptions struct {
	Command string
}

func newCmdCommand(f *factory.Factory) *cobra.Command {
	opts := &commandOptions{}

	cmd := &cobra.Command{
		Use:   "command <oid>",
		Short: "Send a command to the VPN server of an organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"command": opts.Command,
			}

			respBody, err := client.Post("/api/invpn/server/"+args[0]+"/command", body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Command %q sent to server of org %s\n", opts.Command, args[0])
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.Command, "cmd", "", "Command to send (e.g. hold_release)")
	_ = cmd.MarkFlagRequired("cmd")

	return cmd
}
