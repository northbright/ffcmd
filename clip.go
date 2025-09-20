package ffcmd

/*
type Clip struct {
	File  string
	Start string
	End   string
}

func NewClip(file, start, end string) *Clip {
	return &Clip{File: file, Start: start, End: end}
}

func NewClipByDuration(duration float32) *Clip {
	ts := timestamp.NewFromSecond(duration)
	end := ts.String()

	return NewClip("", "", end)
}

func (c *Clip) NewGetDurationCmd(varName, ffprobePath string) (string, error) {
	var (
		err            error
		str            string
		tsStart, tsEnd *timestamp.Timestamp
		startSec       float32
	)

	// Check if the name of variable in shell is empty.
	if varName == "" {
		return "", fmt.Errorf("empty variable name")
	}

	// Remove '$' prefix.
	varName, _ = strings.CutPrefix(varName, "$")

	if ffprobePath == "" {
		ffprobePath = "ffprobe"
	}

	if c.Start == "" {
		c.Start = "00:00:00.000"
	}

	if tsStart, err = timestamp.New(c.Start); err != nil {
		return "", fmt.Errorf("invalid start time format")
	}

	startSec = tsStart.Second()

	if c.End == "" {
		if c.File == "" {
			return "", fmt.Errorf("both end time and video filename are empty, can not get end timestamp")
		}

		str = fmt.Sprintf(`d=$(%s -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "%s"); %s=$(echo $d - %.03f | bc)`, ffprobePath, c.File, varName, startSec)
	} else {
		if tsEnd, err = timestamp.New(c.End); err != nil {
			return "", fmt.Errorf("invalid end time format")
		}

		if tsEnd, err = tsEnd.Sub(tsStart); err != nil {
			return "", fmt.Errorf("tsEnd.Sub() error: %v", err)
		}

		str = fmt.Sprintf(`%s=%.03f`, varName, tsEnd.Second())
	}

	return str, nil
}

func NewDurationToTimestampCmd(durationVarName, timestampVarName string) string {
	return fmt.Sprintf(`d=$(printf "%%.3f" $%s) && sec=$(echo $d | awk -F. '{ print $1 }') && frac=$(echo $d | awk -F. '{ print substr($2, 1, 3) }') && hh=$((sec / 3600)) && mm=$((sec %% 3600 / 60)) && ss=$((sec %% 3600 %% 60)) && printf -v %s "%%02d:%%02d:%%02d,%%03d" $hh $mm $ss $frac`, durationVarName, timestampVarName)
}

func (c *Clip) NewCreateOneSubSRTCmd(srtFile, text, ffprobePath string) (string, error) {
	if srtFile == "" {
		return "", fmt.Errorf("empty SRT file name")
	}

	if text == "" {
		return "", fmt.Errorf("empty subtitle text")
	}

	text = strings.ReplaceAll(text, `"`, `\"`)

	// Shell commands slice.
	var cmds []string

	// Create a command to get clip duration and store it in the var("d").
	getDurationCmd, err := c.NewGetDurationCmd("d", ffprobePath)
	if err != nil {
		return "", err
	}
	cmds = append(cmds, getDurationCmd)

	// Create a command to convert duration to subtitle end timestamp.
	getEndTimestampCmd := NewDurationToTimestampCmd("d", "end")

	// Add get clip duration command.
	cmds = append(cmds, getEndTimestampCmd)

	// Add command to write subtitle text to a SRT file.
	writeSRTCmd := fmt.Sprintf(`printf "1\n00:00:00,000 --> %%s\n%s" $end > "%s"`, text, srtFile)
	cmds = append(cmds, writeSRTCmd)

	// Rturn the final command.
	return ConcatCmds("", cmds...), nil
}

func (c *Clip) NewRemoveOneSubSRTCmd(srtFile string) (string, error) {
	if srtFile == "" {
		return "", fmt.Errorf("empty SRT file name")
	}

	return fmt.Sprintf(`rm "%s"`, srtFile), nil
}
*/
