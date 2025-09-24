package config

import (
	"context"
	"encoding/json/v2"
	"os"

	"github.com/xmx/aegis-common/jsmod"
	"github.com/xmx/aegis-control/library/jsonc"
	"github.com/xmx/jsos/jsvm"
)

type decoder interface {
	decode(ctx context.Context, data []byte, cfg *Config, name string) (*Config, error)
}

type jsondecoder struct{}

func (jsondecoder) decode(_ context.Context, raw []byte, cfg *Config, _ string) (*Config, error) {
	err := json.Unmarshal(raw, cfg)
	return cfg, err
}

type jsoncdecoder struct{}

func (jsoncdecoder) decode(_ context.Context, raw []byte, cfg *Config, _ string) (*Config, error) {
	data := jsonc.Translate(raw)
	err := json.Unmarshal(data, cfg)
	return cfg, err
}

type jsdecoder struct {
	eng jsvm.Engineer
}

func (j *jsdecoder) decode(ctx context.Context, data []byte, cfg *Config, name string) (*Config, error) {
	vari := jsmod.NewVariable[*Config]("aegis/server/config")
	vari.Set(cfg)

	eng := j.engine(ctx)
	eng.Require().Register(vari)
	_, err := eng.RunScript(name, string(data))

	return vari.Get(), err
}

func (j *jsdecoder) engine(ctx context.Context) jsvm.Engineer {
	if j.eng != nil {
		return j.eng
	}

	eng := jsvm.New(ctx)
	require := eng.Require()
	require.Register(jsmod.Modules()...)
	stdout, _ := eng.Output()
	stdout.Attach(os.Stdout)
	j.eng = eng

	return eng
}
