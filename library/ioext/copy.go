package ioext

import "io"

func Copy(dst io.WriteCloser, src io.Reader) (int64, error) {
	//goland:noinspection GoUnhandledErrorResult
	defer dst.Close()
	return io.Copy(dst, src)
}
