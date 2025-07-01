package ffcmd

var (
	// The dir contains FFmpeg binaries(e.g. ffmpeg, ffprobe).
	// Default value is empty: it'll search FFmpeg binaries in $PATH when run the command.
	FFmpegBinDir = ""
)
