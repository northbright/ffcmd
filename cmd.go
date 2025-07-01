package ffcmd

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
)

// Cmd is the command interface.
type Cmd interface {
	String() (string, error)
}

// commandContext returns an *exec.Cmd by given shell command.
func commandContext(ctx context.Context, cmd Cmd) (*exec.Cmd, error) {
	str, err := cmd.String()
	if err != nil {
		return nil, fmt.Errorf("get command string error: %v", err)
	}

	switch runtime.GOOS {
	case "darwin":
		return exec.CommandContext(ctx, "zsh", "-c", str), nil
	case "linux":
		return exec.CommandContext(ctx, "bash", "-c", str), nil
	default:
		return nil, fmt.Errorf("not supported os")
	}
}
