package router

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

func newCmdSSH(f *factory.Factory) *cobra.Command {
	var (
		server     string
		timeoutSec int
		port       int
	)

	cmd := &cobra.Command{
		Use:   "ssh <id>",
		Short: "Print an SSH command to log into a device via ngrok",
		Long: `Open a TCP ngrok tunnel to the device's SSH port and print a ready-to-use
ssh command. The ngrok server must have the embedded SSH reverse proxy
enabled; if it isn't, the tunnel URL won't carry a tunnel-id and this command
will fail.`,
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

			fmt.Fprintf(f.IO.ErrOut, "Establishing ngrok tunnel to device %s...\n", args[0])

			body := map[string]any{
				"name":     "ngrok connect",
				"type":     23,
				"objectId": args[0],
				"priority": 30,
				"timeout":  timeoutSec * 1000,
				"data": map[string]any{
					"server": server,
					"proto":  "tcp",
					"port":   port,
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

			tunnelURL := gjson.GetBytes(respBody, "result.data.response").String()
			if tunnelURL == "" {
				return fmt.Errorf("no tunnel URL returned")
			}

			u, err := url.Parse(tunnelURL)
			if err != nil {
				return fmt.Errorf("parse tunnel URL %q: %w", tunnelURL, err)
			}
			// The ngrok server (INC-1311) returns SSH tunnels as
			//   tcp://<uuid>.<domain>:<sshport>
			// where <uuid> is the routing token (used as the SSH username
			// prefix) and <sshport> is the embedded SSH proxy's public port.
			// The tunnel's internal dynamic TCP port is never exposed and
			// does not appear in the URL.
			tunnelID, host := splitFirstLabel(u.Hostname())
			if !isHexUUID(tunnelID) {
				return fmt.Errorf("tunnel URL leading label %q is not a tunnel uuid; ngrok SSH proxy not enabled on this deployment", tunnelID)
			}
			sshPort, perr := strconv.Atoi(u.Port())
			if perr != nil || sshPort <= 0 {
				return fmt.Errorf("tunnel URL ssh port %q is not a positive integer", u.Port())
			}
			fmt.Fprintln(f.IO.Out, "Tunnel ready:")
			fmt.Fprintf(f.IO.Out, "  Host:           %s\n", host)
			fmt.Fprintf(f.IO.Out, "  SSH proxy port: %d\n", sshPort)
			fmt.Fprintf(f.IO.Out, "  Tunnel ID:      %s\n", tunnelID)
			if sshPort == 22 {
				fmt.Fprintf(f.IO.Out, "  Connect:        ssh %s+<device-user>@%s\n", tunnelID, host)
			} else {
				fmt.Fprintf(f.IO.Out, "  Connect:        ssh -p %d %s+<device-user>@%s\n", sshPort, tunnelID, host)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&server, "server", "ngrok.j3r0lin.com:4443", "Ngrok server address")
	cmd.Flags().IntVar(&timeoutSec, "timeout", 60, "Task timeout in seconds")
	cmd.Flags().IntVar(&port, "port", 22, "Device-side SSH port")

	return cmd
}

// splitFirstLabel cuts host into its leading DNS label and the rest.
// "a-22.ngrok.10.5.17.73.nip.io" -> ("a-22", "ngrok.10.5.17.73.nip.io").
// Returns (host, "") when there is no dot.
func splitFirstLabel(host string) (string, string) {
	i := strings.IndexByte(host, '.')
	if i < 0 {
		return host, ""
	}
	return host[:i], host[i+1:]
}

// isHexUUID reports whether s is exactly 32 lowercase hex characters, the
// format of the ngrok tunnel uuid (128-bit routing token).
func isHexUUID(s string) bool {
	if len(s) != 32 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}
