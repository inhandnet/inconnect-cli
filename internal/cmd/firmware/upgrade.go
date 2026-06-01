package firmware

import (
	"fmt"

	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type upgradeOptions struct {
	FirmwareID string
	Timeout    int
}

func newCmdUpgrade(f *factory.Factory) *cobra.Command {
	opts := &upgradeOptions{}

	cmd := &cobra.Command{
		Use:   "upgrade <device-id>",
		Short: "Upgrade firmware on a device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"firmwareId": opts.FirmwareID,
				"timeout":    opts.Timeout,
			}

			respBody, err := client.Post("/api/device/"+args[0]+"/upgrade", body)
			if err != nil {
				if respBody != nil {
					_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				}
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Firmware upgrade started for device %s\n", args[0])
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.FirmwareID, "firmware-id", "", "Firmware package ID (required)")
	cmd.Flags().IntVar(&opts.Timeout, "timeout", 600, "Upgrade timeout in seconds")
	_ = cmd.MarkFlagRequired("firmware-id")

	return cmd
}
