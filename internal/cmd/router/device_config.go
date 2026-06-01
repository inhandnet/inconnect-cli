package router

import (
	"net/url"
	"os"

	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdDeviceConfig(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "device-config",
		Short: "Manage the full device config stored on the platform",
		Long: `Get or push the full device configuration tracked by the platform.

Related commands that are easy to confuse:
  - "router running-config" — the config currently LIVE on the device (read).
  - "router config-send"     — push the platform-rendered VPN-only config.`,
	}

	cmd.AddCommand(
		newCmdDeviceConfigGet(f),
		newCmdDeviceConfigSet(f),
		newCmdDeviceConfigExport(f),
	)

	return cmd
}

func newCmdDeviceConfigGet(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <device-id>",
		Short: "Get the config the platform has STORED for a device",
		Long: `Get the device configuration stored on the platform (the platform's copy,
with a version number) — not necessarily what the device is running right now.

To fetch the config currently LIVE on the device, use "router running-config".`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			body, err := client.Get("/api/devices/"+args[0]+"/config", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}

func newCmdDeviceConfigSet(f *factory.Factory) *cobra.Command {
	var content, contentFile, desc string

	cmd := &cobra.Command{
		Use:   "set <device-id>",
		Short: "Push a full device config YOU supply to a device",
		Long: `Push a full device configuration that you provide (via --content or
--content-file) to the device and apply it. Requires the device to be online.

You supply the entire config content here. To push only the platform-rendered
VPN config (certs/CA/firewall) without supplying content, use
"router config-send" instead.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			cfg := content
			if contentFile != "" {
				data, readErr := os.ReadFile(contentFile)
				if readErr != nil {
					return readErr
				}
				cfg = string(data)
			}

			body := map[string]any{
				"deviceType":    0,
				"deviceContent": cfg,
			}
			if desc != "" {
				body["deviceDesc"] = desc
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			respBody, err := client.Do("POST", "/api/devices/"+args[0]+"/config/set", q, body)
			if err != nil {
				if respBody != nil {
					_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				}
				return err
			}

			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&content, "content", "", "Config content string")
	cmd.Flags().StringVar(&contentFile, "content-file", "", "Read config content from file")
	cmd.Flags().StringVar(&desc, "desc", "", "Config description")

	return cmd
}

func newCmdDeviceConfigExport(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "export <device-id>",
		Short: "Export device running config metadata",
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

			body, err := client.Get("/api/devices/"+args[0]+"/config/export", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}

