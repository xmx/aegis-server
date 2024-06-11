package profile

import (
	"encoding/json"
	"os"
)

func Load(filename string, v any) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	dec := json.NewDecoder(file)

	return dec.Decode(v)
}
