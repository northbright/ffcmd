package ffcmd

import "fmt"

type arg struct {
	k string
	v string
}

type Filter struct {
	name string
	args []arg
}

func NewFilter(name string) (*Filter, error) {
	if name == "" {
		return nil, fmt.Errorf("empty filter name")
	}
	return &Filter{name: name, args: []arg{}}, nil
}

func (f *Filter) AddArg(k, v string) {
	f.args = append(f.args, arg{k: k, v: v})
}

func (f *Filter) Cmd() string {
	l := len(f.args)
	cmd := fmt.Sprintf("%s=", f.name)
	for i, arg := range f.args {
		cmd += fmt.Sprintf("%s=%s", arg.k, arg.v)
		if i < l-1 {
			cmd += ":"
		}
	}
	return cmd
}

type FilterPipe struct {
	inputs   []string
	name     string
	filters  []*Filter
	outputFn func() string
}

func NewFilterPipe(name string, outputFunc func() string) (*FilterPipe, error) {
	if name == "" {
		return nil, fmt.Errorf("empty filter pipe name")
	}
	return &FilterPipe{name: name, outputFn: outputFunc}, nil
}

func (fp *FilterPipe) AddInput(input string) {
	fp.inputs = append(fp.inputs, input)
}

func (fp *FilterPipe) AddInputByID(id int, stream string, streamID int) {
	fp.inputs = append(fp.inputs, fmt.Sprintf("[%d:%s:%d]", id, stream, streamID))
}

func (fp *FilterPipe) AddInputByOutput(fpOut *FilterPipe) {
	fp.inputs = append(fp.inputs, fpOut.Output())
}

func (fp *FilterPipe) Inputs() string {
	inputs := ""
	for _, in := range fp.inputs {
		inputs += in
	}
	return inputs
}

func (fp *FilterPipe) AddFilter(f *Filter) {
	fp.filters = append(fp.filters, f)
}

func (fp *FilterPipe) Output() string {
	return fp.outputFn()
}

func (fp *FilterPipe) Cmd() string {
	cmd := fp.Inputs()
	l := len(fp.filters)

	for i, f := range fp.filters {
		cmd += f.Cmd()
		if i < l-1 {
			cmd += ","
		}
	}

	cmd += fp.Output()
	return cmd
}

type Cmd struct {
	inputs []string
	pipes  map[string]*FilterPipe
	output string
}

func New() *Cmd {
	return &Cmd{inputs: []string{}, pipes: make(map[string]*FilterPipe), output: ""}
}

func (c *Cmd) AddInput(input string) int {
	id := len(c.inputs)
	c.inputs = append(c.inputs, input)
	return id
}
