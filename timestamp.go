package ffcmd

import (
	"fmt"
	"regexp"
	"strconv"
)

type Timestamp struct {
	hh   int
	mm   int
	ss   int
	msec string
}

func NewTimestamp(str string) (*Timestamp, error) {
	hh, mm, ss, msec, err := ParseHHMMSS(str)
	if err != nil {
		return nil, err
	}

	return &Timestamp{hh: hh, mm: mm, ss: ss, msec: msec}, nil
}

func (ts *Timestamp) Second() float32 {
	sec := ts.hh*3600 + ts.mm*60 + ts.ss

	str := fmt.Sprintf("%d.%s", sec, ts.msec)
	second, _ := strconv.ParseFloat(str, 32)
	return float32(second)
}

func (ts *Timestamp) String() string {
	str := fmt.Sprintf("%d:%02d:%02d", ts.hh, ts.mm, ts.ss)

	if ts.msec != "" {
		str += fmt.Sprintf(".%s", ts.msec)
	}

	return str
}

func (ts *Timestamp) StringForSRT() string {
	str := fmt.Sprintf("%d:%02d:%02d", ts.hh, ts.mm, ts.ss)

	if ts.msec != "" {
		str += fmt.Sprintf(",%s", ts.msec)
	} else {
		str += ",000"
	}

	return str
}

// ParseHHMMSS returns hour, minute, second and millisecond by given string in the form of "HH:MM:SS.ms" or "HH:MM:SS,ms".
func ParseHHMMSS(hhmmss string) (hh, mm, ss int, msec string, err error) {
	re := regexp.MustCompile(`^(\d+):([0-5][0-9]):([0-5][0-9])([\.|,](\d+){1,3})?$`)
	arr := re.FindStringSubmatch(hhmmss)
	l := len(arr)

	if l != 6 {
		return 0, 0, 0, "", fmt.Errorf("incorrect hhmmss input")
	} else {
		hh, err = strconv.Atoi(arr[1])
		if err != nil {
			return 0, 0, 0, "", fmt.Errorf("failed to convert hh string to int")
		}

		mm, err = strconv.Atoi(arr[2])
		if err != nil {
			return 0, 0, 0, "", fmt.Errorf("failed to convert mm string to int")
		}

		ss, err = strconv.Atoi(arr[3])
		if err != nil {
			return 0, 0, 0, "", fmt.Errorf("failed to convert ss string to int")
		}

		if arr[5] != "" {
			msec = arr[5]
		}

		return hh, mm, ss, msec, nil
	}
}
