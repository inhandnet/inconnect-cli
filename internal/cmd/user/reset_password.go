package user

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type resetPasswordOptions struct {
	Language int
}

func newCmdResetPassword(f *factory.Factory) *cobra.Command {
	opts := &resetPasswordOptions{}

	cmd := &cobra.Command{
		Use:   "reset-password <uid>",
		Short: "Reset a user's password (NOT usable by admins — see below)",
		Long: `Reset a user's password.

<uid> is the account ID (the 'uid' field from 'user list', not 'id').

Note: the backend endpoint expects a verification 'code' from the self-service
"forgot password" email flow, which this command does not supply. As a result it
currently cannot perform a standalone admin-initiated reset.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			q.Set("language", strconv.Itoa(opts.Language))

			body := map[string]any{}

			respBody, err := client.Do("PUT", "/api2/users/"+args[0]+"/reset_password", q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Password reset for user %s\n", args[0])
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().IntVar(&opts.Language, "language", 1, "Notification language (1=English, 2=Chinese)")

	return cmd
}
