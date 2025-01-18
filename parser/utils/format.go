package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// FormatGoFile formats a Go file using gofmt
func FormatGoFile(filePath string) error {
	if err := RunCmd("gofmt", "-w", filePath); err != nil {
		return fmt.Errorf("failed to format %s: %w", filePath, err)
	}
	return nil
}

// RunCmd runs a command
func RunCmd(parts ...string) error {
	if len(parts) == 0 {
		return fmt.Errorf("no command provided")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command: %w\nOutput: %s", err, out.String())
	}

	return nil
}
