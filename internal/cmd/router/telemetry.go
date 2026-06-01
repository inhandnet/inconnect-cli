package router

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdTrafficDay(f *factory.Factory) *cobra.Command {
	var month string

	cmd := &cobra.Command{
		Use:   "traffic-day <id>",
		Short: "Get a device's per-day data traffic for a month",
		Long: `Get this device's daily data traffic (one entry per day) for a month.

This is the device-level traffic tracked by site. For organization/router/account
VPN data usage (a different service), use the top-level "data-usage" commands.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			m := strings.ReplaceAll(month, "-", "")
			if m == "" {
				m = time.Now().Format("200601")
			}

			q := url.Values{}
			q.Set("month", m)
			q.Set("device_id", args[0])
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			body, err := client.Get("/api/traffic_day", q)
			if err != nil {
				if body != nil {
					_ = iostreams.FormatOutput(body, f.IO, f.IO.Output)
				}
				return err
			}
			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&month, "month", "", "Month, YYYY-MM or YYYYMM (default current month)")
	return cmd
}

func newCmdOnlineTrend(f *factory.Factory) *cobra.Command {
	var after, before string

	cmd := &cobra.Command{
		Use:   "online-trend <id>",
		Short: "Get a device's online/offline trend over a time range",
		Long: `Get this device's online/offline state over time (a time series of
[timestamp, 0|1] where 1 = online).

For an instantaneous online/offline count across the organization, use
"router stats" instead.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			// online_tendency wants epoch seconds; default to the last 24h.
			start := cmdutil.ParseTimeFlagUnix(after)
			if start == "" {
				start = strconv.FormatInt(time.Now().Add(-24*time.Hour).Unix(), 10)
			}
			end := cmdutil.ParseTimeFlagUnix(before)
			if end == "" {
				end = strconv.FormatInt(time.Now().Unix(), 10)
			}

			q := url.Values{}
			q.Set("start_time", start)
			q.Set("end_time", end)
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			reqBody := map[string]any{"resourceIds": []string{args[0]}}

			respBody, err := client.Do("POST", "/api/online_tendency", q, reqBody)
			if err != nil {
				if respBody != nil {
					_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				}
				return err
			}
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&after, "after", "", "Start time, e.g. 2026-06-01 (default 24h ago)")
	cmd.Flags().StringVar(&before, "before", "", "End time, e.g. 2026-06-02 (default now)")
	return cmd
}

func newCmdSignal(f *factory.Factory) *cobra.Command {
	var after, before, fields string

	cmd := &cobra.Command{
		Use:   "signal <id>",
		Short: "Get a device's cellular signal time series (strength + quality)",
		Long: `Get this device's cellular (modem) signal time series from site, mirroring
the web UI which queries two endpoints together:

  strength  GET /api/devices/{id}/signal          -> [time, rssi] / asu
  quality   GET /api/devices/{id}/signal-quality   -> asu, rssi, rscp, rsrp,
            rsrq, sinr, ssRsrp, ssRsrq, ssSinr, ecio, pci, cid

Output is a JSON object {"strength":..., "quality":...}.
The backend defaults to the last 7 days when --after/--before are omitted.
Use --fields to restrict the quality metrics (comma-separated).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			a := cmdutil.ParseTimeFlag(after)
			b := cmdutil.ParseTimeFlag(before)
			oid, _ := cmd.Flags().GetString("oid")

			// Signal strength: /signal uses begin/end.
			sq := url.Values{}
			if a != "" {
				sq.Set("begin", a)
			}
			if b != "" {
				sq.Set("end", b)
			}
			if oid != "" {
				sq.Set("oid", oid)
			}
			strengthBody, err := client.Get("/api/devices/"+args[0]+"/signal", sq)
			if err != nil {
				if strengthBody != nil {
					_ = iostreams.FormatOutput(strengthBody, f.IO, f.IO.Output)
				}
				return err
			}

			// Signal quality: /signal-quality uses after/before (+fields).
			qq := url.Values{}
			if a != "" {
				qq.Set("after", a)
			}
			if b != "" {
				qq.Set("before", b)
			}
			if fields != "" {
				qq.Set("fields", fields)
			}
			if oid != "" {
				qq.Set("oid", oid)
			}
			qualityBody, err := client.Get("/api/devices/"+args[0]+"/signal-quality", qq)
			if err != nil {
				if qualityBody != nil {
					_ = iostreams.FormatOutput(qualityBody, f.IO, f.IO.Output)
				}
				return err
			}

			combined, err := json.Marshal(map[string]json.RawMessage{
				"strength": json.RawMessage(strengthBody),
				"quality":  json.RawMessage(qualityBody),
			})
			if err != nil {
				return err
			}
			return iostreams.FormatOutput(combined, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&after, "after", "", "Start time, e.g. 2026-06-01 (default 7 days ago)")
	cmd.Flags().StringVar(&before, "before", "", "End time, e.g. 2026-06-02 (default now)")
	cmd.Flags().StringVar(&fields, "fields", "", "Comma-separated quality metrics to return (default all)")
	return cmd
}
