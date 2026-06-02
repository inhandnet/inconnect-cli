package user

import (
	"fmt"
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdBindMac(f *factory.Factory) *cobra.Command {
	var macs []string

	cmd := &cobra.Command{
		Use:   "bind-mac <id>",
		Short: "Bind MAC addresses to a user",
		Long:  "Bind MAC addresses to a user.\n\n<id> is the VPN user ID (the 'id' field from 'user list', not 'uid').",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"macs": macs,
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			respBody, err := client.Do("PUT", "/api/invpn/users/"+args[0]+"/bindMacAddress", q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "MAC addresses bound for user %s\n", args[0])
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringSliceVar(&macs, "mac", nil, "MAC address (can specify multiple)")

	return cmd
}

func newCmdIssueKeypair(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "issue-keypair <id>",
		Short: "Issue a new key pair for a user",
		Long:  "Issue a new key pair for a user.\n\n<id> is the VPN user ID (the 'id' field from 'user list', not 'uid').",
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

			respBody, err := client.Do("PUT", "/api/invpn/users/"+args[0]+"/issueNewKeyPair", q, nil)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "New key pair issued for user %s\n", args[0])
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}
}

func newCmdBatchIssueKeypair(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-issue-keypair <id> [id...]",
		Short: "Issue new key pairs for multiple users",
		Long:  "Issue new key pairs for multiple users.\n\nArguments are VPN user IDs (the 'id' field from 'user list', not 'uid').",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"ids": args,
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			respBody, err := client.Do("POST", "/api/invpn/users/issueNewKeyPair", q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "New key pairs issued for %d users\n", len(args))
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	return cmd
}
