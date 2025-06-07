package ffcmd

import (
	"fmt"
	"regexp"
)

type option struct {
	k string
	v any
}

type Filter struct {
	name    string
	options []option
}

func NewFilter(name string) *Filter {
	return &Filter{name: name, options: []option{}}
}

func (f *Filter) Option(k string, v any) *Filter {
	f.options = append(f.options, option{k: k, v: v})
	return f
}

func (f *Filter) String() string {
	l := len(f.options)
	str := fmt.Sprintf("%s=", f.name)
	for i, option := range f.options {
		str += fmt.Sprintf("%s=%v", option.k, option.v)
		if i < l-1 {
			str += ":"
		}
	}
	return str
}

type FilterChain struct {
	inputs  []any
	output  string
	filters []*Filter
}

func NewFilterChain(output string) *FilterChain {
	return &FilterChain{inputs: []any{}, output: output}
}

func (fc *FilterChain) AddInput(input string) {
	fc.inputs = append(fc.inputs, input)
}

func (fc *FilterChain) AddStreamInput(inputID int, streamType string, streamID int) {
	fc.inputs = append(fc.inputs, fmt.Sprintf("[%d:%s:%d]", inputID, streamType, streamID))
}

func (fc *FilterChain) AddInputFromOutput(fcOut *FilterChain) {
	fc.inputs = append(fc.inputs, fcOut)
}

func (fc *FilterChain) Inputs() string {
	inputs := ""
	for _, in := range fc.inputs {
		switch vv := in.(type) {
		case string:
			inputs += vv
		case *FilterChain:
			inputs += vv.Output()
		default:
		}
	}
	return inputs
}

func (fc *FilterChain) Output() string {
	if len(fc.filters) == 0 {
		return fc.Inputs()
	} else {
		return fc.output
	}
}

func (fc *FilterChain) Chain(f *Filter) *FilterChain {
	if f != nil {
		fc.filters = append(fc.filters, f)
	}
	return fc
}

func (fc *FilterChain) String() string {
	l := len(fc.filters)

	if l == 0 {
		// No filter in the chain, just return empty string as do nothing in the chain.
		return ""
	}

	str := fc.Inputs()
	for i, f := range fc.filters {
		str += f.String()
		if i < l-1 {
			str += ","
		}
	}

	str += fc.Output()
	return str
}

type Cmd struct {
	inputs          []string
	output          string
	fg              []*FilterChain
	selectedStreams []string
}

func New(output string) *Cmd {
	return &Cmd{inputs: []string{}, output: output, fg: []*FilterChain{}}
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

func (c *Cmd) Map(inputID int, streamType string, streamID int) {
	c.selectedStreams = append(c.selectedStreams, fmt.Sprintf("[%d:%s:%d]", inputID, streamType, streamID))
}

func (c *Cmd) MapFilterChainOutput(fc *FilterChain) error {
	re := regexp.MustCompile(`\[\w+\]`)
	streams := re.FindAllString(fc.Output(), -1)
	if len(streams) == 0 {
		fmt.Errorf("no ouput found in last filterchain")
	}

	for _, stream := range streams {
		c.selectedStreams = append(c.selectedStreams, stream)
	}
	return nil
}

func (c *Cmd) String() (string, error) {
	for i, fc := range c.fg {
		fmt.Printf("%d: fc: %v\n", i, fc)
	}

	str := `ffmpeg \
`
	for _, in := range c.inputs {
		str += fmt.Sprintf("-i \"%s\"\n", in)
	}

	str += "-filter_complex \n"

	l := len(c.fg)
	for i, fc := range c.fg {
		str += fc.String()
		if i < l-1 {
			str += ";\n"
		}
	}

	str += "\"\n"

	for _, stream := range c.selectedStreams {
		str += fmt.Sprintf("-map \"%s\"\n", stream)
	}

	str += c.output

	return str, nil
}
