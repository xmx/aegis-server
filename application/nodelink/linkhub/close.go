package linkhub

import (
	"sync/atomic"

	"github.com/xmx/muxconn"
)

type safeClose struct {
	mux muxconn.Muxer
	flg atomic.Bool
}

func newSafeClose(mux muxconn.Muxer) *safeClose {
	return &safeClose{mux: mux}
}

func (cl *safeClose) close() {
	if cl.flg.CompareAndSwap(false, true) {
		_ = cl.mux.Close()
	}
}

func (cl *safeClose) closed() bool {
	return cl.flg.Load()
}
