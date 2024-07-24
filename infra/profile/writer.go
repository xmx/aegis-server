package profile

import (
	"io"
	"sync"
	"sync/atomic"
)

type LogWriteCloser interface {
	io.WriteCloser
	Append(ws io.WriteCloser) (putOK bool)
	Remove(ws io.WriteCloser) (delOK bool)
}

type logWriteCloser struct {
	value atomic.Value
	mutex sync.Mutex
	elems map[io.WriteCloser]struct{}
}

func (lwc *logWriteCloser) Write(p []byte) (int, error) {
	if wrt, ok := lwc.value.Load().(io.Writer); ok && wrt != nil {
		return wrt.Write(p)
	}

	return len(p), nil
}

func (lwc *logWriteCloser) Close() error {
	lwc.mutex.Lock()
	closers := make([]io.Closer, 0, len(lwc.elems))
	for wrt := range lwc.elems {
		closers = append(closers, wrt)
	}
	clear(lwc.elems)
	lwc.value.Store(&atomicWriter{w: io.Discard})
	lwc.mutex.Unlock()

	var err error
	for _, c := range closers {
		if exx := c.Close(); exx != nil && err == nil {
			err = exx
		}
	}

	return err
}

func (lwc *logWriteCloser) Append(wc io.WriteCloser) bool {
	if wc == nil {
		return false
	}

	lwc.mutex.Lock()
	defer lwc.mutex.Unlock()

	if _, exist := lwc.elems[wc]; exist {
		return false
	}
	lwc.elems[wc] = struct{}{}

	size := len(lwc.elems)
	if size == 1 {
		lwc.value.Store(lwc.warp(wc))
	} else {
		wcs := make([]io.Writer, 0, size)
		for w := range lwc.elems {
			wcs = append(wcs, w)
		}
		wrt := io.MultiWriter(wcs...)
		lwc.value.Store(lwc.warp(wrt))
	}

	return true
}

func (lwc *logWriteCloser) Remove(wc io.WriteCloser) bool {
	if wc == nil {
		return false
	}

	lwc.mutex.Lock()
	defer lwc.mutex.Unlock()

	if _, exist := lwc.elems[wc]; !exist {
		return false
	}
	delete(lwc.elems, wc)

	size := len(lwc.elems)
	if size == 0 {
		lwc.value.Store(&atomicWriter{w: io.Discard})
		return true
	}

	wcs := make([]io.Writer, 0, size)
	for w := range lwc.elems {
		wcs = append(wcs, w)
	}
	if size == 1 {
		lwc.value.Store(lwc.warp(wcs[0]))
	} else {
		wrt := io.MultiWriter(wcs...)
		lwc.value.Store(lwc.warp(wrt))
	}

	return true
}

func (lwc *logWriteCloser) warp(w io.Writer) io.Writer {
	return &atomicWriter{w: w}
}

type nopWriteCloser struct {
	w io.Writer
}

func (ncw *nopWriteCloser) Write(p []byte) (int, error) {
	return ncw.w.Write(p)
}

func (ncw *nopWriteCloser) Close() error {
	return nil
}

type atomicWriter struct {
	w io.Writer
}

func (a *atomicWriter) Write(p []byte) (int, error) {
	return a.w.Write(p)
}
