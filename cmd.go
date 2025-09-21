package ffcmd

// ConcatCmds concatenates commands to a new command.
// cmd: current command. Leave it empty to return a command which contains cmds only.
// cmds: commands to concatenate.
func ConcatCmds(cmd string, cmds ...string) string {
	str := ""
	l := len(cmds)

	if cmd != "" {
		str += cmd
		if l > 0 {
			str += " && "
		}
	}

	for i, cmd := range cmds {
		if cmd != "" {
			str += cmd
			if i != l-1 {
				str += " && "
			}
		}
	}

	return str
}
