package ffcmd

import (
	"fmt"
	"sort"
	"strings"
)

// FilterChain represents the filterchain of ffmpeg.
type FilterChain struct {
	inputs  []any
	outputs []string
	filters []string
}

// filterChainOutputData stores the filterchan and the output ID to generate output label as another filterchain's input.
type filterChainOutputData struct {
	fc *FilterChain
	id int
}

// NewFilterChain returns a filterchain by pre-defined outputs(labels) in the "[OUTPUT_LABEL]" format.
// e.g. concatFilterChain := ffcmd.NewFilterChain("[out_v]", "[out_a]")
func NewFilterChain(outputs ...string) *FilterChain {
	return &FilterChain{inputs: []any{}, outputs: outputs, filters: []string{}}
}

// AddInput adds raw string(label) as the input.
func (fc *FilterChain) AddInput(input string) {
	fc.inputs = append(fc.inputs, input)
}

// AddInputByID adds input by specifying input ID, stream type and ID in the stream(e.g. "[0:v:0]", [0:a:1]").
// inputID: 0-based input ID.
// streamType: "v" for video, "a" for audio, "s" for subtitles, "d" for data. See ffmpeg's doc for more.
// streamID: index of the stream for the stream type.
func (fc *FilterChain) AddInputByID(inputID int, streamType string, streamID int) {
	fc.inputs = append(fc.inputs, fmt.Sprintf("[%d:%s:%d]", inputID, streamType, streamID))
}

// AddInputByOutput adds another filterchain's output as input.
// It's useful when another filterchain's output is generated dynamically.
func (fc *FilterChain) AddInputByOutput(fcOut *FilterChain, outputID int) {
	fc.inputs = append(fc.inputs, &filterChainOutputData{fcOut, outputID})
}

// Input returns the input string by 0-based index.
func (fc *FilterChain) Input(id int) string {
	if id < 0 || id >= len(fc.inputs) {
		return ""
	}

	switch vv := fc.inputs[id].(type) {
	case string:
		return vv
	case *filterChainOutputData:
		fc := vv.fc
		id := vv.id
		return fc.Output(id)
	default:
		return ""
	}
}

// Inputs returns all inputs.
func (fc *FilterChain) Inputs() []string {
	var inputs []string

	for _, in := range fc.inputs {
		switch vv := in.(type) {
		case string:
			inputs = append(inputs, vv)
		case *filterChainOutputData:
			fc := vv.fc
			id := vv.id
			input := fc.Output(id)
			inputs = append(inputs, input)
		default:
		}
	}
	return inputs
}

// Output returns the output by 0-based index.
func (fc *FilterChain) Output(id int) string {
	if len(fc.filters) == 0 {
		return fc.Input(id)
	} else {
		if id < 0 || id >= len(fc.outputs) {
			return ""
		}
		return fc.outputs[id]
	}
}

// Outputs returns all outputs.
func (fc *FilterChain) Outputs() []string {
	if len(fc.filters) == 0 {
		return fc.Inputs()
	} else {
		return fc.outputs
	}
}

// Chain chains filter and returns a filterchain to chain next filter(e.g. fc.Chain("fps=30").Chain("scale=1280:720"))
func (fc *FilterChain) Chain(filter string) *FilterChain {
	if filter != "" {
		fc.filters = append(fc.filters, filter)
	}
	return fc
}

// String returns the filterchain's string for ffmpeg command.
func (fc *FilterChain) String() string {
	l := len(fc.filters)

	if l == 0 {
		// No filter in the chain, just return empty string as do nothing in the chain.
		return ""
	}

	str := ""
	for _, input := range fc.Inputs() {
		str += input
	}

	for i, filter := range fc.filters {
		str += filter
		if i < l-1 {
			str += ","
		}
	}

	for _, output := range fc.Outputs() {
		str += output
	}
	return str
}

// FFmpegCmd represents the ffmpeg command.
type FFmpegCmd struct {
	ffmpegPath      string
	inputs          []string
	output          string
	fg              []*FilterChain
	selectedStreams map[string]struct{}
	preCmds         []string
	postCmds        []string
	overwrite       bool
	options         []string
}

// NewFFmpegCmd returns a new ffmpeg command.
// ffmpegPath: path of ffmpeg binary. It'll be set to "ffmpeg" if it's empty.
// output: ffmpeg output(e.g. "output.mp4")
// overwrite: if overwrite output when run ffmpeg command.
// It'll failed to generate output if output exists and overwrite is set to false.
// options: optional options for ffmpeg(e.g. "-r 30", "-shortest", "-crf 20").
func NewFFmpegCmd(ffmpegPath string, output string, overwrite bool, options ...string) *FFmpegCmd {
	ff := &FFmpegCmd{
		ffmpegPath:      ffmpegPath,
		inputs:          []string{},
		output:          output,
		fg:              []*FilterChain{},
		selectedStreams: make(map[string]struct{}),
		overwrite:       overwrite,
		options:         options,
	}

	return ff
}

