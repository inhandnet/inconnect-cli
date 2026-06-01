package router

import (
	"fmt"
	"net/url"

	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newCmdKick(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "kick <id>",
		Short: "Force-disconnect a device",
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

			_, err = client.Get("/api/device/"+args[0]+"/kick", q)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Device %s kicked\n", args[0])
			return nil
		},
	}
}

func newCmdReboot(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "reboot <id>",
		Short: "Reboot a device",
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

			// timeout is the server-side wait (ms) for the device to ack, per
			// DeferredResult; the web UI uses 15000.
			body := map[string]any{
				"method":  "reboot",
				"timeout": 15000,
			}

			_, err = client.Do("POST", "/api/device/"+args[0]+"/methods", q, body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Reboot command sent to device %s\n", args[0])
			return nil
		},
	}
}
