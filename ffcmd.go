package ffcmd

import (
	"fmt"
)

type FilterChain struct {
	inputs  []any
	outputs []string
	filters []string
}

type filterChainOutputData struct {
	fc *FilterChain
	id int
}

func NewFilterChain(outputs ...string) *FilterChain {
	return &FilterChain{inputs: []any{}, outputs: outputs, filters: []string{}}
}

func (fc *FilterChain) AddInput(input string) {
	fc.inputs = append(fc.inputs, input)
}

func (fc *FilterChain) AddInputByID(inputID int, streamType string, streamID int) {
	fc.inputs = append(fc.inputs, fmt.Sprintf("[%d:%s:%d]", inputID, streamType, streamID))
}

func (fc *FilterChain) AddInputByOutput(fcOut *FilterChain, outputID int) {
	fc.inputs = append(fc.inputs, &filterChainOutputData{fcOut, outputID})
}

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

func (fc *FilterChain) Outputs() []string {
	if len(fc.filters) == 0 {
		return fc.Inputs()
	} else {
		return fc.outputs
	}
}

func (fc *FilterChain) Chain(filter string) *FilterChain {
	if filter != "" {
		fc.filters = append(fc.filters, filter)
	}
	return fc
}

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

type Cmd struct {
	inputs          []string
	output          string
	fg              []*FilterChain
	selectedStreams map[string]struct{}
}

func New(output string) *Cmd {
	return &Cmd{inputs: []string{}, output: output, fg: []*FilterChain{}, selectedStreams: make(map[string]struct{})}
}

// AddInput adds input and returns index of the input.
func (c *Cmd) AddInput(in string) int {
	id := len(c.inputs)
	c.inputs = append(c.inputs, in)
	return id
}

func (c *Cmd) Chain(fc *FilterChain) *Cmd {
	c.fg = append(c.fg, fc)
	return c
}

func (c *Cmd) Map(stream string) {
	if _, ok := c.selectedStreams[stream]; !ok {
		c.selectedStreams[stream] = struct{}{}
	}
}

func (c *Cmd) MapByID(inputID int, streamType string, streamID int) {
	stream := fmt.Sprintf("[%d:%s:%d]", inputID, streamType, streamID)
	if _, ok := c.selectedStreams[stream]; !ok {
		c.selectedStreams[stream] = struct{}{}
	}
}

func (c *Cmd) MapByOutput(fc *FilterChain, id int) {
	stream := fc.Output(id)
	c.Map(stream)
}

func (c *Cmd) MapByOutputs(fc *FilterChain) {
	for _, stream := range fc.Outputs() {
		c.Map(stream)
	}
}

func (c *Cmd) String() (string, error) {
	str := "ffmpeg \\\n"

	for _, in := range c.inputs {
		str += fmt.Sprintf("-i \"%s\" \\\n", in)
	}

	str += "-filter_complex \" \\\n"

	l := len(c.fg)
	for i, fc := range c.fg {
		s := fc.String()
		if s == "" {
			continue
		}

		str += s
		if i < l-1 {
			str += ";\n"
		} else {
			// Complex filtergraph outputs streams with labeled pads must be mapped once and exactly once.
			c.MapByOutputs(fc)
		}
	}

	str += "\" \\\n"

	for stream, _ := range c.selectedStreams {
		str += fmt.Sprintf("-map \"%s\" \\\n", stream)
	}

	str += c.output

	return str, nil
}
