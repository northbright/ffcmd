package ffcmd

// Cmd is the command interface.
type Cmd interface {
	String() (string, error)
}
