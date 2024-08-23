package profile

import (
	"encoding/json"
	"os"
)

func JSON(path string) (*Config, error) {
	cfg := new(Config)
	for name, err := range readdir(path, "*.json") {
		if err != nil {
			return nil, err
		}
		if err = unmarshalJSON(name, cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func unmarshalJSON(name string, v any) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()
	dec := json.NewDecoder(file)

	return dec.Decode(v)
}
