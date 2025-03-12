package asciicast

import (
	"io"
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
