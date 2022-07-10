package util

import (
	"fmt"
	"os/exec"
	"runtime"
)

func OpenInBrowser(url string) error {
	return openInBrowser(runtime.GOOS, url)
}

func openInBrowser(os string, url string) error {
	cmd, err := getBrowserCmd(os, url)
	if err != nil {
		return fmt.Errorf("failed build open browser cmd: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	return nil
}

func getBrowserCmd(os string, url string) (*exec.Cmd, error) {
	switch os {
	case "linux":
		return exec.Command("xdg-open", url), nil
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url), nil
	case "darwin":
		return exec.Command("open", url), nil
	default:
		return nil, fmt.Errorf("unsupported platform %v", os)
	}

}
