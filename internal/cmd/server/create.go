package server

import (
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type createOptions struct {
	OID         string
	Group       string
	Version     string
	Subnet      string
	Address     string
	Netmask     string
	IfconfigPool string
	Host        string
	Port        int
	NodePort    int
	ServiceType string
	Proto       string
	Deploy      bool
}

func newCmdCreate(f *factory.Factory) *cobra.Command {
	opts := &createOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a VPN server",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			oid := opts.OID
			if oid == "" {
				oid, _ = cmd.Flags().GetString("oid")
			}

			body := map[string]any{
				"oid": oid,
			}
			if opts.Group != "" {
				body["group"] = opts.Group
			}
			if opts.Version != "" {
				body["version"] = opts.Version
			}
			if opts.Subnet != "" {
				body["subnet"] = opts.Subnet
			}
			if opts.Address != "" {
				body["address"] = opts.Address
			}
			if opts.Netmask != "" {
				body["netmask"] = opts.Netmask
			}
			if opts.IfconfigPool != "" {
				body["ifconfigPool"] = opts.IfconfigPool
			}
			if opts.Host != "" {
				body["host"] = opts.Host
			}
			if opts.Port > 0 {
				body["port"] = opts.Port
			}
			if opts.NodePort > 0 {
				body["nodePort"] = opts.NodePort
			}
			if opts.ServiceType != "" {
				body["serviceType"] = opts.ServiceType
			}
			if opts.Proto != "" {
				body["proto"] = opts.Proto
			}

			q := url.Values{}
			if !opts.Deploy {
				q.Set("deploy", "false")
			}

			respBody, err := client.Do("POST", "/api/invpn/server", q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteCreated(f, "Server", respBody)
			return iostreams.FormatOutput(redactBody(cmd, respBody), f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.OID, "org-id", "", "Organization ID (required)")
	cmd.Flags().StringVar(&opts.Group, "group", "", "Server group name (default: \"default\")")
	cmd.Flags().StringVar(&opts.Version, "version", "", "Server version")
	cmd.Flags().StringVar(&opts.Subnet, "subnet", "", "Server subnet (e.g. 10.8.0.0/16)")
	cmd.Flags().StringVar(&opts.Address, "address", "", "Server address")
	cmd.Flags().StringVar(&opts.Netmask, "netmask", "", "Server netmask")
	cmd.Flags().StringVar(&opts.IfconfigPool, "ifconfig-pool", "", "Ifconfig pool")
	cmd.Flags().StringVar(&opts.Host, "host", "", "Server host")
	cmd.Flags().IntVar(&opts.Port, "port", 0, "Server port")
	cmd.Flags().IntVar(&opts.NodePort, "node-port", 0, "Kubernetes node port")
	cmd.Flags().StringVar(&opts.ServiceType, "service-type", "", "Kubernetes service type")
	cmd.Flags().StringVar(&opts.Proto, "proto", "", "Protocol (tcp or udp)")
	cmd.Flags().BoolVar(&opts.Deploy, "deploy", false, "Deploy server (provision K8s pod) after creation; default creates the DB record only")

	return cmd
}
