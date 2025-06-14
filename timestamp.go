package ffcmd

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Timestamp represents the timestamp for video and SRT file.
// It's 0-based.
type Timestamp struct {
	hh  int
	mm  int
	ss  int
	mmm int
}

// NewTimestamp parses str in the form of "HH:MM:SS(.mmm)" or "HH:MM:SS(,mmm)" and returns a new Timestamp.
func NewTimestamp(str string) (*Timestamp, error) {
	re := regexp.MustCompile(`^(\d{2}):([0-5][0-9]):([0-5][0-9])([\.|,](\d{3}))?$`)
	arr := re.FindStringSubmatch(str)
	l := len(arr)

	var hh, mm, ss, mmm int

	if l != 6 {
		return nil, fmt.Errorf("incorrect input")
	} else {
		hh, _ = strconv.Atoi(arr[1])
		mm, _ = strconv.Atoi(arr[2])
		ss, _ = strconv.Atoi(arr[3])

		if arr[5] != "" {
			mmm, _ = strconv.Atoi(arr[5])
		}

		return &Timestamp{hh: hh, mm: mm, ss: ss, mmm: mmm}, nil
	}
}

// NewTimestampFromSecond converts the seconds in float to timestamp.
func NewTimestampFromSecond(second float32) (*Timestamp, error) {
	integer, frac := math.Modf(float64(second))
	sec := int(integer)

	str := fmt.Sprintf("%.3f", frac)
	str = strings.TrimPrefix(str, "0.")
	mmm, _ := strconv.Atoi(str)

	hh := sec / 3600
	mm := sec / 3600 % 60
	ss := sec % 3600 % 60

	return &Timestamp{hh: hh, mm: mm, ss: ss, mmm: mmm}, nil
}

// Str returns the timestamp string.
// If forSRT is true, it returns string in the format: "HH:MM:SS,mmm".
// Otherwise, the format is "HH:MM:SS.mmm".
// mmm is the millisecond.
func (ts *Timestamp) Str(forSRT bool) string {
	sep := "."
	if forSRT {
		sep = ","
	}

	return fmt.Sprintf("%02d:%02d:%02d%s%03d", ts.hh, ts.mm, ts.ss, sep, ts.mmm)
}

// String returns the timestamp string in the format: "HH:MM:SS.mmm".
func (ts *Timestamp) String() string {
	return ts.Str(false)
}

// StringForSRT returns the timestamp string for SRT file which is in the format: "HH:MM:SS,mmm".
func (ts *Timestamp) StringForSRT() string {
	return ts.Str(true)
}

// Second returns the second string in the format: "s.mmm" format.
// mmm is the millisecond.
// It may used for the "start" / "end" option of "trim" filter.
func (ts *Timestamp) Second() string {
	return fmt.Sprintf("%d.%03d", ts.hh*3600+ts.mm*60+ts.ss, ts.mmm)
}
