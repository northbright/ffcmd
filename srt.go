package ffcmd

import (
	"fmt"
	"strings"
)

type CreateOneSubSRTCmd struct {
	srtFile   string
	videoFile string
	text      string
	start     string
	end       string
}

func NewCreateOneSubSRTCmd(srtFile, videoFile, text, start, end string) (*CreateOneSubSRTCmd, error) {
	if srtFile == "" {
		return nil, fmt.Errorf("empty SRT file name")
	}

	if text == "" {
		return nil, fmt.Errorf("empty subtitle text")
	}

	// Escape exclamation mark.
	text = strings.ReplaceAll(text, "!", `\!`)

	return &CreateOneSubSRTCmd{srtFile: srtFile, videoFile: videoFile, text: text, start: start, end: end}, nil
}

func (cmd *CreateOneSubSRTCmd) String() (string, error) {
	var start, end string
	str := ""

	if cmd.start == "" {
		start = "00:00:00,000"
	} else {
		ts, err := NewTimestamp(cmd.start)
		if err != nil {
			fmt.Printf("cmd.start: %s\n", cmd.start)
			return "", fmt.Errorf("invalid start time format")
		}
		start = ts.String()
	}

	if cmd.end == "" {
		if cmd.videoFile == "" {
			return "", fmt.Errorf("both end time and video filename are empty, can not get end timestamp")
		}

		str = fmt.Sprintf(`ffprobe -v error -select_streams v:0 -show_entries stream=duration -of csv=s=,:p=0 "%s" | awk -F. '{ print $1 }' | read sec; hh=$((sec / 3600)); mm=$((sec %% 3600 / 60)); ss=$((sec %% 3600 %% 60)); printf -v end "%%02d:%%02d:%%02d,000" hh mm ss; echo -ne "1\n%s --> $end\n%s" > "%s"`, cmd.videoFile, start, cmd.text, cmd.srtFile)
	} else {
		ts, err := NewTimestamp(cmd.end)
		if err != nil {
			return "", fmt.Errorf("invalid end time format")
		}
		end = ts.String()

		str = fmt.Sprintf(`echo -ne "1\n%s --> %s\n%s" > "%s"`, start, end, cmd.text, cmd.srtFile)
	}

	return str, nil
}

type RemoveOneSubSRTCmd struct {
	srtFile string
}

func NewRemoveOneSubSRTCmd(srtFile string) (*RemoveOneSubSRTCmd, error) {
	if srtFile == "" {
		return nil, fmt.Errorf("empty SRT filename")
	}

	return &RemoveOneSubSRTCmd{srtFile: srtFile}, nil
}

func (cmd *RemoveOneSubSRTCmd) String() (string, error) {
	return fmt.Sprintf(`rm "%s"`, cmd.srtFile), nil
}
