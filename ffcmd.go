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
	inputs   []string
	filters  []*Filter
	outputFn func() string
}

func NewFilterChain(outputFunc func() string) *FilterChain {
	return &FilterChain{inputs: []string{}, outputFn: outputFunc}
}

func (fc *FilterChain) AddInput(input string) {
	fc.inputs = append(fc.inputs, input)
}

func (fc *FilterChain) AddStreamInput(inputID int, streamType string, streamID int) {
	fc.inputs = append(fc.inputs, fmt.Sprintf("[%d:%s:%d]", inputID, streamType, streamID))
}

func (fc *FilterChain) AddOutputAsInput(fpOut *FilterChain) {
	fc.inputs = append(fc.inputs, fpOut.Output())
}

func (fc *FilterChain) Inputs() string {
	inputs := ""
	for _, in := range fc.inputs {
		inputs += in
	}
	return inputs
}

func (fc *FilterChain) Chain(f *Filter) *FilterChain {
	fc.filters = append(fc.filters, f)
	return fc
}

func (fc *FilterChain) Output() string {
	if len(fc.filters) == 0 {
		// No filters in the chain, return inputs directly.
		return fc.Inputs()
	} else {
		return fc.outputFn()
	}
}

func (fc *FilterChain) String() string {
	str := fc.Inputs()
	l := len(fc.filters)

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
	inputs []string
	output string
	fg     []*FilterChain
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

func (c *Cmd) String() (string, error) {
	str := `ffmpeg \
`
	for _, in := range c.inputs {
		str += fmt.Sprintf("-i \"%s\"\n", in)
	}

	str += "-filter_complex \n"

	l := len(c.fg)
	var selectedStreams []string

	for i, fc := range c.fg {
		str += fc.String()
		if i < l-1 {
			str += ";\n"
		} else {
			// Use outputs of final filterchain as selected streams
			re := regexp.MustCompile(`\[\w+\]`)
			selectedStreams = re.FindAllString(fc.Output(), -1)
			if len(selectedStreams) == 0 {
				return "", fmt.Errorf("no ouput found in last filterchain")
			}
		}
	}

	str += "\"\n"

	for _, stream := range selectedStreams {
		str += fmt.Sprintf("-map \"%s\"\n", stream)
	}

	str += c.output

	return str, nil
}
