package ffcmd

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

type ReadOutputFunc func(stdout, stderr io.ReadCloser) error

// Cmd is the command interface.
type Cmd interface {
	String() (string, error)
	Run(dir string, fn ReadOutputFunc) error
}

func RunCmd(dir, cmdStr string, fn ReadOutputFunc) error {
	cmd := exec.Command("bash", "-c", cmdStr)

	// Set working dir.
	cmd.Dir = dir
	log.Printf("----------- cmd.Dir: %s\n", cmd.Dir)

	// Create stdout, stderr pipes.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd.Start() error: %v", err)
	}

	if fn != nil {
		if err := fn(stdout, stderr); err != nil {
			return fmt.Errorf("read output function error: %v", err)
		}
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("cmd.Wait() error: %v", err)
	}

	return nil
}
