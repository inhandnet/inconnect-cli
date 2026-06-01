package router

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

func newCmdExec(f *factory.Factory) *cobra.Command {
	var timeoutSec int

	cmd := &cobra.Command{
		Use:   "exec <id> <command>...",
		Short: "Run a shell command on a router remotely",
		Long: `Run a shell command on a router remotely and print its output.

Examples:
  ics router exec <id> show log
  ics router exec <id> "ifconfig eth0"`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			command := strings.Join(args[1:], " ")

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			body := map[string]any{
				"name":     "remote console: " + command,
				"type":     2,
				"objectId": args[0],
				"priority": 30,
				"timeout":  timeoutSec * 1000,
				"data": map[string]any{
					"deviceDesc":    "CMD",
					"sensorId":      0,
					"deviceType":    0,
					"deviceContent": command,
				},
			}

			respBody, err := client.Do("POST", "/api2/tasks/run", q, body)
			if err != nil {
				if respBody != nil {
					_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				}
				return err
			}

			if taskErr := gjson.GetBytes(respBody, "result.error").String(); taskErr != "" {
				return fmt.Errorf("%s", taskErr)
			}

			fmt.Fprintln(f.IO.Out, gjson.GetBytes(respBody, "result.data.response").String())
			return nil
		},
	}

	cmd.Flags().IntVar(&timeoutSec, "timeout", 30, "Command timeout in seconds")

	return cmd
}

func newCmdRunningConfig(f *factory.Factory) *cobra.Command {
	var timeoutSec int

	cmd := &cobra.Command{
		Use:   "running-config <id>",
		Short: "Fetch the config currently LIVE on the device (decoded)",
		Long: `Fetch the configuration currently running on the device right now (live,
server-side decoded). Requires the device to be online.

Related commands that are easy to confuse:
  - "router device-config get" — the platform's STORED copy (may differ from live).
  - "router config-send"        — push the platform-rendered VPN-only config.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			// GET /config/raw triggers a GET_RUNNING_CONFIG task on the api-server
			// and returns the server-side decoded config, handling device online
			// check and content-type decoding (hex/etc) that we can't do client-side.
			q := url.Values{}
			q.Set("timeout", strconv.Itoa(timeoutSec))
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			respBody, err := client.Get("/api/devices/"+args[0]+"/config/raw", q)
			if err != nil {
				if respBody != nil {
					_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				}
				return err
			}

			// The endpoint returns the config as a bare JSON string.
			fmt.Fprintln(f.IO.Out, gjson.ParseBytes(respBody).String())
			return nil
		},
	}

	cmd.Flags().IntVar(&timeoutSec, "timeout", 60, "Task timeout in seconds")

	return cmd
}
