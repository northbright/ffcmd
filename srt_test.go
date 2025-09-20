//go:build !windows

package ffcmd_test

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/northbright/ffcmd"
	"github.com/northbright/pathelper"
)

type ClipData struct {
	File     string
	Start    string
	End      string
	Subtitle string
}

var clipData = []ClipData{
	ClipData{File: "01.MP4", Start: "", End: "00:00:03", Subtitle: "Mido's tickling Mimao and he's enjoying..."},
	ClipData{File: "02.MOV", Start: "", End: "", Subtitle: "Mimao's playing the toy."},
	ClipData{File: "03.MOV", Start: "00:00:01", End: "00:00:08", Subtitle: "It's hard to brush Maomi's teeth!"},
	ClipData{File: "04.MOV", Start: "00:00:02", End: "", Subtitle: "\"What's this??!!\"\nCan I eat it?"},
}

func ExampleNewGetDurationCmd() {
	var cmds []string

	// Append the first command to the command slice.
	cmds = append(cmds, "total=0")

	for i, c := range clipData {
		// Create a var name.
		varName := fmt.Sprintf("d%02d", i+1)

		// Create a command to get the duration and set the var's value to the duration.
		getDurationCmd, err := ffcmd.NewGetDurationCmd(c.File, c.Start, c.End, varName, "")
		if err != nil {
			log.Printf("ffcmd.NewGetDurationCmd() error: %v", err)
			return
		}
		cmds = append(cmds, getDurationCmd)

		// Create a command to increase total by adding the duration.
		increaseCmd := fmt.Sprintf("total=$(echo $total + $%s | bc)", varName)
		cmds = append(cmds, increaseCmd)
	}

	// Create a command to show total duration.
	cmds = append(cmds, "printf '\\ntotal=%.3f' $total")

	// Concatenates all commands to a final command.
	finalCmd := ffcmd.ConcatCmds("", cmds...)

	fmt.Printf("final command:\n%s\n", finalCmd)

	// Create an exec.Command.
	cmd := exec.Command("bash", "-c", finalCmd)

	cmd.Dir = "examples"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("cmd.Run() error: %v", err)
		return
	}

	log.Printf("cmd.Run() succeeded")

	// Output:
	// final command:
	// total=0 && d01=3.000 && total=$(echo $total + $d01 | bc) && d=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "02.MOV"); d02=$(echo $d - 0.000 | bc) && total=$(echo $total + $d02 | bc) && d03=7.000 && total=$(echo $total + $d03 | bc) && d=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "04.MOV"); d04=$(echo $d - 2.000 | bc) && total=$(echo $total + $d04 | bc) && printf '\ntotal=%.3f' $total
	//
	// total=24.800
}

func ExampleNewCreateOneSubSRTCmd() {
	var cmds []string

	for i, c := range clipData {
		// Create a var name. The shell var is used to store the duration of the clip.
		varName := fmt.Sprintf("d%02d", i+1)

		// Create a command to get the duration and set the var's value to the duration.
		getDurationCmd, err := ffcmd.NewGetDurationCmd(c.File, c.Start, c.End, varName, "")
		if err != nil {
			log.Printf("ffcmd.NewGetDurationCmd() error: %v", err)
			return
		}
		cmds = append(cmds, getDurationCmd)

		// Create a command to create a SRT file which contains only 1 subtitle.
		srtFile := fmt.Sprintf("%s.srt", pathelper.BaseWithoutExt(c.File))
		srtCmd, err := ffcmd.NewCreateOneSubSRTCmd(srtFile, c.Subtitle, varName)
		if err != nil {
			log.Printf("ffcmd.NewCreateOneSubSRTCmd() error: %v", err)
			return
		}
		cmds = append(cmds, srtCmd)

		// Create a command to show the SRT file.
		cmds = append(cmds, fmt.Sprintf("cat %s && echo", srtFile))

		// Create a command to remove the SRT file after use.
		/*rmSRTCmd, err := ffcmd.NewRemoveOneSubSRTCmd(srtFile)
		if err != nil {
			log.Printf("ffcmd.NewRemoveOneSubSRTCmd() error: %v", err)
			return
		}
		cmds = append(cmds, rmSRTCmd)
		*/
	}

	// Concatenates all commands to a final command.
	finalCmd := ffcmd.ConcatCmds("", cmds...)

	fmt.Printf("final command:\n%s\n", finalCmd)

	// Create an exec.Command.
	cmd := exec.Command("bash", "-c", finalCmd)

	cmd.Dir = "examples"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("cmd.Run() error: %v", err)
		return
	}

	log.Printf("cmd.Run() succeeded")

	// Output:
	// final command:
	// d01=3.000 && d=$(printf "%.3f" $d01) && sec=$(echo $d | awk -F. '{ print $1 }') && frac=$(echo $d | awk -F. '{ print substr($2, 1, 3) }') && hh=$((sec / 3600)) && mm=$((sec % 3600 / 60)) && ss=$((sec % 3600 % 60)) && printf -v end "%02d:%02d:%02d,%03d" $hh $mm $ss $frac && printf "1\n00:00:00,000 --> %s\nMido's tickling Mimao and he's enjoying..." $end > "01.srt" && cat 01.srt && echo && rm "01.srt" && d=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "02.MOV"); d02=$(echo $d - 0.000 | bc) && d=$(printf "%.3f" $d02) && sec=$(echo $d | awk -F. '{ print $1 }') && frac=$(echo $d | awk -F. '{ print substr($2, 1, 3) }') && hh=$((sec / 3600)) && mm=$((sec % 3600 / 60)) && ss=$((sec % 3600 % 60)) && printf -v end "%02d:%02d:%02d,%03d" $hh $mm $ss $frac && printf "1\n00:00:00,000 --> %s\nMimao's playing the toy." $end > "02.srt" && cat 02.srt && echo && rm "02.srt" && d03=7.000 && d=$(printf "%.3f" $d03) && sec=$(echo $d | awk -F. '{ print $1 }') && frac=$(echo $d | awk -F. '{ print substr($2, 1, 3) }') && hh=$((sec / 3600)) && mm=$((sec % 3600 / 60)) && ss=$((sec % 3600 % 60)) && printf -v end "%02d:%02d:%02d,%03d" $hh $mm $ss $frac && printf "1\n00:00:00,000 --> %s\nIt's hard to brush Maomi's teeth!" $end > "03.srt" && cat 03.srt && echo && rm "03.srt" && d=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "04.MOV"); d04=$(echo $d - 2.000 | bc) && d=$(printf "%.3f" $d04) && sec=$(echo $d | awk -F. '{ print $1 }') && frac=$(echo $d | awk -F. '{ print substr($2, 1, 3) }') && hh=$((sec / 3600)) && mm=$((sec % 3600 / 60)) && ss=$((sec % 3600 % 60)) && printf -v end "%02d:%02d:%02d,%03d" $hh $mm $ss $frac && printf "1\n00:00:00,000 --> %s\n\"What's this??!!\"
	// Can I eat it?" $end > "04.srt" && cat 04.srt && echo && rm "04.srt"
	// 1
	// 00:00:00,000 --> 00:00:03,000
	// Mido's tickling Mimao and he's enjoying...
	// 1
	// 00:00:00,000 --> 00:00:09,633
	// Mimao's playing the toy.
	// 1
	// 00:00:00,000 --> 00:00:07,000
	// It's hard to brush Maomi's teeth!
	// 1
	// 00:00:00,000 --> 00:00:05,167
	// "What's this??!!"
	// Can I eat it?
}
