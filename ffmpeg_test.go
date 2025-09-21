package ffcmd_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/northbright/ffcmd"
	"github.com/northbright/timestamp"
)

func Example() {
	type Clip struct {
		File     string `json:"file"`
		Start    string `json:"start"`
		End      string `json:"end"`
		Subtitle string `json:"subtitle"`
	}

	type Output struct {
		File string `json:"file"`
		W    int    `json:"w"`
		H    int    `json:"h"`
	}

	type Highlights struct {
		Clips              []*Clip `json:"clips"`
		BGM                string  `json:"bgm"`
		BGMFadeOutDuration int     `json:"bgm_fade_out_duration"`
		Out                *Output `json:"output"`
	}

	jsonStr := `
{
    "clips": [
        {
		    "file": "01.MP4",
			"start": "",
			"end": "00:00:03",
			"subtitle": "Mido's tickling Mimao and he's enjoying..."
	    },
        {
            "file": "02.MOV",
            "start": "",
            "end": "",
            "subtitle": "Mimao's playing the toy."
        },
        {
            "file": "03.MOV",
			"start": "00:00:01",
			"end": "00:00:08",
            "subtitle": "It's hard to brush Maomi's teeth!"
        },
        {
            "file": "04.MOV",
            "start": "00:00:02",
            "end": "",
            "subtitle": "\"What's this??!!\"\nCan I eat it?"
        }
	],
	"bgm": "penguinmusic-Better Day.mp3",
	"bgm_fade_out_duration": 5,
	"output": { "file": "output.mp4", "w": 720, "h": 960 }
}
`

	hl := &Highlights{}

	if err := json.Unmarshal([]byte(jsonStr), hl); err != nil {
		log.Printf("json.Unmarshal() error: %v", err)
		return
	}

	// Creates a ffmpeg command.
	ffmpeg := ffcmd.NewFFmpegCmd("", "output.mp4", true, "-preset slow", "-shortest")
	// Create a "concat" filterchain.
	concatFC := ffcmd.NewFilterChain("[out_v]")

	// Segments count to concat.
	n := 0

	// Total duration var of all clips.
	totalDurationVarName := "total"
	ffmpeg.AddPreCmd(fmt.Sprintf("%s=0", totalDurationVarName))

	// Loop all video clips.
	for i, c := range hl.Clips {
		// Create clip video filter chain.
		clip_v := ffcmd.NewFilterChain(fmt.Sprintf("[clip_%02d_v]", i))

		// Add video file as ffmpeg input and get the input index.
		// Add video stream of the file([X:v:0], X is the ffmpeg input id) as clip's input.
		id := ffmpeg.AddInput(c.File)
		clip_v.AddInputByID(id, "v", 0)

		// Create and chain scale, pad, setsar filters.
		scale := fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", hl.Out.W, hl.Out.H)
		pad := fmt.Sprintf("pad=%d:%d:(ow-iw)/2:(oh-ih)/2", hl.Out.W, hl.Out.H)
		setsar := "setsar=1:1"

		clip_v.Chain(scale).Chain(pad).Chain(setsar)

		// Check if need to chain trim, setpts filter.
		if c.Start != c.End {
			// Create clip video filters.
			trim := "trim="

			if c.Start != "" {
				start, err := timestamp.New(c.Start)
				if err != nil {
					log.Printf("get start timestamp error: %v", err)
					return
				}
				trim += fmt.Sprintf("start=%s:", start.SecondStr())
			}

			if c.End != "" {
				end, err := timestamp.New(c.End)
				if err != nil {
					log.Printf("get end timestamp error: %v", err)
					return
				}
				trim += fmt.Sprintf("end=%s", end.SecondStr())
			} else {
				trim = strings.TrimSuffix(trim, ":")
			}

			setpts := "setpts=PTS-STARTPTS"

			// Chain trim and setpts filter.
			clip_v.Chain(trim).Chain(setpts)
		}

		// Get clip duration.
		clipDurationVarName := fmt.Sprintf("d%02d", i+1)

		// Create a command to get the duration and set the var's value to the duration.
		getDurationCmd, err := ffcmd.NewGetDurationCmd(c.File, c.Start, c.End, clipDurationVarName, "")
		if err != nil {
			log.Printf("ffcmd.NewGetDurationCmd() error: %v", err)
			return
		}
		// Add command to get the duration of the video clip as ffmpeg's pre-commands(setup commands).
		ffmpeg.AddPreCmd(getDurationCmd)

		// Add command to get total duration clips' duration.
		increaseCmd := fmt.Sprintf("%s=$(echo $%s + $%s | bc)", totalDurationVarName, totalDurationVarName, clipDurationVarName)
		ffmpeg.AddPreCmd(increaseCmd)

		// Check if need to chain subtitles filter.
		if c.Subtitle != "" {
			srtFile := strings.Replace(c.File, filepath.Ext(c.File), ".srt", -1)

			// Create a command to create the SRT file.
			createCmd, err := ffcmd.NewCreateOneSubSRTCmd(srtFile, c.Subtitle, clipDurationVarName)
			if err != nil {
				log.Printf("ffcmd.NewCreateOneSubSRTCmd() error: %v", err)
				return
			}

			// Add command to create SRT file as ffmpeg's pre-command(setup commmand).
			ffmpeg.AddPreCmd(createCmd)

			removeCmd, err := ffcmd.NewRemoveOneSubSRTCmd(srtFile)
			if err != nil {
				log.Printf("ffcmd.NewRemoveOneSubSRTCmd() error: %v", err)
				return
			}
			// Add command to remove created file as ffmpeg's post-command(clean-up command).
			ffmpeg.AddPostCmd(removeCmd)

			// Create and chain subtitles filter.
			subtitles := fmt.Sprintf("subtitles='%s':force_style='Fontsize=13'", srtFile)
			clip_v.Chain(subtitles)
		}

		// Add clip video filterchain to filtergraph.
		ffmpeg.Chain(clip_v)

		// Add clip video filter chain's output as concat filterchain's input.
		concatFC.AddInputByOutput(clip_v, 0)

		// Increase segment count.
		n += 1
	}

	// Create concat filters.
	concat := fmt.Sprintf("concat=n=%d:v=1:a=0", n)

	// Chain concat filter.
	concatFC.Chain(concat)

	// Add concat filterchain to filtergraph.
	ffmpeg.Chain(concatFC)

	// Select concat filterchain's output.
	ffmpeg.MapByOutput(concatFC, 0)

	// Get BGM fade out start time.
	getBGMFadeOutStartCmd := fmt.Sprintf("st=$(echo $%s - %d | bc) && st=$(printf \"%%.3f\" $st)", totalDurationVarName, hl.BGMFadeOutDuration)
	ffmpeg.AddPreCmd(getBGMFadeOutStartCmd)

	// Add BGM as input and create BGM filterchain.
	bgm := ffcmd.NewFilterChain("[bgm]")
	bgm.AddInputByID(ffmpeg.AddInput(hl.BGM), "a", 0)

	// Create afade filter.
	aFade := fmt.Sprintf("afade=t=out:st=$st:d=%d", hl.BGMFadeOutDuration)
	bgm.Chain(aFade)

	// Add BGM filterchain to the filtergraph.
	ffmpeg.Chain(bgm)

	// Select output of BGM filterchain.
	ffmpeg.MapByOutput(bgm, 0)

	// Output the raw FFmpeg string.
	str, err := ffmpeg.String()
	if err != nil {
		fmt.Printf("ffmpeg.String() error: %v", err)
		return
	}

	fmt.Println(str)

	// Press Ctrl + C to stop the test.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Create an exec.Command.
	cmd := exec.CommandContext(ctx, "bash", "-c", str)

	// To stop the process and its subprocesses,
	// request that the process group id be set (Setpgid: true) to the PID of the newly spawned process (Pgid: 0).
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Dir = "./examples"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Cancel = func() error {
		// Kill all processes in the group via `kill -9 -$PGID`.
		// Note the "-" to signal the group.
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	go func() {
		<-ctx.Done()
		log.Printf("context is done")
	}()

	if err = cmd.Run(); err != nil {
		log.Printf("cmd.Run() error: %v", err)
		return
	}

	log.Printf("cmd.Run() succeeded")

	// Output:
	// total=0 && d01=3.000 && total=$(echo $total + $d01 | bc) && d=$(printf "%.3f" $d01) && sec=$(echo $d | awk -F. '{ print $1 }') && frac=$(echo $d | awk -F. '{ print substr($2, 1, 3) }') && hh=$((sec / 3600)) && mm=$((sec % 3600 / 60)) && ss=$((sec % 3600 % 60)) && printf -v end "%02d:%02d:%02d,%03d" $hh $mm $ss $frac && printf "1\n00:00:00,000 --> %s\nMido's tickling Mimao and he's enjoying..." $end > "01.srt" && d=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "02.MOV"); d02=$(echo $d - 0.000 | bc) && total=$(echo $total + $d02 | bc) && d=$(printf "%.3f" $d02) && sec=$(echo $d | awk -F. '{ print $1 }') && frac=$(echo $d | awk -F. '{ print substr($2, 1, 3) }') && hh=$((sec / 3600)) && mm=$((sec % 3600 / 60)) && ss=$((sec % 3600 % 60)) && printf -v end "%02d:%02d:%02d,%03d" $hh $mm $ss $frac && printf "1\n00:00:00,000 --> %s\nMimao's playing the toy." $end > "02.srt" && d03=7.000 && total=$(echo $total + $d03 | bc) && d=$(printf "%.3f" $d03) && sec=$(echo $d | awk -F. '{ print $1 }') && frac=$(echo $d | awk -F. '{ print substr($2, 1, 3) }') && hh=$((sec / 3600)) && mm=$((sec % 3600 / 60)) && ss=$((sec % 3600 % 60)) && printf -v end "%02d:%02d:%02d,%03d" $hh $mm $ss $frac && printf "1\n00:00:00,000 --> %s\nIt's hard to brush Maomi's teeth!" $end > "03.srt" && d=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "04.MOV"); d04=$(echo $d - 2.000 | bc) && total=$(echo $total + $d04 | bc) && d=$(printf "%.3f" $d04) && sec=$(echo $d | awk -F. '{ print $1 }') && frac=$(echo $d | awk -F. '{ print substr($2, 1, 3) }') && hh=$((sec / 3600)) && mm=$((sec % 3600 / 60)) && ss=$((sec % 3600 % 60)) && printf -v end "%02d:%02d:%02d,%03d" $hh $mm $ss $frac && printf "1\n00:00:00,000 --> %s\n\"What's this??!!\"
	// Can I eat it?" $end > "04.srt" && st=$(echo $total - 5 | bc) && st=$(printf "%.3f" $st) && echo "y" | ffmpeg \
	// -i "01.MP4" \
	// -i "02.MOV" \
	// -i "03.MOV" \
	// -i "04.MOV" \
	// -i "penguinmusic-Better Day.mp3" \
	// -filter_complex " \
	// [0:v:0]scale=720:960:force_original_aspect_ratio=decrease,pad=720:960:(ow-iw)/2:(oh-ih)/2,setsar=1:1,trim=end=3.000,setpts=PTS-STARTPTS,subtitles='01.srt':force_style='Fontsize=13'[clip_00_v];
	// [1:v:0]scale=720:960:force_original_aspect_ratio=decrease,pad=720:960:(ow-iw)/2:(oh-ih)/2,setsar=1:1,subtitles='02.srt':force_style='Fontsize=13'[clip_01_v];
	// [2:v:0]scale=720:960:force_original_aspect_ratio=decrease,pad=720:960:(ow-iw)/2:(oh-ih)/2,setsar=1:1,trim=start=1.000:end=8.000,setpts=PTS-STARTPTS,subtitles='03.srt':force_style='Fontsize=13'[clip_02_v];
	// [3:v:0]scale=720:960:force_original_aspect_ratio=decrease,pad=720:960:(ow-iw)/2:(oh-ih)/2,setsar=1:1,trim=start=2.000,setpts=PTS-STARTPTS,subtitles='04.srt':force_style='Fontsize=13'[clip_03_v];
	// [clip_00_v][clip_01_v][clip_02_v][clip_03_v]concat=n=4:v=1:a=0[out_v];
	// [4:a:0]afade=t=out:st=$st:d=5[bgm]" \
	// -map "[bgm]" \
	// -map "[out_v]" \
	// -preset slow \
	// -shortest \
	// output.mp4 && rm "01.srt" && rm "02.srt" && rm "03.srt" && rm "04.srt"
}
