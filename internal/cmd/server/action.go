package server

import (
	"fmt"
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdStop(f *factory.Factory) *cobra.Command {
	var group string

	cmd := &cobra.Command{
		Use:   "stop <org-id>",
		Short: "Stop VPN server(s) for an organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"action": "stop",
			}

			path := "/api/invpn/org/" + args[0] + "/servers"
			if group != "" {
				q := url.Values{}
				q.Set("group", group)
				path = "/api/invpn/org/" + args[0] + "/server"
				respBody, err := client.Do("POST", path, q, body)
				if err != nil {
					_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
					return err
				}
				fmt.Fprintf(f.IO.ErrOut, "Server (group=%s) stopped for org %s\n", group, args[0])
				return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
			}

			respBody, err := client.Post(path, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "All servers stopped for org %s\n", args[0])
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&group, "group", "", "Server group (if omitted, stops all servers)")

	return cmd
}

func newCmdRecover(f *factory.Factory) *cobra.Command {
	var group string

	cmd := &cobra.Command{
		Use:   "recover <org-id>",
		Short: "Recover VPN server(s) for an organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"action": "recover",
			}

			path := "/api/invpn/org/" + args[0] + "/servers"
			if group != "" {
				q := url.Values{}
				q.Set("group", group)
				path = "/api/invpn/org/" + args[0] + "/server"
				respBody, err := client.Do("POST", path, q, body)
				if err != nil {
					_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
					return err
				}
				fmt.Fprintf(f.IO.ErrOut, "Server (group=%s) recovered for org %s\n", group, args[0])
				return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
			}

			respBody, err := client.Post(path, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "All servers recovered for org %s\n", args[0])
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&group, "group", "", "Server group (if omitted, recovers all servers)")

	return cmd
}

func newCmdIssueKeypair(f *factory.Factory) *cobra.Command {
	var serverID, deployAt string

	cmd := &cobra.Command{
		Use:   "issue-keypair <org-id>",
		Short: "Issue new key pair for server(s)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if deployAt != "" {
				body["deployAt"] = deployAt
			}

			path := "/api/invpn/orgs/" + args[0] + "/servers/issueNewKeyPair"
			if serverID != "" {
				path = "/api/invpn/orgs/" + args[0] + "/servers/" + serverID + "/issueNewKeyPair"
			}

			respBody, err := client.Put(path, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "New key pair issued for org %s\n", args[0])
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&serverID, "server-id", "", "Specific server ID (if omitted, issues for all servers)")
	cmd.Flags().StringVar(&deployAt, "deploy-at", "", "Scheduled deployment date (ISO 8601)")

	return cmd
}
