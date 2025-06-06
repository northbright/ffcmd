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
		Start    int
		End      int
		subtitle string
	}

	clips := []clip{
		{file: "01.MOV", trimStart: 10, trimEnd: 0, subtitle: "Goal from Jacky!"},
		{file: "02.MOV", trimStart: 0, trimEnd: 0, subtitle: "Goal from Sonny!"},
	}

	cmd := ffcmd.New("output.mp4")
	cmd.AddInput("opening.png")

	op_v := ffcmd.NewFilterChain(func() string { return "[op_v]" })
	op_v.AddStreamInput(cmd.AddInput("opening.png"), "v", 0)

	ed_v := ffcmd.NewFilterChain(func() string { return "[ed_v]" })
	ed_v.AddStreamInput(cmd.AddInput("ending.JPG"), "v", 0)

	cmd.Chain(op_v)

	for i, c := range clips {
		if c.Start != 0 && c.End != 0 {
			// Video filter chain
			clip_v := ffcmd.NewFilterChain(func() string { return fmt.Sprintf("[clip_%02d_v]", i) })
			clip_v.AddStreamInput(cmd.AddInput(file), "v", 0)
			cmd.Chain(clip_v)

			trim := ffcmd.NewFilter("trim").Option("start", c.Start).Option("end", c.End)
			setpts := ffcmd.NewFilter("setpts").Option("expr", "PTS-STARTPTS")

			clip_v.Chain(trim).Chain(setpts)

			if c.subtitle != "" {
				srtFile := strings.Replace(c.file, filepath.Ext(c.file), ".srt", -1)
				subtitles := ffcmd.NewFilter("subtitles").Option("file", srtFile)
				clip_v.Chain(subtitles)
			}

			// Audio filter chain
			clip_a := ffcmd.NewFilterChain(func() string { return fmt.Sprintf("[clip_%02d_a]", i) })
			clip_v.AddStreamInput(cmd.AddInput(file), "a", 0)
			cmd.Chain(clip_a)

			trim := ffcmd.NewFilter("atrim").Option("start", c.Start).Option("end", c.End)
			setpts := ffcmd.NewFilter("asetpts").Option("expr", "PTS-STARTPTS")

			clip_a.Chain(trim).Chain(setpts)
		} else {
			if c.subtitle != "" {
				clip_v := ffcmd.NewFilterChain(func() string { return fmt.Sprintf("[clip_%02d_v]", i) })
				clip_v.AddStreamInput(cmd.AddInput(file), "v", 0)
				cmd.Chain(clip_v)

				srtFile := strings.Replace(c.file, filepath.Ext(c.file), ".srt", -1)
				subtitles := ffcmd.NewFilter("subtitles").Option("file", srtFile)
				clip_v.Chain(subtitles)
			}
		}
	}

	cmd.Chain(ed_v)

	fps := ffcmd.NewFilter("fps").Option("fps", 30)
	loop := ffcmd.NewFilter("loop").Option("loop", 90).Option("size", 1)
	scale := ffcmd.NewFilter("scale").Option("w", 1280).Option("h", 720).Option("force_original_aspect_ratio", "decrease")
	pad := ffcmd.NewFilter("pad").Option("w", 1280).Option("h", 720).Option("x", "(ow-iw)/2").Option("y", "(oh-ih)/2")
	format := ffcmd.NewFilter("format").Option("pix_fmts", "yuv420p")
	subtitles := ffcmd.NewFilter("subtitles").Option("filename", "1.srt").Option("force_style", "'Fontsize=16'")
	fade := ffcmd.NewFilter("fade").Option("t", "out").Option("st", 2).Option("d", 1)

	op_v.Chain(fps).Chain(loop).Chain(scale).Chain(pad).Chain(format).Chain(subtitles).Chain(fade)

	str := op_v.String()
	fmt.Println(str)

	cmd.Chain(op_v)

	// Output:
}
