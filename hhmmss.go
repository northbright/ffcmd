package ffcmd

import (
	"fmt"
	"regexp"
	"strconv"
)

func HHMMSSToSec(hhmmss string) (float32, error) {
	var sec float32

	re := regexp.MustCompile(`^(\d{2}):([0-5][0-9]):([0-5][0-9])(\.(\d+))?$`)
	arr := re.FindStringSubmatch(hhmmss)
	l := len(arr)

	if l != 6 {
		return 0, fmt.Errorf("incorrect hhmmss input")
	} else {
		hh := arr[1]
		mm := arr[2]
		ss := arr[3]

		hour, err := strconv.Atoi(hh)
		if err != nil {
			return 0, fmt.Errorf("failed to convert hh string to int")
		}

		minute, err := strconv.Atoi(mm)
		if err != nil {
			return 0, fmt.Errorf("failed to convert mm string to int")
		}

		second, err := strconv.Atoi(ss)
		if err != nil {
			return 0, fmt.Errorf("failed to convert ss string to int")
		}

		sec = float32(hour*3600 + minute*60 + second)

		if arr[5] != "" {
			millisecond, err := strconv.Atoi(arr[5])
			if err != nil {
				return 0, fmt.Errorf("failed to convert millisecond string to float")
			}
			sec += float32(millisecond) / 1000
		}

		return sec, nil
	}
}
