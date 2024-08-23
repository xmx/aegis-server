package asciicast

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
)

type Writer interface {
	// Writer https://docs.asciinema.org/manual/asciicast/v2/#o-output-data-written-to-a-terminal
	io.Writer

	// Input https://docs.asciinema.org/manual/asciicast/v2/#i-input-data-read-from-a-terminal
	Input(data string) error

	// Marker https://docs.asciinema.org/manual/asciicast/v2/#m-marker
	Marker(data string) error

	// Resize https://docs.asciinema.org/manual/asciicast/v2/#r-resize
	Resize(cols, rows int) error
}

// Header asciicast header is JSON-encoded object containing recording meta-data.
//
// https://docs.asciinema.org/manual/asciicast/v2/#header
type Header struct {
	// Version Must be set to 2. Integer.
	Version int `json:"version"`

	// Width Initial terminal width, i.e number of columns. Integer.
	Width int `json:"width,omitempty"`

	// Height Initial terminal height, i.e. number of rows. Integer.
	Height int `json:"height,omitempty"`

	// Timestamp Unix timestamp of the beginning of the recording session. Integer.
	Timestamp int64 `json:"timestamp,omitempty"`

	// Title of the asciicast, as given via -t option to asciinema rec. String.
	Title string `json:"title,omitempty"`

	// Duration of the whole recording in seconds (when it's known upfront). Float.
	Duration float64 `json:"duration,omitempty"`

	// IdleTimeLimit Idle time limit, as given via -i option to asciinema rec. Float.
	//
	// This should be used by an asciicast player to reduce all terminal
	// inactivity (delays between frames) to maximum of idle_time_limit value.
	IdleTimeLimit float64 `json:"idle_time_limit,omitempty"`

	// Command that was recorded, as given via -c option to asciinema rec. String.
	Command string `json:"command,omitempty"`

	// Env Map of captured environment variables. Object (String -> String).
	//
	// Official asciinema recorder captures only SHELL and TERM by default. All implementations
	// of asciicast-compatible terminal recorder should not capture any additional environment
	// variables unless explicitly requested by the user.
	Env map[string]string `json:"env,omitempty"`

	// Theme Color theme of the recorded terminal. Object, with the following attributes.
	Theme *Theme `json:"theme,omitempty"`
}

func (h *Header) EnvSetShell(v string) *Header {
	h.setEnv("SHELL", v)
	return h
}

func (h *Header) EnvSetTerm(v string) *Header {
	h.setEnv("TERM", v)
	return h
}

func (h *Header) setEnv(k, v string) {
	if k == "" || v == "" {
		return
	}
	if h.Env == nil {
		h.Env = make(map[string]string, 4)
	}
	h.Env[k] = v
}

// Theme Color theme of the recorded terminal. Object, with the following attributes.
//
// https://docs.asciinema.org/manual/asciicast/v2/#theme
type Theme struct {
	// FG normal text color.
	FG string `json:"fg,omitempty"`

	// BG normal background color.
	BG string `json:"bg,omitempty"`

	// list of 8 or 16 colors, separated by colon character.
	Palette string `json:"palette,omitempty"`
}

type CodeType string

const (
	CodeOutput CodeType = "o"
	CodeInput  CodeType = "i"
	CodeMarker CodeType = "m"
	CodeResize CodeType = "r"
)

func NewWriter(wrt io.Writer, h *Header) Writer {
	now := time.Now()
	if h != nil {
		h.Version = 2
		h.Timestamp = now.Unix()
		if h.Height <= 0 {
			h.Height = 24
		}
		if h.Width <= 0 {
			h.Width = 80
		}
	}

	return &castWriter{
		encoder: json.NewEncoder(wrt),
		wrt:     wrt,
		startAt: now,
		header:  h,
	}
}

type castWriter struct {
	mutex   sync.Mutex
	encoder *json.Encoder
	wrt     io.Writer
	startAt time.Time
	written bool
	header  *Header
}

func (cw *castWriter) Write(p []byte) (int, error) {
	n := len(p)
	if err := cw.write(CodeOutput, string(p)); err != nil {
		return 0, err
	}

	return n, nil
}

func (cw *castWriter) Resize(cols, rows int) error {
	if cols <= 0 || rows <= 0 {
		return nil
	}

	// formatted as "{COLS}x{ROWS}"
	data := fmt.Sprintf("%dx%d", cols, rows)

	return cw.write(CodeResize, data)
}

func (cw *castWriter) Marker(data string) error {
	return cw.write(CodeMarker, data)
}

func (cw *castWriter) Input(data string) error {
	return cw.write(CodeInput, data)
}

func (cw *castWriter) write(code CodeType, data string) error {
	cw.mutex.Lock()
	defer cw.mutex.Unlock()

	if cw.header != nil && !cw.written {
		cw.written = true
		if err := cw.writeHeader(); err != nil {
			return err
		}
	}

	nano := time.Since(cw.startAt).Nanoseconds()
	offset := float64(nano) / float64(time.Second)
	lines := []any{offset, code, data}

	return cw.encoder.Encode(lines)
}

func (cw *castWriter) writeHeader() error {
	return cw.encoder.Encode(cw.header)
}
