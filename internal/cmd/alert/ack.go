package alert

import (
	"fmt"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdAck(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "ack <id> [id...]",
		Short: "Acknowledge alerts",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]interface{}{
				"resourceIds": args,
			}

			respBody, err := client.Put("/api/invpn/alerts/acknowledge", body)
			if err != nil {
				if respBody != nil {
					_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				}
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Acknowledged %d alert(s)\n", len(args))
			return nil
		},
	}
}
