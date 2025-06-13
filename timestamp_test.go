package ffcmd_test

import (
	"fmt"
	"log"

	"github.com/northbright/ffcmd"
)

func ExampleTimestamp() {
	arr := []string{
		"00:00:00",
		"10:20:30.500",
		"20:30:40,900",
	}

	for _, str := range arr {
		ts, err := ffcmd.NewTimestamp(str)
		if err != nil {
			log.Printf("NewTimestamp() error: %v", err)
			return
		}
		fmt.Printf("%s -> String(): %s, StringForSRT(): %s, Second(): %s\n", str, ts.String(), ts.StringForSRT(), ts.Second())
	}

	fArr := []float32{
		0.0,
		3.14,
		3882.46000,
	}

	for _, f := range fArr {
		ts, err := ffcmd.NewTimestampFromSecond(f)
		if err != nil {
			log.Printf("NewTimestampFromSecond() error: %v", err)
			return
		}
		fmt.Printf("%f -> String(): %s, StringForSRT(): %s, Second(): %s\n", f, ts.String(), ts.StringForSRT(), ts.Second())
	}

	// Output:
	// 00:00:00 -> String(): 00:00:00.000, StringForSRT(): 00:00:00,000, Second(): 0.000
	// 10:20:30.500 -> String(): 10:20:30.500, StringForSRT(): 10:20:30,500, Second(): 37230.500
	// 20:30:40,900 -> String(): 20:30:40.900, StringForSRT(): 20:30:40,900, Second(): 73840.900
	// 0.000000 -> String(): 00:00:00.000, StringForSRT(): 00:00:00,000, Second(): 0.000
	// 3.140000 -> String(): 00:00:03.140, StringForSRT(): 00:00:03,140, Second(): 3.140
	// 3882.459961 -> String(): 01:01:42.460, StringForSRT(): 01:01:42,460, Second(): 3702.460
}
