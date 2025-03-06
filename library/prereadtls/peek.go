package prereadtls

import (
	"bufio"
	"net"
	"sync/atomic"
)

type peekedConn struct {
	net.Conn
	br *bufio.Reader
}

func (pc *peekedConn) Read(p []byte) (int, error) {
	return pc.br.Read(p)
}

type onceAcceptListener struct {
	conn net.Conn
	once atomic.Bool
}

func (oal *onceAcceptListener) Accept() (net.Conn, error) {
	if oal.once.CompareAndSwap(false, true) {
		return oal.conn, nil
	}

	return nil, net.ErrClosed
}

func (oal *onceAcceptListener) Close() error {
	// 注意：不要执行 oal.conn.Close()
	return nil
}

func (oal *onceAcceptListener) Addr() net.Addr {
	return oal.conn.RemoteAddr()
}

func NewOnceAccept(conn net.Conn) net.Listener {
	return &onceAcceptListener{conn: conn}
}
