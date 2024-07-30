package babel

import (
	_ "embed"
	"sync"

	"github.com/dop251/goja"
)

var (
	//go:embed babel.min.js
	gzipBabelJS      string
	onceCompileBabel = sync.OnceValues(compileBabel)
	pool             = &transPool{pool: sync.Pool{New: newBabelFunc}}
)

func ES5(script string) (string, error) {
	opt := map[string]any{
		"plugins": []string{"transform-modules-commonjs"},
	}
	return Transform(script, opt)
}

func Transform(script string, opt map[string]any) (string, error) {
	if opt == nil {
		opt = map[string]any{}
	}

	tran := pool.Get()
	defer pool.Put(tran)

	return tran.Transform(script, opt)
}

func compileBabel() (*goja.Program, error) {
	//reader := bytes.NewReader(gzipBabelJS)
	//zip, err := gzip.NewReader(reader)
	//if err != nil {
	//	return nil, err
	//}
	////goland:noinspection GoUnhandledErrorResult
	//defer zip.Close()
	//
	//code, err := io.ReadAll(zip)
	//if err != nil {
	//	return nil, err
	//}

	return goja.Compile("babel.min.js", gzipBabelJS, false)
}

type transformer struct {
	vm   *goja.Runtime
	err  error
	this goja.Value
	call goja.Callable
}

func (t *transformer) Transform(script string, opt map[string]any) (string, error) {
	if err := t.err; err != nil {
		return "", err
	}

	val, err := t.call(t.this, t.vm.ToValue(script), t.vm.ToValue(opt))
	if err != nil {
		return "", err
	}
	code := val.ToObject(t.vm).Get("code").String()

	return code, nil
}

func newBabelFunc() any {
	prog, err := onceCompileBabel()
	if err != nil {
		return &transformer{err: err}
	}

	vm := goja.New()
	logFunc := func(goja.FunctionCall) goja.Value { return nil }
	_ = vm.Set("console", map[string]func(goja.FunctionCall) goja.Value{
		"log":   logFunc,
		"error": logFunc,
		"warn":  logFunc,
	})
	if _, err = vm.RunProgram(prog); err != nil {
		return &transformer{err: err}
	}

	this := vm.Get("Babel")
	var call goja.Callable
	obj := this.ToObject(vm).Get("transform")
	if err = vm.ExportTo(obj, &call); err != nil {
		return &transformer{err: err}
	}

	return &transformer{vm: vm, this: this, call: call}
}

type transPool struct {
	pool sync.Pool
}

func (t *transPool) Get() *transformer {
	return t.pool.Get().(*transformer)
}

func (t *transPool) Put(trans *transformer) {
	t.pool.Put(trans)
}
