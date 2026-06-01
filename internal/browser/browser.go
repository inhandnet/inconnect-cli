//go:build windows

package browser

import "os/exec"

func Cmd(targetURL string) *exec.Cmd {
	return exec.Command("rundll32", "url.dll,FileProtocolHandler", targetURL)
}

func Open(targetURL string) {
	cmd := Cmd(targetURL)
	if cmd != nil {
		_ = cmd.Start()
	}
}
