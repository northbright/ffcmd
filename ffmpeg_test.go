package ffcmd_test

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/northbright/ffcmd"
)

func Example() {
	type Clip struct {
		File     string
		Start    string
		End      string
		Subtitle string
	}

	type ImageClip struct {
		File            string
		Duration        int
		FadeOutDuration int
		Subtitle        string
	}

	type Output struct {
		File string
		W    int
		H    int
		FPS  int
	}

	op := ImageClip{
		File:            "op.jpg",
		Duration:        3,
		FadeOutDuration: 1,
		Subtitle:        "Good Times with Maomi & Mimao",
	}

	ed := ImageClip{
		File:            "ed.jpg",
		Duration:        3,
		FadeOutDuration: 1,
		Subtitle:        "Mimao likes lying on father's bed...😂",
	}

	clips := []Clip{
		{File: "01.MP4", Start: "", End: "00:00:05", Subtitle: "Mido's tickling Mimao and she's enjoying..."},
		{File: "02.MOV", Start: "", End: "", Subtitle: "Mimao's playing the toy."},
		{File: "03.MOV", Start: "00:00:01", End: "00:00:09", Subtitle: "It's hard to brush Maomi's teeth..."},
	}

	out := Output{
		File: "output.mp4",
		W:    720,
		H:    960,
		FPS:  30,
	}

	// Create ffmpeg command with output file.
	ffmpeg := ffcmd.New("output.mp4")

	// Create op video filterchain.
	op_v := ffcmd.NewFilterChain("[op_v]")

	// Add "op.jpg" as ffmpeg input and get the input index.
	// Add video stream of "op.jpg"([0:v:0]) as op video chain's input.
	op_v.AddInputByID(ffmpeg.AddInput(op.File), "v", 0)

	// Create op video filters.
	fps := fmt.Sprintf("fps=%d", out.FPS)
	loop := fmt.Sprintf("loop=loop=%d:size=1", op.Duration*out.FPS)
	scale := fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", out.W, out.H)
	pad := fmt.Sprintf("pad=%d:%d:(ow-iw)/2:(oh-ih)/2", out.W, out.H)
	setsar := "setsar=1:1"
	format := "format=pix_fmts=yuv420p"

	// Chain op video filters.
	op_v.Chain(fps).Chain(loop).Chain(scale).Chain(pad).Chain(setsar).Chain(format)

	// Check if need to chain subtitles filter.
	if op.Subtitle != "" {
		srtFile := strings.Replace(op.File, filepath.Ext(op.File), ".srt", -1)
		createCmd, err := ffcmd.NewCreateOneSubSRTCmdForImageClip(srtFile, op.Subtitle, float32(op.Duration))
		if err != nil {
			log.Printf("ffcmd.NewCreateOneSubSRTCmdForImageClip() error: %v", err)
			return
		}
		// Add command to create SRT file as ffmpeg's pre-commands(set-up commmands).
		ffmpeg.AddPreCmd(createCmd)

		removeCmd, err := ffcmd.NewRemoveOneSubSRTCmd(srtFile)
		if err != nil {
			log.Printf("ffcmd.NewRemoveOneSubSRTCmd() error: %v", err)
			return
		}
		// Add command to remove created file as ffmpeg's post-commands(clean-up commands).
		ffmpeg.AddPostCmd(removeCmd)

		subtitles := fmt.Sprintf("subtitles='%s'", srtFile)

		// Chain subtitles filter.
		op_v.Chain(subtitles)
	}

	// Chain fade filter.
	fade := fmt.Sprintf("fade=t=out:st=%d:d=%d", op.Duration-op.FadeOutDuration, op.FadeOutDuration)
	op_v.Chain(fade)

	// Create op audio filterchain.
	op_a := ffcmd.NewFilterChain("[op_a]")

	// Create op audio fiters.
	aevalsrc := fmt.Sprintf("aevalsrc=0:d=%d", op.Duration)

	// Chain ed audio filters.
	op_a.Chain(aevalsrc)

	// Add op video / audio filterchain to filtergraph.
	ffmpeg.Chain(op_v)
	ffmpeg.Chain(op_a)

	// Create ed video filterchain.
	ed_v := ffcmd.NewFilterChain("[ed_v]")

	// Add "ed.jpg" as ffmpeg input and get the input index.
	// Add video stream of "ed.jpg"([1:v:0]) as ed's input.
	ed_v.AddInputByID(ffmpeg.AddInput(ed.File), "v", 0)

	// Create ed video filters.
	loop = fmt.Sprintf("loop=loop=%d:size=1", ed.Duration*out.FPS)

	// Chain ed video filters.
	ed_v.Chain(fps).Chain(loop).Chain(scale).Chain(pad).Chain(setsar).Chain(format)

	// Check if need to chain subtitles filter.
	if ed.Subtitle != "" {
		srtFile := strings.Replace(ed.File, filepath.Ext(ed.File), ".srt", -1)
		createCmd, err := ffcmd.NewCreateOneSubSRTCmdForImageClip(srtFile, ed.Subtitle, float32(ed.Duration))
		if err != nil {
			log.Printf("ffcmd.NewCreateOneSubSRTCmdForImageClip() error: %v", err)
			return
		}
		// Add command to create SRT file as ffmpeg's pre-commands(set-up commmands).
		ffmpeg.AddPreCmd(createCmd)

		removeCmd, err := ffcmd.NewRemoveOneSubSRTCmd(srtFile)
		if err != nil {
			log.Printf("ffcmd.NewRemoveOneSubSRTCmd() error: %v", err)
			return
		}
		// Add command to remove created file as ffmpeg's post-commands(clean-up commands).
		ffmpeg.AddPostCmd(removeCmd)

		subtitles := fmt.Sprintf("subtitles='%s'", srtFile)

		// Chain subtitles filter.
		ed_v.Chain(subtitles)
	}

	// Chain fade filter.
	fade = fmt.Sprintf("fade=t=out:st=%d:d=%d", ed.Duration-ed.FadeOutDuration, ed.FadeOutDuration)
	ed_v.Chain(fade)

	// Create audio filterchain.
	ed_a := ffcmd.NewFilterChain("[ed_a]")

	// Create ed audio fiters.
	aevalsrc = fmt.Sprintf("aevalsrc=0:d=%d", ed.Duration)

	// Chain ed audio filters.
	ed_a.Chain(aevalsrc)

	// Add ed video / audio filterchain to filtergraph.
	ffmpeg.Chain(ed_v)
	ffmpeg.Chain(ed_a)

	// Create concat filter chain.
	concatFC := ffcmd.NewFilterChain("[outv]", "[outa]")

	// Add op video and audio filterchain's output as concat filterchain's input.
	concatFC.AddInputByOutput(op_v, 0)
	concatFC.AddInputByOutput(op_a, 0)

	// Segments count to concat.
	// Initialized to 2: op + ed.
	n := 2

	// Loop all video clips.
	for i, c := range clips {
		// Create clip video filter chain.
		clip_v := ffcmd.NewFilterChain(fmt.Sprintf("[clip_%02d_v]", i))

		// Create clip audio filter chain.
		clip_a := ffcmd.NewFilterChain(fmt.Sprintf("[clip_%02d_a]", i))

		// Add video file as ffmpeg input and get the input index.
		// Add video / audio stream of the file([X:v:0] / [X:a:0], X is the ffmpeg input id) as clip's input.
		id := ffmpeg.AddInput(c.File)
		clip_v.AddInputByID(id, "v", 0)
		clip_a.AddInputByID(id, "a", 0)

		// Create and chain scale, pad, setsar filters.
		scale := fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", out.W, out.H)
		pad := fmt.Sprintf("pad=%d:%d:(ow-iw)/2:(oh-ih)/2", out.W, out.H)
		setsar := "setsar=1:1"

		clip_v.Chain(scale).Chain(pad).Chain(setsar)

		// Check if need to chain trim, setpts / atrim, asetpts filter.
		if c.Start != c.End {
			// Create clip video / audio filters.
			trim := "trim="
			atrim := "atrim="

			if c.Start != "" {
				start, err := ffcmd.NewTimestamp(c.Start)
				if err != nil {
					log.Printf("get start timestamp error: %v", err)
					return
				}
				trim += fmt.Sprintf("start=%s:", start.Second())
				atrim += fmt.Sprintf("start=%s:", start.Second())
			}

			if c.End != "" {
				end, err := ffcmd.NewTimestamp(c.End)
				if err != nil {
					log.Printf("get end timestamp error: %v", err)
					return
				}
				trim += fmt.Sprintf("end=%s", end.Second())
				atrim += fmt.Sprintf("end=%s", end.Second())
			} else {
				trim = strings.TrimSuffix(trim, ":")
				atrim = strings.TrimSuffix(atrim, ":")
			}

			setpts := "setpts=PTS-STARTPTS"

			// Chain trim and setpts filter.
			clip_v.Chain(trim).Chain(setpts)

			asetpts := "asetpts=PTS-STARTPTS"

			// Chain atrim and asetpts filter.
			clip_a.Chain(atrim).Chain(asetpts)
		}

		// Check if need to chain subtitles filter.
		if c.Subtitle != "" {
			srtFile := strings.Replace(c.File, filepath.Ext(c.File), ".srt", -1)
			createCmd, err := ffcmd.NewCreateOneSubSRTCmd(srtFile, c.File, c.Subtitle, c.Start, c.End)
			if err != nil {
				log.Printf("ffcmd.NewCreateOneSubSRTCmd() error: %v", err)
				return
			}
			// Add command to create SRT file as ffmpeg's pre-commands(set-up commmands).
			ffmpeg.AddPreCmd(createCmd)

			removeCmd, err := ffcmd.NewRemoveOneSubSRTCmd(srtFile)
			if err != nil {
				log.Printf("ffcmd.NewRemoveOneSubSRTCmd() error: %v", err)
				return
			}
			// Add command to remove created file as ffmpeg's post-commands(clean-up commands).
			ffmpeg.AddPostCmd(removeCmd)

			subtitles := fmt.Sprintf("subtitles='%s'", srtFile)

			// Chain subtitles filter.
			clip_v.Chain(subtitles)
		}

		// Add clip video / audio filterchain to filtergraph.
		ffmpeg.Chain(clip_v)
		ffmpeg.Chain(clip_a)

		// Add clip video / audio filter chain's output as concat filterchain's input.
		concatFC.AddInputByOutput(clip_v, 0)
		concatFC.AddInputByOutput(clip_a, 0)

		// Increase segment count.
		n += 1
	}

	// Add ed video and audio filterchain's output as concat filterchain's input.
	concatFC.AddInputByOutput(ed_v, 0)
	concatFC.AddInputByOutput(ed_a, 0)

	// Create concat filters.
	concat := fmt.Sprintf("concat=n=%d:v=1:a=1", n)

	// Chain concat filters.
	concatFC.Chain(concat)

	// Add concat filterchain to filtergraph.
	ffmpeg.Chain(concatFC)

	// Add BGM as command input.
	id := ffmpeg.AddInput("bgm.m4a")

	// Create filterchain to merge BGM and original audio streams.
	bgmFC := ffcmd.NewFilterChain("[outa_merged_bgm]")
	bgmFC.AddInputByID(id, "a", 0)
	bgmFC.AddInputByOutput(concatFC, 1)

	// Create amerge filter.
	amerge := "amerge=inputs=2"

	// Create pan filter.
	pan := "pan=stereo|c0<c0+c2|c1<c1+c3"

	// Chain filters.
	bgmFC.Chain(amerge).Chain(pan)

	// Add BGM filterchain.
	ffmpeg.Chain(bgmFC)

	// Select output streams.
	// If none stream is selected, it'll auto select last filterchain's labeled outputs.
	ffmpeg.MapByOutput(concatFC, 0)
	ffmpeg.MapByOutput(bgmFC, 0)

	str, err := ffmpeg.String()
	if err != nil {
		fmt.Printf("cmd.String() error: %v", err)
		return
	}

	fmt.Println(str)

	// Output:
	// echo -ne "1\n00:00:08,000 --> 00:00:22,000\nGoal from Wugui\!" > "01.srt" && ffprobe -v error -select_streams v:0 -show_entries stream.Duration -of csv=s=,:p=0 "02.MOV" | awk -F. '{ print $1 }' | read sec; hh=$((sec / 3600)); mm=$((sec % 3600 / 60)); ss=$((sec % 3600 % 60)); printf -v end "%02d:%02d:%02d,000" hh mm ss; echo -ne "1\n00:00:00,000 --> $end\nSonny passed the ball to Wugui and Goal\!\!" > "02.srt" && ffmpeg \
	// -i "opening.png" \
	// -i "ending.JPG" \
	// -i "01.MOV" \
	// -i "02.MOV" \
	// -i "./bgm.m4a" \
	// -filter_complex " \
	// [0:v:0]fps=30,loop=loop=90:size=1,scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,format=pix_fmts=yuv420p,subtitles=opening.srt:force_style='Fontsize=16',fade=t=out:st=2:d=1[op_v];
	// aevalsrc=0:d=3[op_a];
	// [1:v:0]fps=30,loop=loop=150:size=1,scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2,format=pix_fmts=yuv420p,subtitles=ending.srt:force_style='Fontsize=16',fade=t=out:st=4:d=1[ed_v];
	// aevalsrc=0:d=5[ed_a];
	// [2:v:0]trim=start=8.000:end=22.000,setpts=PTS-STARTPTS,subtitles='01.srt'[clip_00_v];
	// [2:a:0]atrim=start=8.000:end=22.000,asetpts=PTS-STARTPTS[clip_00_a];
	// [3:v:0]subtitles='02.srt'[clip_01_v];
	// [op_v][op_a][clip_00_v][clip_00_a][clip_01_v][3:a:0][ed_v][ed_a]concat=n=4:v=1:a=1[outv][outa];
	// [4:a:0][outa]amerge=inputs=2,pan=stereo|c0<c0+c2|c1<c1+c3[outa_merged_bgm]" \
	// -map "[outv]" \
	// -map "[outa_merged_bgm]" \
	// output.mp4 && rm "01.srt" && rm "02.srt"
}
