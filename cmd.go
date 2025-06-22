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
	Command() (*exec.Cmd, error)
	CommandContext(ctx context.Context) (*exec.Cmd, error)
}

// GenCommand returns an exec.Cmd with given command string.
func GenCommand(cmdStr string) (*exec.Cmd, error) {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("bash", "-c", cmdStr), nil
	case "linux":
		return exec.Command("bash", "-c", cmdStr), nil
	default:
		return nil, fmt.Errorf("not supported OS")
	}
}

// GenCommandContext returns an exec.Cmd like [GenCommand] but includes a context.
func GenCommandContext(ctx context.Context, cmdStr string) (*exec.Cmd, error) {
	switch runtime.GOOS {
	case "darwin":
		return exec.CommandContext(ctx, "bash", "-c", cmdStr), nil
	case "linux":
		return exec.CommandContext(ctx, "bash", "-c", cmdStr), nil
	default:
		return nil, fmt.Errorf("not supported OS")
	}
}
