package ffcmd_test

import (
	"fmt"

	"github.com/northbright/ffcmd"
)

func ExampleTimestamp() {
	arr := []string{
		"00:00:00",
		"00:00:00.000",
		"01:02:03.456",
		"02:60:00",
		"03:00:60",
		"04:05:08,950",
		"100:59:59.500",
	}

	for _, str := range arr {
		ts, err := ffcmd.NewTimestamp(str)
		if err != nil {
			fmt.Printf("%s: create timestamp error: %s\n", str, err)
		} else {
			fmt.Printf("%s --> Second(): %.2f, String(): %s, SRTString(): %s\n", str, ts.Second(), ts.String(), ts.StringForSRT())
		}
	}

	// Output:
	// 00:00:00 --> Second(): 0.00, String(): 0:00:00, SRTString(): 0:00:00,000
	// 00:00:00.000 --> Second(): 0.00, String(): 0:00:00.000, SRTString(): 0:00:00,000
	// 01:02:03.456 --> Second(): 3723.46, String(): 1:02:03.456, SRTString(): 1:02:03,456
	// 02:60:00: create timestamp error: incorrect hhmmss input
	// 03:00:60: create timestamp error: incorrect hhmmss input
	// 04:05:08,950 --> Second(): 14708.95, String(): 4:05:08.950, SRTString(): 4:05:08,950
	// 100:59:59.500 --> Second(): 363599.50, String(): 100:59:59.500, SRTString(): 100:59:59,500
}
