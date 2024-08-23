package ioext

import (
	"io"
	"sync"
	"sync/atomic"
)

type AttachWriter interface {
	io.Writer
	Attach(w io.Writer) bool
	Leave(w io.Writer) bool
}

func NewAttachWriter() AttachWriter {
	return &attachWriter{
		val:     atomic.Pointer[atomicWriter]{},
		writers: make(map[io.Writer]struct{}, 8),
	}
}

type attachWriter struct {
	val     atomic.Pointer[atomicWriter]
	mutex   sync.Mutex
	writers map[io.Writer]struct{}
}

func (aw *attachWriter) Write(p []byte) (int, error) {
	return aw.loadWriter().Write(p)
}

func (aw *attachWriter) Attach(w io.Writer) bool {
	switch w.(type) {
	case nil, *attachWriter:
		return false
	}

	aw.mutex.Lock()
	defer aw.mutex.Unlock()

	if _, ok := aw.writers[w]; ok {
		return false
	}

	aw.writers[w] = struct{}{}
	aw.refreshWriter()

	return true
}

func (aw *attachWriter) Leave(w io.Writer) bool {
	if w == nil {
		return false
	}

	aw.mutex.Lock()
	defer aw.mutex.Unlock()

	if _, ok := aw.writers[w]; !ok {
		return false
	}
	delete(aw.writers, w)
	aw.refreshWriter()

	return true
}

func (aw *attachWriter) loadWriter() *atomicWriter {
	if aw := aw.val.Load(); aw != nil {
		return aw
	}
	return aw.discordWriter()
}

func (aw *attachWriter) refreshWriter() {
	size := len(aw.writers)
	writers := make([]io.Writer, 0, size)
	for w := range aw.writers {
		writers = append(writers, w)
	}

	if size == 0 {
		aw.val.Store(aw.discordWriter())
	} else if size == 1 {
		aw.val.Store(&atomicWriter{w: writers[0]})
	} else {
		aw.val.Store(&atomicWriter{w: io.MultiWriter(writers...)})
	}
}

func (*attachWriter) discordWriter() *atomicWriter {
	return &atomicWriter{w: io.Discard}
}

type atomicWriter struct {
	w io.Writer
}

func (a *atomicWriter) Write(p []byte) (int, error) {
	return a.w.Write(p)
}
