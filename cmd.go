package ffcmd

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
