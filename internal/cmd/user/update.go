package user

import (
	"net/url"

	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type updateOptions struct {
	Name     string
	RoleID   string
	ExpireAt int64
	Lock     bool
}

func newCmdUpdate(f *factory.Factory) *cobra.Command {
	opts := &updateOptions{}

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a VPN user account",
		Long:  "Update a VPN user account.\n\n<id> is the VPN user ID (the 'id' field from 'user list', not 'uid').",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if opts.Name != "" {
				body["name"] = opts.Name
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

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			respBody, err := client.Do("PUT", "/api/invpn/user/"+args[0], q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteUpdated(f, "User", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "User name")
	cmd.Flags().StringVar(&opts.RoleID, "role-id", "", "Role ID")
	cmd.Flags().Int64Var(&opts.ExpireAt, "expire-at", 0, "Expiration timestamp (epoch millis)")
	cmd.Flags().BoolVar(&opts.Lock, "lock", false, "Lock or unlock the user")

	return cmd
}
