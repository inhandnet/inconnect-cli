package main

import (
	"fmt"
	"os"

	"github.com/inhandnet/ics-cli/internal/cmd"
	"github.com/inhandnet/ics-cli/internal/cmd/alert"
	cmdapi "github.com/inhandnet/ics-cli/internal/cmd/api"
	"github.com/inhandnet/ics-cli/internal/cmd/auditlog"
	"github.com/inhandnet/ics-cli/internal/cmd/auth"
	"github.com/inhandnet/ics-cli/internal/cmd/banner"
	"github.com/inhandnet/ics-cli/internal/cmd/billing"
	cmdconfig "github.com/inhandnet/ics-cli/internal/cmd/config"
	"github.com/inhandnet/ics-cli/internal/cmd/datausage"
	"github.com/inhandnet/ics-cli/internal/cmd/drc"
	"github.com/inhandnet/ics-cli/internal/cmd/endpoint"
	"github.com/inhandnet/ics-cli/internal/cmd/firmware"
"github.com/inhandnet/ics-cli/internal/cmd/mail"
	"github.com/inhandnet/ics-cli/internal/cmd/network"
	"github.com/inhandnet/ics-cli/internal/cmd/org"
	"github.com/inhandnet/ics-cli/internal/cmd/registerlog"
	"github.com/inhandnet/ics-cli/internal/cmd/role"
	"github.com/inhandnet/ics-cli/internal/cmd/router"
	"github.com/inhandnet/ics-cli/internal/cmd/server"
	"github.com/inhandnet/ics-cli/internal/cmd/system"
	"github.com/inhandnet/ics-cli/internal/cmd/task"
	"github.com/inhandnet/ics-cli/internal/cmd/user"
	"github.com/inhandnet/ics-cli/internal/factory"
)

func main() {
	f := factory.New()
	rootCmd := cmd.NewCmdRoot(f)

	rootCmd.AddCommand(
		auth.NewCmdAuth(f),
		cmdconfig.NewCmdConfig(f),
		network.NewCmdNetwork(f),
		server.NewCmdServer(f),
		router.NewCmdRouter(f),
		endpoint.NewCmdEndpoint(f),
		alert.NewCmdAlert(f),
		task.NewCmdTask(f),
		firmware.NewCmdFirmware(f),
		drc.NewCmdDRC(f),
		user.NewCmdUser(f),
		role.NewCmdRole(f),
datausage.NewCmdDataUsage(f),
		org.NewCmdOrg(f),
		billing.NewCmdBilling(f),
		banner.NewCmdBanner(f),
		mail.NewCmdMail(f),
		registerlog.NewCmdRegisterLog(f),
		auditlog.NewCmdAuditLog(f),
		system.NewCmdSystem(f),
		cmdapi.NewCmdAPI(f),
	)

	if _, err := rootCmd.ExecuteC(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
