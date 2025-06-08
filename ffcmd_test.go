package ffcmd_test

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/northbright/ffcmd"
)

func Example() {
	type clip struct {
		file     string
		start    int
		end      int
		subtitle string
	}

	clips := []clip{
		{file: "01.MOV", start: 8, end: 22, subtitle: "乌龟灵活跑位后接老汤助攻破门"},
		{file: "02.MOV", start: 0, end: 7, subtitle: "老汤助攻乌龟门前推射"},
	}

	// Create ffmpeg command with output file.
	cmd := ffcmd.New("output.mp4")

	// Create op video / audio filterchain and add them to filtergraph.
	op_v := ffcmd.NewFilterChain("[op_v]")
	op_a := ffcmd.NewFilterChain("[op_a]")
	cmd.Chain(op_v)
	cmd.Chain(op_a)

	// Add "opening.png" as ffmpeg input and get the input index.
	// Add video stream of "opening.png"([0:v:0]) as op video chain's input.
	op_v.AddInputByID(cmd.AddInput("opening.png"), "v", 0)

	// Create op video filters.
	fps := ffcmd.NewFilter("fps").Option("fps", 30)
	loop := ffcmd.NewFilter("loop").Option("loop", 90).Option("size", 1)
	scale := ffcmd.NewFilter("scale").Option("w", 1280).Option("h", 720).Option("force_original_aspect_ratio", "decrease")
	pad := ffcmd.NewFilter("pad").Option("w", 1280).Option("h", 720).Option("x", "(ow-iw)/2").Option("y", "(oh-ih)/2")
	format := ffcmd.NewFilter("format").Option("pix_fmts", "yuv420p")
	subtitles := ffcmd.NewFilter("subtitles").Option("filename", "opening.srt").Option("force_style", "'Fontsize=16'")
	fade := ffcmd.NewFilter("fade").Option("t", "out").Option("st", 2).Option("d", 1)

	// Chain op video filters.
	op_v.Chain(fps).Chain(loop).Chain(scale).Chain(pad).Chain(format).Chain(subtitles).Chain(fade)

	// Create op audio fiters.
	aevalsrc := ffcmd.NewFilter("aevalsrc").Option("exprs", 0).Option("d", 3)

	// Chain ed audio filters.
	op_a.Chain(aevalsrc)

	// Create ed video / audio filterchain and add them to filtergraph.
	ed_v := ffcmd.NewFilterChain("[ed_v]")
	ed_a := ffcmd.NewFilterChain("[ed_a]")
	cmd.Chain(ed_v)
	cmd.Chain(ed_a)

	// Add "ending.JPG" as ffmpeg input and get the input index.
	// Add video stream of "ending.JPG"([1:v:0]) as ed's input.
	ed_v.AddInputByID(cmd.AddInput("ending.JPG"), "v", 0)

	// Create ed video filters.
	loop = ffcmd.NewFilter("loop").Option("loop", 150).Option("size", 1)
	subtitles = ffcmd.NewFilter("subtitles").Option("filename", "ending.srt").Option("force_style", "'Fontsize=16'")
	fade = ffcmd.NewFilter("fade").Option("t", "out").Option("st", 4).Option("d", 1)

	// Chain ed video filters.
	ed_v.Chain(fps).Chain(loop).Chain(scale).Chain(pad).Chain(format).Chain(subtitles).Chain(fade)

	// Create ed audio fiters.
	aevalsrc = ffcmd.NewFilter("aevalsrc").Option("exprs", 0).Option("d", 5)

	// Chain ed audio filters.
	ed_a.Chain(aevalsrc)

	// Create concat filter chain.
	concatFC := ffcmd.NewFilterChain("[outv][outa]")

	// Add op video and audio filterchain's output as concat filterchain's input.
	concatFC.AddInputByOutput(op_v)
	concatFC.AddInputByOutput(op_a)

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
		id := cmd.AddInput(c.file)
		clip_v.AddInputByID(id, "v", 0)
		clip_a.AddInputByID(id, "a", 0)

		// Add clip video / audio filterchain to filtergraph.
		cmd.Chain(clip_v)
		cmd.Chain(clip_a)

		// Add clip video / audio filter chain's output as concat filterchain's input.
		concatFC.AddInputByOutput(clip_v)
		concatFC.AddInputByOutput(clip_a)

		// Increase segment count.
		n += 1

		// Check if need to chain trim, setpts / atrim, asetpts filter.
		if c.start != c.end {
			// Create clip video filters.
			trim := ffcmd.NewFilter("trim").Option("start", c.start).Option("end", c.end)
			setpts := ffcmd.NewFilter("setpts").Option("expr", "PTS-STARTPTS")

			// Chain trim and setpts filter.
			clip_v.Chain(trim).Chain(setpts)

			// Create clip audio filters.
			atrim := ffcmd.NewFilter("atrim").Option("start", c.start).Option("end", c.end)
			asetpts := ffcmd.NewFilter("asetpts").Option("expr", "PTS-STARTPTS")

			// Chain atrim and asetpts filter.
			clip_a.Chain(atrim).Chain(asetpts)
		}

		// Check if need to chain subtitles filter.
		if c.subtitle != "" {
			srtFile := strings.Replace(c.file, filepath.Ext(c.file), ".srt", -1)
			subtitles := ffcmd.NewFilter("subtitles").Option("filename", srtFile)

			// Chain subtitles filter.
			clip_v.Chain(subtitles)
		}
	}

	// Add ed video and audio filterchain's output as concat filterchain's input.
	concatFC.AddInputByOutput(ed_v)
	concatFC.AddInputByOutput(ed_a)

	// Create concat filters.
	concat := ffcmd.NewFilter("concat").Option("n", n).Option("v", 1).Option("a", 1)

	// Chain concat filters.
	concatFC.Chain(concat)

	// Add concat filterchain to filtergraph.
	cmd.Chain(concatFC)

	str, err := cmd.String()
	if err != nil {
		fmt.Printf("cmd.String() error: %v", err)
		return
	}

	fmt.Println(str)

	// Output:
}
