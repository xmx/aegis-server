package config

import (
	"context"
	"encoding/json/v2"
	"io"
	"os"
	"path/filepath"

	"github.com/xmx/aegis-control/library/jsonc"
)

type Loader interface {
	Load(context.Context) (*Config, error)
}

func Ext(fp string, limit ...int64) Loader {
	var lim int64
	if len(limit) > 0 {
		lim = limit[0]
	}

	dec := &activeFile{file: fp, lim: lim}
	ext := filepath.Ext(fp)
	switch ext {
	case ".jsonc":
		dec.ums = jsoncmarshal{}
	default:
		dec.ums = jsonmarshal{}
	}

	return dec
}

type unmarshaler interface {
	Unmarshal([]byte, any) error
}

type jsonmarshal struct{}

func (jsonmarshal) Unmarshal(raw []byte, v any) error {
	return json.Unmarshal(raw, v)
}

type jsoncmarshal struct{}

func (jsoncmarshal) Unmarshal(raw []byte, v any) error {
	data := jsonc.Translate(raw)
	return json.Unmarshal(data, v)
}

type activeFile struct {
	file string
	lim  int64
	ums  unmarshaler
}

func (af activeFile) Load(context.Context) (*Config, error) {
	cfg := new(Config)
	name := af.file
	if err := af.unmarshal(name, cfg); err != nil {
		return nil, err
	}

	act := cfg.Active
	if act == "" {
		return cfg, nil
	}

	dir := filepath.Dir(name)
	base := filepath.Base(name)
	ext := filepath.Ext(base)
	length, size := len(base), len(ext)
	afile := base[:length-size] + "-" + act + ext
	join := filepath.Join(dir, afile)
	if err := af.unmarshal(join, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (af activeFile) unmarshal(name string, v any) error {
	fd, err := os.Open(name)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer fd.Close()

	rd := af.limiter(fd)
	data, err := io.ReadAll(rd)
	if err != nil {
		return err
	}

	return af.ums.Unmarshal(data, v)
}

func (af activeFile) limiter(r io.Reader) io.Reader {
	if af.lim <= 0 {
		return r
	}

	return io.LimitReader(r, af.lim)
}
