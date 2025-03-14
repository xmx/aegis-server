package jsmod

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/grafana/sobek"
)

func NewHTTP(clis ...*http.Client) *httpClient {
	cli := http.DefaultClient
	if len(clis) != 0 && clis[0] != nil {
		cli = clis[0]
	}

	return &httpClient{cli: cli}
}

type httpClient struct {
	vm  *sobek.Runtime
	cli *http.Client
}

func (hc *httpClient) SetGlobal(vm *sobek.Runtime) error {
	hc.vm = vm
	obj := vm.NewObject()
	if err := obj.Set("get", hc.get); err != nil {
		return err
	}

	return vm.Set("http", obj)
}

func (hc *httpClient) get(call sobek.FunctionCall) sobek.Value {
	args := call.Arguments
	size := len(args)
	if size == 0 {
		return sobek.Null()
	}

	strURL, _ := args[0].Export().(string)
	if strURL == "" {
		return sobek.Null()
	}
	resp, err := hc.cli.Get(strURL)
	if err != nil {
		return hc.vm.NewGoError(err)
	}
	res := &response{res: resp}

	return res.object(hc.vm)
}

type response struct {
	res *http.Response
}

func (r *response) object(vm *sobek.Runtime) sobek.Value {
	obj := vm.NewObject()
	_ = obj.Set("statusCode", r.res.StatusCode)
	_ = obj.Set("text", r.text)
	_ = obj.Set("json", r.json)
	_ = obj.Set("close", r.close)

	return obj
}

func (r *response) text() (string, error) {
	defer r.res.Body.Close()
	bs, err := io.ReadAll(r.res.Body)
	return string(bs), err
}

func (r *response) json() (any, error) {
	defer r.res.Body.Close()
	hm := make(map[string]any, 16)
	if err := json.NewDecoder(r.res.Body).Decode(&hm); err != nil {
		return nil, err
	}

	return hm, nil
}

func (r *response) close() error {
	return r.res.Body.Close()
}
