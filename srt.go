package ffcmd

import (
	"fmt"
	"strings"
)

// NewCreateOneSubSRTCmd returns a command which writes only one subtitle to a SRT file.
// srtFile: SRT file name.
// text: the only one subtitle text.
// durationVarName: name of the shell variable which contains the duration of the video.
// It's the output of [NewGetDurationCmd].
func NewCreateOneSubSRTCmd(srtFile, text, durationVarName string) (string, error) {
	if srtFile == "" {
		return "", fmt.Errorf("empty SRT file name")
	}

	if text == "" {
		return "", fmt.Errorf("empty subtitle text")
	}

	text = strings.ReplaceAll(text, `"`, `\"`)

	// Shell commands slice.
	var cmds []string

	// Create a command to convert duration to subtitle end timestamp.
	getEndTimestampCmd := NewDurationToTimestampCmd(durationVarName, "end")

	// Add get clip duration command.
	cmds = append(cmds, getEndTimestampCmd)

	// Add command to write subtitle text to a SRT file.
	writeSRTCmd := fmt.Sprintf(`printf "1\n00:00:00,000 --> %%s\n%s" $end > "%s"`, text, srtFile)
	cmds = append(cmds, writeSRTCmd)

	// Rturn the final command.
	return ConcatCmds("", cmds...), nil
}

// NewRemoveOneSubSRTCmd returns a command which remove the SRT file.
// srtFile: SRT file name.
func NewRemoveOneSubSRTCmd(srtFile string) (string, error) {
	if srtFile == "" {
		return "", fmt.Errorf("empty SRT file name")
	}

	return fmt.Sprintf(`rm "%s"`, srtFile), nil
}
