package router

import (
	"net/url"

	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type createOptions struct {
	SerialNumber string
	Name         string
	Model        string
	LanInterface string
	ModelID      string
	NetworkID    string
	Subnet       string
	CustomFields []string
}

func newCmdCreate(f *factory.Factory) *cobra.Command {
	opts := &createOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a VPN router",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			body := map[string]any{
				"serialNumber": opts.SerialNumber,
			}
			if opts.Name != "" {
				body["name"] = opts.Name
			}
			if opts.Model != "" {
				body["model"] = opts.Model
			}
			if opts.LanInterface != "" {
				body["lanInterface"] = opts.LanInterface
			}
			if opts.ModelID != "" {
				body["modelId"] = opts.ModelID
			}
			if opts.NetworkID != "" {
				body["networkId"] = opts.NetworkID
			}
			if opts.Subnet != "" {
				body["subnet"] = opts.Subnet
			}
			if len(opts.CustomFields) > 0 {
				cf := parseKeyValues(opts.CustomFields)
				if len(cf) > 0 {
					body["customFields"] = cf
				}
			}

			respBody, err := client.Do("POST", "/api/invpn/router", q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteCreated(f, "Router", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.SerialNumber, "serial", "", "Device serial number (required, 15 chars)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Router name")
	cmd.Flags().StringVar(&opts.Model, "model", "", "Router model (e.g. IR900)")
	cmd.Flags().StringVar(&opts.LanInterface, "lan-interface", "", "LAN interface")
	cmd.Flags().StringVar(&opts.ModelID, "model-id", "", "Model ID")
	cmd.Flags().StringVar(&opts.NetworkID, "network-id", "", "Network ID to assign")
	cmd.Flags().StringVar(&opts.Subnet, "subnet", "", "Router subnet")
	cmd.Flags().StringSliceVar(&opts.CustomFields, "custom-field", nil, "Custom field (key=value, can specify multiple)")
	_ = cmd.MarkFlagRequired("serial")

	return cmd
}

func parseKeyValues(pairs []string) map[string]string {
	m := make(map[string]string)
	for _, p := range pairs {
		for i := 0; i < len(p); i++ {
			if p[i] == '=' {
				m[p[:i]] = p[i+1:]
				break
			}
		}
	}
	return m
}
