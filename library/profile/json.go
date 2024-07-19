package profile

import (
	"encoding/json"
	"os"
)

func JSON(path string, v any) error {
	for name, err := range Readdir(path, "*.json") {
		if err != nil {
			return err
		}
		if err = unmarshalJSON(name, v); err != nil {
			return err
		}
	}

	return nil
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
