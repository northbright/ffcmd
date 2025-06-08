package ffcmd_test

import (
	"fmt"

	"github.com/northbright/ffcmd"
)

func ExampleHHMMSSToSec() {
	arr := []string{
		"00:00:00",
		"00:00:00.000",
		"01:02:03.456",
		"02:60:00",
		"03:00:60",
	}

	for _, hhmmss := range arr {
		sec, err := ffcmd.HHMMSSToSec(hhmmss)
		if err != nil {
			fmt.Printf("%s --> error: %v\n", hhmmss, err)
		} else {
			fmt.Printf("%s --> %.3f\n", hhmmss, sec)
		}
	}

	// Output:
	// 00:00:00 --> 0.000
	// 00:00:00.000 --> 0.000
	// 01:02:03.456 --> 3723.456
	// 02:60:00 --> error: incorrect hhmmss input
	// 03:00:60 --> error: incorrect hhmmss input
}
