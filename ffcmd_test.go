package ffcmd_test

import (
	"fmt"

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

	// Create op video filterchain.
	op_v := ffcmd.NewFilterChain("[op_v]")

	// Add "opening.png" as ffmpeg input and get the input index.
	// Add video stream of "opening.png"([0:v:0]) as op video chain's input.
	op_v.AddInputByID(cmd.AddInput("opening.png"), "v", 0)

	// Create op video filters.
	fps := "fps=30"
	loop := "loop=loop=90:size=1"
	scale := "scale=1280:720:force_original_aspect_ratio=decrease"
	pad := "pad=1280:720:(ow-iw)/2:(oh-ih)/2"
	format := "format=pix_fmts=yuv420p"
	subtitles := "subtitles=opening.srt:force_style='Fontsize=16'"
	fade := "fade=t=out:st=2:d=1"

	// Chain op video filters.
	op_v.Chain(fps).Chain(loop).Chain(scale).Chain(pad).Chain(format).Chain(subtitles).Chain(fade)

	// Create op audio filterchain.
	op_a := ffcmd.NewFilterChain("[op_a]")

	// Create op audio fiters.
	aevalsrc := "aevalsrc=0:d=3"

	// Chain ed audio filters.
	op_a.Chain(aevalsrc)

	// Add op video / audio filterchain to filtergraph.
	cmd.Chain(op_v)
	cmd.Chain(op_a)

	// Create ed video filterchain.
	ed_v := ffcmd.NewFilterChain("[ed_v]")

	// Add "ending.JPG" as ffmpeg input and get the input index.
	// Add video stream of "ending.JPG"([1:v:0]) as ed's input.
	ed_v.AddInputByID(cmd.AddInput("ending.JPG"), "v", 0)

	// Create ed video filters.
	loop = "loop=loop=150:size=1"
	subtitles = "subtitles=ending.srt:force_style='Fontsize=16'"
	fade = "fade=t=out:st=4:d=1"

	// Chain ed video filters.
	ed_v.Chain(fps).Chain(loop).Chain(scale).Chain(pad).Chain(format).Chain(subtitles).Chain(fade)

	// Create audio filterchain.
	ed_a := ffcmd.NewFilterChain("[ed_a]")

	// Create ed audio fiters.
	aevalsrc = "aevalsrc=0:d=5"

	// Chain ed audio filters.
	ed_a.Chain(aevalsrc)

	// Add ed video / audio filterchain to filtergraph.
	cmd.Chain(ed_v)
	cmd.Chain(ed_a)

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
		id := cmd.AddInput(c.file)
		clip_v.AddInputByID(id, "v", 0)
		clip_a.AddInputByID(id, "a", 0)

		// Check if need to chain trim, setpts / atrim, asetpts filter.
		if c.start != c.end {
			// Create clip video filters.
			trim := fmt.Sprintf("trim=%d:%d", c.start, c.end)
			setpts := "setpts=PTS-STARTPTS"

			// Chain trim and setpts filter.
			clip_v.Chain(trim).Chain(setpts)

			// Create clip audio filters.
			atrim := fmt.Sprintf("atrim=%d:%d", c.start, c.end)
			asetpts := "asetpts=PTS-STARTPTS"

			// Chain atrim and asetpts filter.
			clip_a.Chain(atrim).Chain(asetpts)
		}

		// Check if need to chain subtitles filter.
		if c.subtitle != "" {
			sub, _ := ffcmd.NewSub(c.subtitle, "00:00:00,000", "00:00:07,000")
			srt, _ := ffcmd.NewSRT(c.file, sub)
			// Add command to create SRT file as ffmpeg's pre-commands(set-up commmands).
			cmd.AddPreCmd(srt.CreateCmd())
			// Add command to remove created file as ffmpeg's post-commands(clean-up commands).
			cmd.AddPostCmd(srt.RemoveCmd())

			subtitles := fmt.Sprintf("subtitles='%s'", srt.Filename())

			// Chain subtitles filter.
			clip_v.Chain(subtitles)
		}

		// Add clip video / audio filterchain to filtergraph.
		cmd.Chain(clip_v)
		cmd.Chain(clip_a)

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
	cmd.Chain(concatFC)

	// Add BGM as command input.
	id := cmd.AddInput("./bgm.m4a")

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
	cmd.Chain(bgmFC)

	// Select output streams.
	// If none stream is selected, it'll auto select last filterchain's labeled outputs.
	cmd.MapByOutput(concatFC, 0)
	cmd.MapByOutput(bgmFC, 0)

	str, err := cmd.String()
	if err != nil {
		fmt.Printf("cmd.String() error: %v", err)
		return
	}

	fmt.Println(str)

	// Output:
}
