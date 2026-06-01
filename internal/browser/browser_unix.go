//go:build !windows

package browser

import (
	"os/exec"
	"runtime"
)

func Cmd(targetURL string) *exec.Cmd {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", targetURL)
	default:
		return exec.Command("xdg-open", targetURL)
	}
}

func Open(targetURL string) {
	cmd := Cmd(targetURL)
	if cmd != nil {
		_ = cmd.Start()
	}
}
