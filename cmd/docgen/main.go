// Command docgen generates Markdown reference docs for every ics subcommand.
//
// It mirrors the command tree assembled in cmd/ics/main.go and writes one
// Markdown file per command into the target directory (default: docs/commands,
// overridable via the first CLI argument). Filenames have spaces replaced with
// underscores so the ics-skills plugin can map a subcommand path to its
// reference file (e.g. "ics router list" -> ics_router_list.md).
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

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
