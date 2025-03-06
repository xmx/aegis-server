package prereadtls

import (
	"bufio"
	"net"
	"time"
)

func Serve(lis net.Listener, tcpSrv, tlsSrv func(conn net.Conn)) error {
	var srv server
	if tlsSrv == nil {
		srv = &onlyTCPServer{tcpSrv: tcpSrv}
	} else {
		srv = &fullServer{tcpSrv: tcpSrv, tlsSrv: tlsSrv}
	}

	var tempDelay time.Duration
	for {
		rw, err := lis.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if maxDelay := time.Second; tempDelay > maxDelay {
					tempDelay = maxDelay
				}
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		go srv.serve(rw)
	}
}

type server interface {
	serve(net.Conn)
}

type onlyTCPServer struct {
	tcpSrv func(conn net.Conn)
}

func (ots *onlyTCPServer) serve(c net.Conn) {
	ots.tcpSrv(c)
}

type fullServer struct {
	tcpSrv func(net.Conn)
	tlsSrv func(net.Conn)
}

func (fsv *fullServer) serve(c net.Conn) {
	br := bufio.NewReader(c)
	conn := &peekedConn{Conn: c, br: br}
	_ = c.SetReadDeadline(time.Now().Add(10 * time.Second))
	if peek, err := br.Peek(1); err == nil && peek[0] == 0x16 { // TLS ClientHello
		fsv.tlsSrv(conn)
	} else {
		fsv.tcpSrv(conn)
	}
}
