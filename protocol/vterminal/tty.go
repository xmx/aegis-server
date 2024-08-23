package vterminal

import "io"

type Typewriter interface {
	io.ReadWriteCloser
	Size() (cols, rows int, err error)
	Resize(cols, rows int) error
}
