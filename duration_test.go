package ffcmd_test

/*
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

func Example_NewGetDurationCmd() {
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
*/
