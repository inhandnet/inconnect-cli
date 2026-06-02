package user

import (
	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newCmdDelete(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <uid>",
		Short: "Delete a user account",
		Long:  "Delete a user account.\n\n<uid> is the account ID (the 'uid' field from 'user list', not 'id').",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			_, err = client.Delete("/api2/users/" + args[0])
			if err != nil {
				return err
			}

			cmdutil.WriteDeleted(f, "User", args[0])
			return nil
		},
	}
}
