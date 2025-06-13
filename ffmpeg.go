package ffcmd

import (
	"fmt"
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

// FFmpeg represents the ffmpeg command.
type FFmpeg struct {
	inputs          []string
	output          string
	fg              []*FilterChain
	selectedStreams map[string]struct{}
	preCmds         []Cmd
	postCmds        []Cmd
}

// New returns a new ffmpeg command.
func New(output string) *FFmpeg {
	return &FFmpeg{inputs: []string{}, output: output, fg: []*FilterChain{}, selectedStreams: make(map[string]struct{})}
}

// AddInput adds input and returns index of the input.
func (ff *FFmpeg) AddInput(in string) int {
	id := len(ff.inputs)
	ff.inputs = append(ff.inputs, in)
	return id
}

// AddPreCmd adds the command(set-up) to run before ffmpeg.
func (ff *FFmpeg) AddPreCmd(cmd Cmd) {
	ff.preCmds = append(ff.preCmds, cmd)
}

// AddPostCmd adds the command(clean-up) to run after ffmpeg.
func (ff *FFmpeg) AddPostCmd(cmd Cmd) {
	ff.postCmds = append(ff.postCmds, cmd)
}

// Chain chains filterchain and return a ffmpeg command to chain next filterchain.
// e.g. ff.Chain(videoFC).Chain(audioFC).Chain(ConcatFC).
func (ff *FFmpeg) Chain(fc *FilterChain) *FFmpeg {
	ff.fg = append(ff.fg, fc)
	return ff
}

// Map selects stream as ffmpeg output.
func (ff *FFmpeg) Map(stream string) {
	if _, ok := ff.selectedStreams[stream]; !ok {
		ff.selectedStreams[stream] = struct{}{}
	}
}

// MapByID selects stream by input index, stream type and index of stream as ffmpeg output.
func (ff *FFmpeg) MapByID(inputID int, streamType string, streamID int) {
	stream := fmt.Sprintf("[%d:%s:%d]", inputID, streamType, streamID)
	if _, ok := ff.selectedStreams[stream]; !ok {
		ff.selectedStreams[stream] = struct{}{}
	}
}

// MapByOutput selects the output stream of filterchain by index as ffmpeg output dynamically.
func (ff *FFmpeg) MapByOutput(fc *FilterChain, id int) {
	stream := fc.Output(id)
	ff.Map(stream)
}

// MapByOutputs selects all the output streams of filterchain as ffmpeg outputs dynamically.
func (ff *FFmpeg) MapByOutputs(fc *FilterChain) {
	for _, stream := range fc.Outputs() {
		ff.Map(stream)
	}
}

// String returns the ffmpeg command string to run.
func (ff *FFmpeg) String() (string, error) {
	str := ""
	for _, cmd := range ff.preCmds {
		s, err := cmd.String()
		if err != nil {
			return "", fmt.Errorf("add pre-cmd error: %v", err)
		}
		str += fmt.Sprintf(`%s && `, s)
	}

	str += "ffmpeg \\\n"

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
			ff.MapByOutputs(fc)
		}
	}

	str += "\" \\\n"

	for stream, _ := range ff.selectedStreams {
		str += fmt.Sprintf("-map \"%s\" \\\n", stream)
	}

	str += ff.output

	for _, cmd := range ff.postCmds {
		s, err := cmd.String()
		if err != nil {
			return "", fmt.Errorf("add post-cmd error: %v", err)
		}
		str += fmt.Sprintf(` && %s`, s)
	}

	return str, nil
}
