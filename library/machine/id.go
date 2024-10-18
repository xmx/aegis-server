package machine

import (
	"io"
	"os"
	"strings"
)

func ID() (string, error) {
	mid, err := machineID()
	if err != nil {
		return "", err
	}
	mid = normalization(mid)

	return mid, nil
}

func readFile(file string) (string, error) {
	src, err := os.Open(file)
	if err != nil {
		return "", err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer src.Close()

	data, err := io.ReadAll(io.LimitReader(src, 1024))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func normalization(mid string) string {
	lower := strings.ToLower(mid)
	runes := make([]rune, 0, len(lower))
	for _, ch := range lower {
		if ('0' <= ch && ch <= '9') || ('a' <= ch && ch <= 'z') {
			runes = append(runes, ch)
		}
	}

	return string(runes)
}
