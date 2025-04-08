package qmux

import (
	"context"
	"io"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/quic-go/quic-go"
)

type Muxer interface {
	io.Closer

	// Listener 获取 [net.Listener]。
	Listener() net.Listener

	// OpenStream 开启同步流，用于自定义扩展。
	OpenStream(ctx context.Context) (quic.Stream, error)

	// DialContext 拨号器。
	DialContext(ctx context.Context, network, address string) (net.Conn, error)

	// InternalURL 内部通信的 URL 地址。
	InternalURL(path string) *url.URL

	// RxBytes 接收到的数据字节数。
	RxBytes() uint64

	// TxBytes 发送出去的数据字节数。
	TxBytes() uint64
}

func New(conn quic.Connection) Muxer {
	dial := &net.Dialer{
		Timeout:         time.Minute,
		KeepAliveConfig: net.KeepAliveConfig{Enable: true},
	}
	dial.SetMultipathTCP(true)

	return &quicMux{
		conn:  conn,
		dial:  dial,
		addr:  "tunnel.internal:80",
		laddr: conn.LocalAddr(),
		raddr: conn.RemoteAddr(),
	}
}

type quicMux struct {
	conn    quic.Connection
	dial    *net.Dialer
	addr    string
	laddr   net.Addr
	raddr   net.Addr
	rxBytes atomic.Uint64
	txBytes atomic.Uint64
	lismu   sync.Mutex
	listen  *quicListener
}

func (qm *quicMux) Close() error {
	return qm.conn.CloseWithError(0, "")
}

func (qm *quicMux) Listener() net.Listener {
	qm.lismu.Lock()
	if qm.listen == nil {
		ctx, cancel := context.WithCancel(context.Background())
		qm.listen = &quicListener{
			mux:    qm,
			ctx:    ctx,
			cancel: cancel,
		}
	}
	listen := qm.listen
	qm.lismu.Unlock()

	return listen
}

func (qm *quicMux) OpenStream(ctx context.Context) (quic.Stream, error) {
	return qm.conn.OpenStreamSync(ctx)
}

func (qm *quicMux) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	if network == "tcp" && address == qm.addr { // HTTP 走内部通道
		stm, err := qm.OpenStream(ctx)
		if err != nil {
			return nil, err
		}
		conn := &quicConn{stm: stm, mux: qm}

		return conn, nil
	}

	return qm.dial.DialContext(ctx, network, address)
}

func (qm *quicMux) InternalURL(path string) *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   qm.addr,
		Path:   path,
	}
}

func (qm *quicMux) RxBytes() uint64 {
	return qm.rxBytes.Load()
}

func (qm *quicMux) TxBytes() uint64 {
	return qm.txBytes.Load()
}

type quicListener struct {
	mux    *quicMux
	ctx    context.Context
	cancel context.CancelFunc
}

func (ql *quicListener) Accept() (net.Conn, error) {
	select {
	case <-ql.ctx.Done():
		return nil, net.ErrClosed
	default:
	}
	stm, err := ql.mux.conn.AcceptStream(ql.ctx)
	if err != nil {
		return nil, err
	}
	conn := &quicConn{stm: stm, mux: ql.mux}

	return conn, nil
}

func (ql *quicListener) Close() error {
	select {
	case <-ql.ctx.Done():
		return net.ErrClosed
	default:
		ql.cancel()
		return nil
	}
}

func (ql *quicListener) Addr() net.Addr {
	return ql.mux.laddr
}

type quicConn struct {
	mux *quicMux
	stm quic.Stream
}

func (qc *quicConn) Read(b []byte) (int, error) {
	n, err := qc.stm.Read(b)
	if n > 0 {
		qc.mux.rxBytes.Add(uint64(n))
	}
	return n, err
}

func (qc *quicConn) Write(b []byte) (int, error) {
	n, err := qc.stm.Write(b)
	if n > 0 {
		qc.mux.txBytes.Add(uint64(n))
	}

	return n, err
}

func (qc *quicConn) Close() error {
	return qc.stm.Close()
}

func (qc *quicConn) LocalAddr() net.Addr {
	return qc.mux.laddr
}

func (qc *quicConn) RemoteAddr() net.Addr {
	return qc.mux.raddr
}

func (qc *quicConn) SetDeadline(t time.Time) error {
	return qc.stm.SetDeadline(t)
}

func (qc *quicConn) SetReadDeadline(t time.Time) error {
	return qc.stm.SetReadDeadline(t)
}

func (qc *quicConn) SetWriteDeadline(t time.Time) error {
	return qc.stm.SetWriteDeadline(t)
}
