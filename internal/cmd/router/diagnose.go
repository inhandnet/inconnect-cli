package router

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

type diagnoseOptions struct {
	After  string
	Before string
}

func newCmdDiagnose(f *factory.Factory) *cobra.Command {
	opts := &diagnoseOptions{}

	cmd := &cobra.Command{
		Use:   "diagnose <router-id>",
		Short: "Aggregate multi-source diagnostics for a router into one report",
		Long: `Pull a router's diagnostics from several sources in one shot:

  - connectionLogs   VPN session logs (vpn-controller)
  - vpnEvents        VPN auth/connection events (vpn-controller)
  - connectionEvents device MQTT online/offline events (site)
  - registerEvents   device registration events (site)

Output is a single JSON object keyed by source. Individual sources that fail
are reported to stderr and left as empty arrays (best-effort).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			rid := args[0]
			oid, _ := cmd.Flags().GetString("oid")

			routerQ := url.Values{}
			if oid != "" {
				routerQ.Set("oid", oid)
			}
			routerBody, err := client.Get("/api/invpn/router/"+rid, routerQ)
			if err != nil {
				return err
			}
			routerRaw := resultRaw(routerBody, "{}")
			sn := gjson.GetBytes(routerBody, "result.serialNumber").String()

			after := cmdutil.ParseTimeFlag(opts.After)
			before := cmdutil.ParseTimeFlag(opts.Before)
			timeQ := func(extra url.Values) url.Values {
				q := url.Values{}
				for k, v := range extra {
					q[k] = v
				}
				if oid != "" {
					q.Set("oid", oid)
				}
				if after != "" {
					q.Set("after", after)
				}
				if before != "" {
					q.Set("before", before)
				}
				return q
			}

			fetch := func(label, path string, q url.Values) json.RawMessage {
				body, err := client.Get(path, q)
				if err != nil {
					fmt.Fprintf(f.IO.ErrOut, "warning: %s failed: %v\n", label, err)
					return json.RawMessage("[]")
				}
				return json.RawMessage(resultRaw(body, "[]"))
			}

			connQ := timeQ(url.Values{"rid": {rid}})
			eventQ := timeQ(url.Values{"rid": {rid}})

			report := map[string]json.RawMessage{
				"router":           json.RawMessage(routerRaw),
				"connectionLogs":   fetch("connection-logs", "/api/invpn/connection-logs", connQ),
				"vpnEvents":        fetch("vpn-events", "/api/invpn/vpn-events", eventQ),
				"connectionEvents": fetch("connection-events", "/api/devices/"+rid+"/connection-events", timeQ(url.Values{})),
			}
			if sn != "" {
				report["registerEvents"] = fetch("register-events", "/api/devices/"+sn+"/register-events", timeQ(url.Values{}))
			} else {
				report["registerEvents"] = json.RawMessage("[]")
			}

			out, err := json.Marshal(report)
			if err != nil {
				return err
			}
			return iostreams.FormatOutput(out, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.After, "after", "", "Only records at/after this time (e.g. 2024-01-01T00:00:00Z)")
	cmd.Flags().StringVar(&opts.Before, "before", "", "Only records before this time")

	return cmd
}

// resultRaw returns the raw JSON of the response's `result` field, or fallback
// when absent. It tolerates both list-wrapped ({"result":[...]}) and
// object-wrapped ({"result":{...}}) responses.
func resultRaw(body []byte, fallback string) string {
	if r := gjson.GetBytes(body, "result"); r.Exists() {
		return r.Raw
	}
	return fallback
}
