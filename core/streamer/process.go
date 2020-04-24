package streamer

import (
	"fmt"
	"os"
	"os/exec"
)

// IProcess is an interface around the FFMPEG process
type IProcess interface {
	Spawn(path, URI string) *exec.Cmd
}

// ProcessLoggingOpts describes options for process logging
type ProcessLoggingOpts struct {
	Enabled    bool   // Option to set logging for transcoding processes
	Directory  string // Directory for the logs
	MaxSize    int    // Maximum size of kept logging files in megabytes
	MaxBackups int    // Maximum number of old log files to retain
	MaxAge     int    // Maximum number of days to retain an old log file.
	Compress   bool   // Indicates if the log rotation should compress the log files
}

// Process is the main type for creating new processes
type Process struct {
	keepFiles   bool
	audio       bool
	loggingOpts ProcessLoggingOpts
}

// Type check
var _ IProcess = (*Process)(nil)

// NewProcess creates a new process able to spawn transcoding FFMPEG processes
func NewProcess(
	keepFiles bool,
	audio bool,
	loggingOpts ProcessLoggingOpts,
) *Process {
	return &Process{keepFiles, audio, loggingOpts}
}

// getHLSFlags are for getting the flags based on the config context
func (p Process) getHLSFlags() string {
	if p.keepFiles {
		return "append_list"
	}
	return "delete_segments+append_list"
}

// Spawn creates a new FFMPEG cmd
func (p Process) Spawn(path, URI string) *exec.Cmd {
	/**
	-filter_complex "drawtext=fontsize=30:fontfile=FreeSerif.ttf:textfile=/tmp/text.txt:reload=1:x=(w-text_w)/2:y=(h-text_h)/2" \
	*/
	os.MkdirAll(path, os.ModePerm)
	processCommands := []string{
		"-y",
		"-fflags",
		"nobuffer",
		"-rtsp_transport",
		"tcp",
		"-i",
		URI,
		"-vsync",
		"0",
		"-c:v",
		"libx264",
		"-crf",
		"30",
		"-preset",
		"fast",
		"-pix_fmt",
		"yuv420p",
		"-flags",
		"+cgop",
		"-g",
		"50",
		"-movflags",
		"frag_keyframe+empty_moov",
	}
	if p.audio {
		processCommands = append(processCommands, "-an")
	}
	processCommands = append(processCommands,
		"-hls_flags",
		p.getHLSFlags(),
		"-f",
		"hls",
		"-segment_list_flags",
		"live",
		"-hls_time",
		"1",
		"-hls_list_size",
		"5",
		"-hls_delete_threshold",
		"5",
		"-hls_wrap",
		"5",
		"-hls_segment_filename",
		fmt.Sprintf("%s/%%d.ts", path),
		fmt.Sprintf("%s/index.m3u8", path),
	)
	cmd := exec.Command("ffmpeg", processCommands...)
	return cmd
}
