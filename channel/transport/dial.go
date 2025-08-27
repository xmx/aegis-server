package transport

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"time"

	"github.com/xtaci/smux"
	"golang.org/x/net/quic"
)

type Dialer interface {
	DialContext(ctx context.Context, addr string, data, reply any) (Muxer, error)
}

func NewDialer(quicCfg *quic.Config) (Dialer, error) {
	udp, err := quic.Listen("udp", "", quicCfg)
	if err != nil {
		return nil, err
	}

	ndl := &net.Dialer{
		Timeout:         10 * time.Second,
		KeepAliveConfig: net.KeepAliveConfig{Enable: true},
	}
	ndl.SetMultipathTCP(true)
	tlsCfg := quicCfg.TLSConfig
	tcp := &tls.Dialer{NetDialer: ndl, Config: tlsCfg}
	dd := &dualDialer{
		cfg: quicCfg,
		tcp: tcp,
		udp: udp,
	}

	return dd, nil
}

// dualDialer 支持 udp 和 tcp。
type dualDialer struct {
	cfg *quic.Config
	tcp *tls.Dialer
	udp *quic.Endpoint
}

func (dd *dualDialer) DialContext(ctx context.Context, addr string, data, reply any) (Muxer, error) {
	if mux, _ := dd.dialUDP(ctx, addr, data, reply); mux != nil {
		return mux, nil
	}

	return dd.dialTCP(ctx, addr, data, reply)
}

func (dd *dualDialer) dialUDP(parent context.Context, addr string, data, reply any) (Muxer, error) {
	mst, err := dd.udp.Dial(parent, "udp", addr, dd.cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(parent, time.Minute)
	defer cancel()

	stm, err := mst.NewStream(ctx)
	if err != nil {
		_ = mst.Close()
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer stm.Close()

	stm.SetReadContext(ctx)
	stm.SetWriteContext(ctx)

	if err = json.NewEncoder(stm).Encode(data); err != nil {
		_ = mst.Close()
		return nil, err
	}
	if err = json.NewDecoder(stm).Decode(reply); err != nil {
		_ = mst.Close()
		return nil, err
	}

	laddr := mst.LocalAddr()
	raddr := mst.RemoteAddr()
	mux := &udpMux{
		qc:    mst,
		laddr: &udpAddr{addr: laddr},
		raddr: &udpAddr{addr: raddr},
	}

	return mux, nil
}

func (dd *dualDialer) dialTCP(ctx context.Context, addr string, data, reply any) (Muxer, error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(data); err != nil {
		return nil, err
	}

	conn, err := dd.tcp.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}

	reqURL := &url.URL{
		Scheme: "https",
		Host:   addr,
		Path:   "/api/channel",
	}
	strURL := reqURL.String()
	req, err := http.NewRequest(http.MethodPost, strURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	now := time.Now()
	deadline := now.Add(time.Minute)
	_ = conn.SetDeadline(deadline)
	if err = req.Write(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(reply); err != nil {
		_ = conn.Close()
		return nil, err
	}

	code := resp.StatusCode
	if code/100 != 2 {
		_ = conn.Close()
		return nil, fmt.Errorf(resp.Status)
	}

	cfg := smux.DefaultConfig()
	sess, err := smux.Client(conn, cfg)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	mux := &tcpMux{sess: sess}

	return mux, nil
}

type udpAddr struct {
	addr netip.AddrPort
}

func (u udpAddr) Network() string {
	return "udp"
}

func (u udpAddr) String() string {
	return u.addr.String()
}
