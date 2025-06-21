package ffcmd

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Cmd is the command interface.
type Cmd interface {
	String() (string, error)
	Command() (*exec.Cmd, error)
}

// GenCommand returns an exec.Cmd with given command string.
func GenCommand(cmdStr string) (*exec.Cmd, error) {
	sh := ""
	arg1 := ""

	switch runtime.GOOS {
	case "darwin":
		sh = "bash"
		arg1 = "-c"
	case "linux":
		sh = "bash"
		arg1 = "-c"
	default:
		return nil, fmt.Errorf("not supported OS")
	}

	cmd := exec.Command(sh, arg1, cmdStr)
	return cmd, nil
}
