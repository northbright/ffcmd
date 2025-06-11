package ffcmd

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Sub struct {
	Text  string
	Start string
	End   string
}

type SRT struct {
	filename string
	subs     []Sub
}

func NewSRT(videoFile string, subs []Sub) (*SRT, error) {
	if videoFile == "" {
		return nil, fmt.Errorf("empty video filename")
	}

	if len(subs) == 0 {
		return nil, fmt.Errorf("empty subtitles")
	}

	filename := strings.Replace(videoFile, filepath.Ext(videoFile), ".srt", -1)

	return &SRT{filename: filename, subs: subs}, nil
}

func NewSRTForTrimedVideo(videoFile string, trimStart, trimEnd, subtitle string) (*SRT, error) {
	return nil, nil
}

func (srt *SRT) Filename() string {
	return srt.filename
}

func (srt *SRT) CreateCmd() string {
	cmd := `echo -ne "`
	for i, sub := range srt.subs {
		cmd += fmt.Sprintf(`%d\n%s --> %s\n%s`, i+1, sub.Start, sub.End, sub.Text)
	}
	cmd += fmt.Sprintf(`" > "%s"`, srt.filename)
	return cmd
}

func (srt *SRT) RemoveCmd() string {
	cmd := fmt.Sprintf("rm \"%s\"", srt.filename)
	return cmd
}
