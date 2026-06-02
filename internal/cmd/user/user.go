package user

import (
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewCmdUser(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage VPN user accounts",
	}

	cmd.AddCommand(
		newCmdList(f),
		newCmdCreate(f),
		newCmdUpdate(f),
		newCmdDelete(f),
		newCmdResetPassword(f),
		newCmdLock(f),
		newCmdUnlock(f),
		newCmdSetFloatAddress(f),
		newCmdBindMac(f),
		newCmdIssueKeypair(f),
		newCmdBatchIssueKeypair(f),
	)

	return cmd
}
