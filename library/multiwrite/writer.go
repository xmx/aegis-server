package multiwrite

import (
	"io"
	"sync"
	"sync/atomic"
)

type Writer interface {
	io.Writer
	Attach(w io.Writer) bool
	Detach(w io.Writer) bool
}

func New(ws []io.Writer) Writer {
	writers := make(map[io.Writer]struct{}, 8)
	for _, w := range ws {
		if w != nil {
			writers[w] = struct{}{}
		}
	}
	mw := &multiWriter{writers: writers}
	mw.refresh()

	return mw
}

type multiWriter struct {
	val     atomic.Pointer[symbolWriter]
	mutex   sync.Mutex
	writers map[io.Writer]struct{}
}

func (mw *multiWriter) Write(p []byte) (int, error) {
	w := mw.val.Load()
	return w.Write(p)
}

func (mw *multiWriter) Attach(w io.Writer) bool {
	if w == nil {
		return false
	}

	mw.mutex.Lock()
	defer mw.mutex.Unlock()

	if _, exist := mw.writers[w]; exist {
		return false
	}
	mw.writers[w] = struct{}{}
	mw.refresh()

	return true
}

func (mw *multiWriter) Detach(w io.Writer) bool {
	if w == nil {
		return false
	}

	mw.mutex.Lock()
	defer mw.mutex.Unlock()

	if _, exist := mw.writers[w]; !exist {
		return false
	}
	delete(mw.writers, w)
	mw.refresh()

	return true
}

// refresh not safe for concurrent.
func (mw *multiWriter) refresh() {
	wrts := make([]io.Writer, 0, 8)
	for w := range mw.writers {
		wrts = append(wrts, w)
	}

	switch len(wrts) {
	case 0:
		mw.val.Store(&symbolWriter{w: io.Discard})
	case 1:
		mw.val.Store(&symbolWriter{w: wrts[0]})
	default:
		mw.val.Store(&symbolWriter{w: io.MultiWriter(wrts...)})
	}
}

type symbolWriter struct {
	w io.Writer
}

func (sw *symbolWriter) Write(p []byte) (int, error) {
	return sw.w.Write(p)
}
