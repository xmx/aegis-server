package wsocket

import "io"

func JSWriter(conn *Conn, kind Kind) io.Writer {
	if kind == "" {
		kind = KindStdout
	}

	return &jsWriter{
		conn: conn,
		kind: kind,
	}
}

type jsWriter struct {
	conn *Conn
	kind Kind
}

func (jw *jsWriter) Write(p []byte) (int, error) {
	n := len(p)
	body := &Body{Kind: jw.kind, Data: string(p)}
	if err := jw.conn.WriteJSON(body); err != nil {
		return 0, err
	}

	return n, nil
}
