package ffcmd

import (
	"fmt"
	"regexp"
	"strconv"
)

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

func (ts *Timestamp) String() string {
	return fmt.Sprintf("%02d:%02d:%02d,%03d", ts.hh, ts.mm, ts.ss, ts.mmm)
}

func (ts *Timestamp) Second() string {
	return fmt.Sprintf("%d.%03d", ts.hh*3600+ts.mm*60+ts.ss, ts.mmm)
}
