// Command docgen generates Markdown reference docs for every inconnect subcommand.
//
// It mirrors the command tree assembled in cmd/inconnect/main.go and writes one
// Markdown file per command into the target directory (default: docs/commands,
// overridable via the first CLI argument). Filenames have spaces replaced with
// underscores so the inconnect-skills plugin can map a subcommand path to its
// reference file (e.g. "inconnect router list" -> inconnect_router_list.md).
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"github.com/inhandnet/inconnect-cli/internal/cmd"
	"github.com/inhandnet/inconnect-cli/internal/cmd/alert"
	cmdapi "github.com/inhandnet/inconnect-cli/internal/cmd/api"
	"github.com/inhandnet/inconnect-cli/internal/cmd/auditlog"
	"github.com/inhandnet/inconnect-cli/internal/cmd/auth"
	"github.com/inhandnet/inconnect-cli/internal/cmd/banner"
	"github.com/inhandnet/inconnect-cli/internal/cmd/billing"
	cmdconfig "github.com/inhandnet/inconnect-cli/internal/cmd/config"
	"github.com/inhandnet/inconnect-cli/internal/cmd/connectionlog"
	"github.com/inhandnet/inconnect-cli/internal/cmd/datausage"
	"github.com/inhandnet/inconnect-cli/internal/cmd/drc"
	"github.com/inhandnet/inconnect-cli/internal/cmd/endpoint"
	"github.com/inhandnet/inconnect-cli/internal/cmd/firmware"
	"github.com/inhandnet/inconnect-cli/internal/cmd/mail"
	"github.com/inhandnet/inconnect-cli/internal/cmd/network"
	"github.com/inhandnet/inconnect-cli/internal/cmd/org"
	"github.com/inhandnet/inconnect-cli/internal/cmd/registerlog"
	"github.com/inhandnet/inconnect-cli/internal/cmd/role"
	"github.com/inhandnet/inconnect-cli/internal/cmd/router"
	"github.com/inhandnet/inconnect-cli/internal/cmd/server"
	"github.com/inhandnet/inconnect-cli/internal/cmd/system"
	"github.com/inhandnet/inconnect-cli/internal/cmd/task"
	"github.com/inhandnet/inconnect-cli/internal/cmd/user"
	"github.com/inhandnet/inconnect-cli/internal/cmd/vpnevent"
	"github.com/inhandnet/inconnect-cli/internal/factory"
)

func main() {
	f := factory.New()
	root := cmd.NewCmdRoot(f)
	root.AddCommand(
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
		connectionlog.NewCmdConnectionLog(f),
		vpnevent.NewCmdVpnEvent(f),
		system.NewCmdSystem(f),
		cmdapi.NewCmdAPI(f),
	)

	dir := "docs/commands"
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	cleanDir := filepath.Clean(dir)

	if err := os.MkdirAll(cleanDir, 0o750); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	filePrepender := func(string) string { return "" }
	linkHandler := func(name string) string { return strings.TrimSuffix(name, ".md") }

	disableAutoGenTag(root)

	if err := doc.GenMarkdownTreeCustom(root, cleanDir, filePrepender, linkHandler); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating docs: %v\n", err)
		os.Exit(1)
	}

	renameFiles(cleanDir)

	count := 0
	_ = filepath.Walk(cleanDir, func(_ string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			count++
		}
		return nil
	})

	cleanGeneratedFiles(cleanDir)

	fmt.Printf("Generated %d command docs in %s/\n", count, cleanDir)
}

func disableAutoGenTag(c *cobra.Command) {
	c.DisableAutoGenTag = true
	for _, sub := range c.Commands() {
		disableAutoGenTag(sub)
	}
}

// renameFiles replaces spaces with underscores in generated filenames so the
// skills plugin can derive a reference path from a subcommand path.
func renameFiles(dir string) {
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		newName := strings.ReplaceAll(info.Name(), " ", "_")
		if newName != info.Name() {
			return os.Rename(path, filepath.Join(filepath.Dir(path), newName))
		}
		return nil
	})
}

// cleanGeneratedFiles strips cobra's auto-generated footer and the SEE ALSO
// cross-link sections, which point at filenames the skills repo doesn't serve.
func cleanGeneratedFiles(dir string) {
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		content := string(data)

		lines := strings.Split(content, "\n")
		var cleaned []string
		for _, line := range lines {
			if strings.HasPrefix(line, "###### Auto generated") {
				continue
			}
			cleaned = append(cleaned, line)
		}
		content = strings.TrimRight(strings.Join(cleaned, "\n"), "\n") + "\n"

		if idx := strings.Index(content, "### SEE ALSO"); idx >= 0 {
			content = strings.TrimRight(content[:idx], "\n") + "\n"
		}

		return os.WriteFile(path, []byte(content), 0o600)
	})
}
