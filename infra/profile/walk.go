package profile

import (
	"os"

	"github.com/xmx/aegis-server/library/jsonc"
)

func JSONC(path string) (*Config, error) {
	cfg := new(Config)
	for name, err := range readdir(path, "*.jsonc") {
		if err != nil {
			return nil, err
		}
		if err = unmarshalJSONC(name, cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func unmarshalJSONC(name string, v any) error {
	raw, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	return jsonc.Unmarshal(raw, v)
}
