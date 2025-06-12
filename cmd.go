package ffcmd

type Cmd interface {
	String() (string, error)
}
