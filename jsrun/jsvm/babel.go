package jsvm

import (
	_ "embed"
	"sync"

	"github.com/dop251/goja"
)

var (
	//go:embed babel.js
	babelJS string

	onceCompileBabel = sync.OnceValues(func() (*goja.Program, error) {
		return goja.Compile("babel.js", babelJS, false)
	})
)
