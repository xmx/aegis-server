package config

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

func Ext(fp string, limit ...int64) Loader {
	var lim int64
	if len(limit) > 0 {
		lim = limit[0]
	}

	dec := &activeFile{file: fp, lim: lim}
	ext := filepath.Ext(fp)
	switch ext {
	case ".js", ".javascript":
		dec.dec = &jsdecoder{}
	case ".jsonc":
		dec.dec = &jsoncdecoder{}
	default:
		dec.dec = &jsondecoder{}
	}

	return dec
}

type activeFile struct {
	file string
	lim  int64
	dec  decoder
}

func (af activeFile) Load(ctx context.Context) (*Config, error) {
	name := af.file
	cfg, err := af.decode(ctx, name, new(Config))
	if err != nil {
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

	return af.decode(ctx, join, cfg)
}

func (af activeFile) decode(ctx context.Context, name string, v *Config) (*Config, error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer fd.Close()

	rd := af.limiter(fd)
	data, err := io.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	return af.dec.decode(ctx, data, v, name)
}

func (af activeFile) limiter(r io.Reader) io.Reader {
	if af.lim <= 0 {
		return r
	}

	return io.LimitReader(r, af.lim)
}
