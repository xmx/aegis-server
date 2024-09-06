package eventsource

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
)

// Accept 将请求转为 sse，如果不符合协议就返回 nil。
func Accept(w http.ResponseWriter, r *http.Request) EventSource {
	f, ok := w.(http.Flusher)
	if !ok || r.Header.Get("Accept") != "text/event-stream" ||
		r.Header.Get("Cache-Control") != "no-cache" {
		return nil
	}

	var gzw *gzip.Writer
	encodings := r.Header.Get("Accept-Encoding")
	for _, str := range strings.Split(encodings, ",") {
		if strings.TrimSpace(str) == "gzip" {
			w.Header().Set("Content-Encoding", "gzip")
			gzw = gzip.NewWriter(w)
		}
	}
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	f.Flush()

	parent := r.Context()
	ctx, cancel := context.WithCancel(parent)

	sse := &sseWriter{
		ctx:    ctx,
		cancel: cancel,
	}
	if gzw == nil {
		sse.write = w
		sse.flush = f
	} else {
		sse.write = gzw
		sse.close = gzw
		sse.flush = &gzipFlusher{write: gzw, flush: f}
	}

	return sse
}

type EventSource interface {
	io.WriteCloser
	Done() <-chan struct{}
	JSON(event string, data any) error
}

type sseWriter struct {
	mutex  sync.Mutex
	write  io.Writer
	flush  http.Flusher
	close  io.Closer
	ctx    context.Context
	cancel context.CancelFunc
}

func (sw *sseWriter) Write(p []byte) (int, error) {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	n := len(p)
	if err := sw.writeData(p); err != nil {
		return 0, err
	}

	return n, nil
}

func (sw *sseWriter) Close() error {
	var err error
	if c := sw.close; c != nil {
		err = c.Close()
	}
	sw.cancel()

	return err
}

func (sw *sseWriter) JSON(event string, data any) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	if event != "" {
		if err = sw.writeEvent([]byte(event)); err != nil {
			return err
		}
	}
	err = sw.writeData(msg)

	return err
}

func (sw *sseWriter) Done() <-chan struct{} {
	return sw.ctx.Done()
}

func (sw *sseWriter) writeEvent(evt []byte) error {
	_, err := sw.write.Write([]byte("event: "))
	if err == nil {
		if _, err = sw.write.Write(evt); err == nil {
			_, err = sw.write.Write([]byte{'\n'})
		}
	}
	if err != nil {
		_ = sw.Close()
	}

	return err
}

func (sw *sseWriter) writeData(data []byte) error {
	_, err := sw.write.Write([]byte("data: "))
	if err == nil {
		if _, err = sw.write.Write(data); err == nil {
			_, err = sw.write.Write([]byte{'\n', '\n'})
		}
	}
	if err != nil {
		_ = sw.Close()
		return err
	}
	sw.flush.Flush()

	return nil
}

type gzipFlusher struct {
	write *gzip.Writer
	flush http.Flusher
}

func (gf *gzipFlusher) Flush() {
	_ = gf.write.Flush()
	gf.flush.Flush()
}
