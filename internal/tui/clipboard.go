package tui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// CopyToClipboard 將指定文字複製到作業系統剪貼簿 (macOS: pbcopy, Windows: powershell/clip, Linux: xclip/wl-copy)
func CopyToClipboard(text string) error {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("pbcopy")
		cmd.Stdin = strings.NewReader(text)
		return cmd.Run()
	case "windows":
		cmd := exec.Command("powershell", "-command", "Set-Clipboard", "-Value", fmt.Sprintf("%q", text))
		return cmd.Run()
	case "linux":
		cmd := exec.Command("xclip", "-selection", "clipboard")
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			return nil
		}
		cmdWl := exec.Command("wl-copy")
		cmdWl.Stdin = strings.NewReader(text)
		return cmdWl.Run()
	default:
		return fmt.Errorf("unsupported OS for clipboard copy: %s", runtime.GOOS)
	}
}
