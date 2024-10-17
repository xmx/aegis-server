package i18n

import (
	"encoding/json"
	"os"
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

func TestExample(t *testing.T) {
	t.Logf("%s", language.SimplifiedChinese)

	builder, err := NewBuilder()
	if err != nil {
		t.Fatal(err)
	}

	detector := NewDetector(builder)
	str := detector.Detect([]language.Tag{language.Russian}, "hello", "Alice")
	t.Log(str)
}

func NewBuilder() (*catalog.Builder, error) {
	file, err := os.Open("i18n.json")
	if err != nil {
		return nil, err
	}
	msgs := make(map[string]map[string]string, 64)
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&msgs); err != nil {
		return nil, err
	}

	builder := catalog.NewBuilder(catalog.Fallback(language.SimplifiedChinese))
	for key, msg := range msgs {
		for lang, str := range msg {
			tag, _ := language.Parse(lang)
			if tag == language.Und {
				continue
			}
			builder.SetString(tag, key, str)
		}
	}

	return builder, nil
}
