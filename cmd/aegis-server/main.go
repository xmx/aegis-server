package main

import (
	"os"

	"github.com/xmx/aegis-common/buildinfo"
)

func main() {
	buildinfo.ANSI(os.Stdout)
}
