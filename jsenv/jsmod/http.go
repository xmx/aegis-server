package jsmod

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func NewHTTP() jsvm.GlobalRegister {
	return new(stdHTTP)
}

type stdHTTP struct {
	vm jsvm.Runtime
}

func (std *stdHTTP) RegisterGlobal(vm jsvm.Runtime) error {
	std.vm = vm
	hm := map[string]any{
		"post": std.post,
	}

	return vm.Runtime().Set("http", hm)
}

func (std *stdHTTP) post(addr string, body any) (string, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, addr, buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	res := new(bytes.Buffer)
	io.Copy(res, resp.Body)
	_ = resp.Body.Close()

	return res.String(), nil
}
