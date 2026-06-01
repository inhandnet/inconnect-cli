package user

import (
	"net/url"
	"strconv"

	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type createOptions struct {
	Name      string
	Email     string
	OID       string
	RoleID    string
	ExpireAt  int64
	Lock      bool
	NetworkID string
	External  bool
	Language  int
}

func newCmdCreate(f *factory.Factory) *cobra.Command {
	opts := &createOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a VPN user account",
		Long: `Create a VPN user account.

Each VPN user has two IDs, both shown by 'user list':
  - id  : VPN user ID — used by update/set-float-address/bind-mac/issue-keypair
  - uid : account ID   — used by lock/unlock/delete/reset-password

Note: --role-id is required by the server. Look up role IDs via 'ics role list'.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"name":  opts.Name,
				"email": opts.Email,
			}
			if opts.OID != "" {
				body["oid"] = opts.OID
			}
			if opts.RoleID != "" {
				body["roleId"] = opts.RoleID
			}
			if opts.ExpireAt > 0 {
				body["expireAt"] = opts.ExpireAt
			}
			if cmd.Flags().Changed("lock") {
				body["lock"] = opts.Lock
			}
			if opts.NetworkID != "" {
				body["networkId"] = opts.NetworkID
			}
			if cmd.Flags().Changed("external") {
				body["external"] = opts.External
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}
			if opts.Language > 0 {
				q.Set("language", strconv.Itoa(opts.Language))
			}

			respBody, err := client.Do("POST", "/api/invpn/user", q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteCreated(f, "User", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "User name (required)")
	cmd.Flags().StringVar(&opts.Email, "email", "", "User email (required)")
	cmd.Flags().StringVar(&opts.OID, "org-id", "", "Organization ID")
	cmd.Flags().StringVar(&opts.RoleID, "role-id", "", "Role ID (required; see 'ics role list')")
	cmd.Flags().Int64Var(&opts.ExpireAt, "expire-at", 0, "Expiration timestamp (epoch millis)")
	cmd.Flags().BoolVar(&opts.Lock, "lock", false, "Create in locked state")
	cmd.Flags().StringVar(&opts.NetworkID, "network-id", "", "Network ID to assign")
	cmd.Flags().BoolVar(&opts.External, "external", false, "External user")
	cmd.Flags().IntVar(&opts.Language, "language", 1, "Notification language (1=English, 2=Chinese)")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("role-id")

	return cmd
}
