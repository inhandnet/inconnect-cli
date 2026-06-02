package user

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

type lockUnlockOptions struct {
	Notice   bool
	Language int
}

func newCmdLock(f *factory.Factory) *cobra.Command {
	opts := &lockUnlockOptions{Notice: true}

	cmd := &cobra.Command{
		Use:   "lock <uid>",
		Short: "Lock a user account",
		Long:  "Lock a user account.\n\n<uid> is the account ID (the 'uid' field from 'user list', not 'id').\nBy default an email notification is sent; pass --notice=false to suppress it.",
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
			q.Set("notice", strconv.FormatBool(opts.Notice))
			if opts.Language > 0 {
				q.Set("language", strconv.Itoa(opts.Language))
			}

			_, err = client.Do("PUT", "/api/users/"+args[0]+"/lock", q, nil)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "User %s locked\n", args[0])
			return nil
		},
	}

	cmd.Flags().BoolVar(&opts.Notice, "notice", true, "Send notification email")
	cmd.Flags().IntVar(&opts.Language, "language", 1, "Notification language (1=English, 2=Chinese)")

	return cmd
}

func newCmdUnlock(f *factory.Factory) *cobra.Command {
	opts := &lockUnlockOptions{Notice: true}

	cmd := &cobra.Command{
		Use:   "unlock <uid>",
		Short: "Unlock a user account",
		Long:  "Unlock a user account.\n\n<uid> is the account ID (the 'uid' field from 'user list', not 'id').\nBy default an email notification is sent; pass --notice=false to suppress it.",
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
			q.Set("notice", strconv.FormatBool(opts.Notice))
			if opts.Language > 0 {
				q.Set("language", strconv.Itoa(opts.Language))
			}

			_, err = client.Do("PUT", "/api/users/"+args[0]+"/unlock", q, nil)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "User %s unlocked\n", args[0])
			return nil
		},
	}

	cmd.Flags().BoolVar(&opts.Notice, "notice", true, "Send notification email")
	cmd.Flags().IntVar(&opts.Language, "language", 1, "Notification language (1=English, 2=Chinese)")

	return cmd
}
