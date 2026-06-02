package auth

import (
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewCmdAuth(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with InConnect",
	}

	cmd.AddCommand(
		newCmdLogin(f),
		newCmdLogout(f),
		newCmdStatus(f),
		newCmdRegister(f),
		newCmdImpersonate(f),
		newCmdSwitchOrg(f),
		newCmdOrgs(f),
	)

	return cmd
}
