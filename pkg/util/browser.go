package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		cmd := exec.Command("xdg-open", url)
		cmd.Env = append(os.Environ(), "DISPLAY=")
		return cmd.Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
