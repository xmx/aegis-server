package vterminal

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
)

func NewPTMX(cmd *exec.Cmd, cols, rows int) (Typewriter, error) {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}
	if cols <= 0 {
		cols = 80
	}
	if rows <= 0 {
		rows = 40
	}
	size := &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)}
	if err = pty.Setsize(ptmx, size); err != nil {
		_ = ptmx.Close()
		return nil, err
	}

	pt := &ptmxTTY{
		ptmx: ptmx,
		cmd:  cmd,
		cols: cols,
		rows: rows,
	}

	return pt, nil
}

type ptmxTTY struct {
	ptmx *os.File
	cmd  *exec.Cmd
	cols int
	rows int
}

func (pt *ptmxTTY) Read(p []byte) (int, error) {
	return pt.ptmx.Read(p)
}

func (pt *ptmxTTY) Write(p []byte) (int, error) {
	return pt.ptmx.Write(p)
}

func (pt *ptmxTTY) Close() error {
	if p := pt.cmd.Process; p != nil {
		_ = p.Kill()
	}
	return pt.ptmx.Close()
}

func (pt *ptmxTTY) Size() (cols, rows int, err error) {
	return pt.cols, pt.rows, nil
}

func (pt *ptmxTTY) Resize(cols, rows int) error {
	size := &pty.Winsize{Rows: uint16(cols), Cols: uint16(rows)}
	if err := pty.Setsize(pt.ptmx, size); err != nil {
		return err
	}
	pt.cols, pt.rows = cols, rows

	return nil
}
