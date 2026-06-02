package router

import (
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewCmdRouter(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "router",
		Short: "Manage VPN routers",
	}

	cmd.AddCommand(
		newCmdList(f),
		newCmdGet(f),
		newCmdCreate(f),
		newCmdUpdate(f),
		newCmdDelete(f),
		newCmdConfigSend(f),
		newCmdSubnet(f),
		newCmdOvpn(f),
		newCmdClientOvpn(f),
		newCmdKick(f),
		newCmdReboot(f),
		newCmdStats(f),
		newCmdModels(f),
		newCmdDeviceConfig(f),
		newCmdTransfer(f),
		newCmdNextVip(f),
		newCmdSetRip(f),
		newCmdNatConf(f),
		newCmdLocations(f),
		newCmdWeb(f),
		newCmdExec(f),
		newCmdRunningConfig(f),
		newCmdTrafficDay(f),
		newCmdOnlineTrend(f),
		newCmdSignal(f),
	)

	return cmd
}
