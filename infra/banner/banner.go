package banner

import (
	"fmt"
	"io"
	"runtime"
)

// ANSI 打印 banner 到指定输出流。
//
// 何为 ANSI 转义序列：https://en.wikipedia.org/wiki/ANSI_escape_code.
func ANSI(w io.Writer) {
	onceParseArgs()

	_, _ = fmt.Fprintf(w, ansiLogo, version, runtime.GOOS, runtime.GOARCH,
		hostname, username, workdir, pid, compileAt, commitAt, path, revision)
}
