package ffcmd

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
)

// Cmd is the interface of shell command for FFmpeg binaries(e.g. ffmpeg, ffprobe).
// String returns the command string of ffmpeg binaries(e.g. "ffmpeg -i input.MOV output.mp4").
type Cmd interface {
	String() (string, error)
}

// CommandContext converts a [Cmd] to an [os/exec.Cmd].
func CommandContext(ctx context.Context, cmd Cmd) (*exec.Cmd, error) {
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

// Command converts a [Cmd] to an [os/exec.Cmd] like [CommandContext] but without a [context.Context].
func Command(cmd Cmd) (*exec.Cmd, error) {
	return CommandContext(context.Background(), cmd)
}
