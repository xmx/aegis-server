package shipx

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/xgfone/ship/v5"
)

func SSE(c *ship.Context) (*EventSource, error) {
	// 粗略的检查一下是不是 EventSource 请求
	w, r := c.Response().ResponseWriter, c.Request()
	f, ok := w.(http.Flusher)
	if !ok || r.Header.Get(ship.HeaderAccept) != "text/event-stream" ||
		r.Header.Get(ship.HeaderCacheControl) != "no-cache" {
		return nil, ship.ErrUnsupportedMediaType
	}

	w.Header().Set(ship.HeaderContentType, "text/event-stream; charset=utf-8")
	w.Header().Set(ship.HeaderCacheControl, "no-cache")
	w.Header().Set(ship.HeaderConnection, "keep-alive")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte{'\n'})
	f.Flush()

	parent := r.Context()
	ctx, cancel := context.WithCancel(parent)

	sse := &EventSource{
		flush:  f,
		write:  w,
		ctx:    ctx,
		cancel: cancel,
	}

	return sse, nil
}

type EventSource struct {
	mutex  sync.Mutex
	flush  http.Flusher
	write  http.ResponseWriter
	ctx    context.Context
	cancel context.CancelFunc
}

func (e *EventSource) Write(p []byte) (int, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	_, err := e.write.Write([]byte("data: "))
	if err == nil {
		if _, err = e.write.Write(p); err == nil {
			_, err = e.write.Write([]byte("\n\n"))
			e.flush.Flush()
		}
	}
	if err != nil {
		_ = e.Close()
		return 0, err
	}

	return len(p), nil
}

func (e *EventSource) Close() error {
	e.cancel()
	return nil
}

func (e *EventSource) JSON(event string, data any) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	var err error
	if event != "" {
		event = "event: " + event + "\n"
		_, err = e.write.Write([]byte(event))
	}
	if err != nil {
		_ = e.Close()
		return err
	}

	if _, err = e.write.Write([]byte("data: ")); err == nil {
		if err = json.NewEncoder(e.write).Encode(data); err == nil {
			_, err = e.write.Write([]byte("\n\n"))
		}
	}
	e.flush.Flush()
	if err != nil {
		_ = e.Close()
		return err
	}

	return nil
}

func (e *EventSource) Done() <-chan struct{} {
	return e.ctx.Done()
}
