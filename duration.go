package ffcmd

import (
	"fmt"
	"strings"

	"github.com/northbright/timestamp"
)

// NewGetDurationCmd gets the duration of a media file(or clip) and store the duration into a shell variable.
// file: file name of a media file.
// start: start time of the clip. Format: "HH:MM:SS.mmm".
// end: end time of the clip.
// If start or end is not empy it'll get the duration of the clip.
// varName: shell variable name to store the duration.
// ffprobePath: path of ffprobe. It'll be set to "ffprobe" if it's empty.
func NewGetDurationCmd(file, start, end, varName, ffprobePath string) (string, error) {
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

	if start == "" {
		start = "00:00:00.000"
	}

	if tsStart, err = timestamp.New(start); err != nil {
		return "", fmt.Errorf("invalid start time format")
	}

	startSec = tsStart.Second()

	if end == "" {
		if file == "" {
			return "", fmt.Errorf("both end time and video filename are empty, can not get end timestamp")
		}

		str = fmt.Sprintf(`d=$(%s -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "%s"); %s=$(echo $d - %.03f | bc)`, ffprobePath, file, varName, startSec)
	} else {
		if tsEnd, err = timestamp.New(end); err != nil {
			return "", fmt.Errorf("invalid end time format")
		}

		if tsEnd, err = tsEnd.Sub(tsStart); err != nil {
			return "", fmt.Errorf("tsEnd.Sub() error: %v", err)
		}

		str = fmt.Sprintf(`%s=%.03f`, varName, tsEnd.Second())
	}

	return str, nil
}

// NewDurationToTimestampCmd converts a duration to a timestamp.
// It reads the duration in the shell variable and store the converted timestamp in another shell variable.
// durationVarName: name of shell variable stored duration.
// timestampVarName: name of shell variable to store timestamp.
func NewDurationToTimestampCmd(durationVarName, timestampVarName string) string {
	return fmt.Sprintf(`d=$(printf "%%.3f" $%s) && sec=$(echo $d | awk -F. '{ print $1 }') && frac=$(echo $d | awk -F. '{ print substr($2, 1, 3) }') && hh=$((sec / 3600)) && mm=$((sec %% 3600 / 60)) && ss=$((sec %% 3600 %% 60)) && printf -v %s "%%02d:%%02d:%%02d,%%03d" $hh $mm $ss $frac`, durationVarName, timestampVarName)
}
