package asciicast

import (
	"encoding/json"
	"io"
	"strconv"
	"sync"
	"time"
)

func NewHeader(width, height int) *Header {
	env := make(Envs, 4)
	env.SetShell("xterm-256color")

	return &Header{
		Version:   2,
		Width:     width,
		Height:    height,
		Env:       env,
		Timestamp: time.Now().Unix(),
	}
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

type Envs map[string]string

func (e Envs) SetShell(v string) {
	e["SHELL"] = v
}

func (e Envs) SetTerm(v string) {
	e["TERM"] = v
}

func (h Header) NewWriter(w io.Writer) Writer {
	return &asciiWriter{
		startAt: time.Now(),
		encoder: json.NewEncoder(w),
		header:  h,
	}
}

type asciiWriter struct {
	mutex   sync.Mutex
	startAt time.Time
	encoder *json.Encoder
	header  Header
	written bool
}

func (w *asciiWriter) Write(p []byte) (int, error) {
	n := len(p)
	err := w.write("o", string(p))
	return n, err
}

func (w *asciiWriter) Input(data string) error {
	return w.write("i", data)
}

func (w *asciiWriter) Marker(data string) error {
	return w.write("m", data)
}

func (w *asciiWriter) Resize(cols, rows int) error {
	// formatted as "{COLS}x{ROWS}"
	return w.write("c", strconv.Itoa(cols)+"x"+strconv.Itoa(rows))
}

func (w *asciiWriter) write(code, data string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if err := w.onceWriteHeader(); err != nil {
		return err
	}

	nano := time.Since(w.startAt).Nanoseconds()
	offset := float64(nano) / float64(time.Second)
	lines := []any{offset, code, data}

	return w.encoder.Encode(lines)
}

func (w *asciiWriter) onceWriteHeader() error {
	if !w.written {
		w.written = true
		return w.encoder.Encode(w.header)
	}

	return nil
}