// AddInput adds input and returns index of the input.
func (ff *FFmpegCmd) AddInput(in string) int {
	id := len(ff.inputs)
	ff.inputs = append(ff.inputs, in)
	return id
}

// AddPreCmd adds the command(set-up) to run before ffmpeg.
func (ff *FFmpegCmd) AddPreCmd(cmd string) {
	ff.preCmds = append(ff.preCmds, cmd)
}

// AddPostCmd adds the command(clean-up) to run after ffmpeg.
func (ff *FFmpegCmd) AddPostCmd(cmd string) {
	ff.postCmds = append(ff.postCmds, cmd)
}

// Chain chains filterchain and return a ffmpeg command to chain next filterchain.
// e.g. ff.Chain(videoFC).Chain(audioFC).Chain(ConcatFC).
func (ff *FFmpegCmd) Chain(fc *FilterChain) *FFmpegCmd {
	ff.fg = append(ff.fg, fc)
	return ff
}

// Map selects stream as ffmpeg output.
func (ff *FFmpegCmd) Map(stream string) {
	if _, ok := ff.selectedStreams[stream]; !ok {
		ff.selectedStreams[stream] = struct{}{}
	}
}

// MapID selects a stream.
// inputID: stream index.
// streamType: stream type(e.g. "v" for video stream, "a" for audio stream).
// streadmID: index of stream.
func (ff *FFmpegCmd) MapID(inputID int, streamType string, streamID int) {
	stream := fmt.Sprintf("%d:%s:%d", inputID, streamType, streamID)
	if _, ok := ff.selectedStreams[stream]; !ok {
		ff.selectedStreams[stream] = struct{}{}
	}
}

// MapOutput selects an output stream of a filterchain.
// fc: filerchain to select.
// id: stream id of output streams of the filterchain.
func (ff *FFmpegCmd) MapOutput(fc *FilterChain, id int) {
	stream := fc.Output(id)

	// If the output is not a label but an input stream with index(e.g. [0:v:0]),
	// remove the "[" and "]" or -map will fail.
	if strings.Contains(stream, ":") {
		stream = strings.ReplaceAll(stream, "[", "")
		stream = strings.ReplaceAll(stream, "]", "")
	}

	ff.Map(stream)
}

// MapOutputs selects all output streams of a filterchain.
// fc: filerchain to select.
func (ff *FFmpegCmd) MapOutputs(fc *FilterChain) {
	for _, stream := range fc.Outputs() {
		ff.Map(stream)
	}
}

// String returns the ffmpeg command string to run.
func (ff *FFmpegCmd) String() (string, error) {
	str := ""

	// Add the commands to run before ffmpeg.
	if len(ff.preCmds) > 0 {
		str = ConcatCmds("", ff.preCmds...)
		str += " && "
	}

	// Check if overwrite output.
	if ff.overwrite {
		str += `echo "y" | `
	}

	bin := "ffmpeg"
	if ff.ffmpegPath != "" {
		bin = ff.ffmpegPath
	}

	str += fmt.Sprintf("%s \\\n", bin)

	for _, in := range ff.inputs {
		str += fmt.Sprintf("-i \"%s\" \\\n", in)
	}

	str += "-filter_complex \" \\\n"

	l := len(ff.fg)
	for i, fc := range ff.fg {
		s := fc.String()
		if s == "" {
			continue
		}

		str += s
		if i < l-1 {
			str += ";\n"
		} else {
			// Complex filtergraph outputs streams with labeled pads must be mapped once and exactly once.
			ff.MapOutputs(fc)
		}
	}

	str += "\" \\\n"

	// Sort streams by names.
	var selectedStreams []string
	for stream := range ff.selectedStreams {
		selectedStreams = append(selectedStreams, stream)
	}
	sort.Strings(selectedStreams)

	for _, stream := range selectedStreams {
		str += fmt.Sprintf("-map \"%s\" \\\n", stream)
	}

	// Add options.
	for _, option := range ff.options {
		str += fmt.Sprintf("%s \\\n", option)
	}

	// Add output.
	str += ff.output

	// Add the commands to run after ffmpeg.
	if len(ff.postCmds) > 0 {
		str = ConcatCmds(str, ff.postCmds...)
	}

	return str, nil
}
