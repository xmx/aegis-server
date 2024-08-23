package vterminal

import (
	"io"

	"golang.org/x/crypto/ssh"
)

func NewSSH(network, addr string, cfg *ssh.ClientConfig, cols, rows int) (Typewriter, error) {
	var err error
	closers := make([]io.Closer, 0, 4)
	defer func() {
		if err != nil {
			for _, closer := range closers {
				_ = closer.Close()
			}
		}
	}()

	cli, err := ssh.Dial(network, addr, cfg)
	if err != nil {
		return nil, err
	}
	closers = append(closers, cli)
	sess, err := cli.NewSession()
	if err != nil {
		return nil, err
	}
	closers = append(closers, sess)

	stdin, err := sess.StdinPipe()
	if err != nil {
		return nil, err
	}
	closers = append(closers, stdin)
	stdout, err := sess.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if cols <= 0 {
		cols = 80
	}
	if rows <= 0 {
		rows = 40
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err = sess.RequestPty("xterm-256color", rows, cols, modes); err != nil {
		return nil, err
	}
	if err = sess.Shell(); err != nil {
		return nil, err
	}

	st := &sshTTY{
		cli:    cli,
		sess:   sess,
		stdin:  stdin,
		stdout: stdout,
		cols:   cols,
		rows:   rows,
	}

	return st, nil
}

type sshTTY struct {
	cli    *ssh.Client
	sess   *ssh.Session
	stdin  io.WriteCloser
	stdout io.Reader
	cols   int
	rows   int
}

func (st *sshTTY) Read(p []byte) (int, error) {
	return st.stdout.Read(p)
}

func (st *sshTTY) Write(p []byte) (int, error) {
	return st.stdin.Write(p)
}

func (st *sshTTY) Close() error {
	err := st.stdin.Close()
	if exx := st.sess.Close(); exx != nil && err == nil {
		err = exx
	}
	if exx := st.cli.Close(); exx != nil && err == nil {
		err = exx
	}

	return err
}

func (st *sshTTY) Size() (int, int, error) {
	return st.cols, st.rows, nil
}

func (st *sshTTY) Resize(cols, rows int) error {
	if cols <= 0 || rows <= 0 {
		return nil
	}

	if err := st.sess.WindowChange(rows, cols); err != nil {
		return err
	}
	st.cols, st.rows = cols, rows

	return nil
}
