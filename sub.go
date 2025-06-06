package ffcmd

import "fmt"

type Sub struct {
	text  string
	start string
	end   string
}

type SRT struct {
	filename string
	subs     []Sub
}

func NewSRT(filename string, subs []Sub) (*SRT, error) {
	if filename == "" {
		return nil, fmt.Errorf("empty SRT filename")
	}

	if len(subs) == 0 {
		return nil, fmt.Errorf("empty subtitles")
	}

	return &SRT{subs: subs}, nil
}

func (srt *SRT) CreateCmd() string {
	cmd := `echo -ne "`
	for i, sub := range srt.subs {
		cmd += fmt.Sprintf("%d\n%s --> %s\n%s\n", i+1, sub.start, sub.end, sub.text)
	}
	cmd += fmt.Sprintf("\" > %s.srt", srt.filename)
	return cmd
}
