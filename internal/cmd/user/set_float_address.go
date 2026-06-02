package user

import (
	"fmt"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdSetFloatAddress(f *factory.Factory) *cobra.Command {
	var enable bool

	cmd := &cobra.Command{
		Use:   "set-float-address <id>",
		Short: "Enable or disable floating IP for a user",
		Long:  "Enable or disable floating IP for a user.\n\n<id> is the VPN user ID (the 'id' field from 'user list', not 'uid').",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"floatAddress": enable,
			}

			respBody, err := client.Put("/api/invpn/users/"+args[0]+"/setFloatAddress", body)
			if err != nil {
				if respBody != nil {
					_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				}
				return err
			}

			action := "disabled"
			if enable {
				action = "enabled"
			}
			fmt.Fprintf(f.IO.ErrOut, "Float address %s for user %s\n", action, args[0])
			return nil
		},
	}

	cmd.Flags().BoolVar(&enable, "enable", true, "Enable or disable float address")

	return cmd
}
