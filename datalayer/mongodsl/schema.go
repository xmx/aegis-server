package mongodsl

import (
	"reflect"
	"strings"
	"time"

	"github.com/xmx/aegis-server/datalayer/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type KeyInfo struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type CollectionInfo struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Keys        []KeyInfo `json:"keys"`
}

func ForCollection(zero repository.CollectionInfer) (*CollectionInfo, error) {
	name, comment := zero.CollectionInfo()

	rtyp := reflect.TypeOf(zero)
	keys := parseStruct(rtyp, "")

	return &CollectionInfo{
		Name:        name,
		Description: comment,
		Keys:        keys,
	}, nil
}

func parseStructField(sf reflect.StructField, prefix string) []KeyInfo {
	if !sf.IsExported() {
		return nil
	}

	btag := sf.Tag.Get("bson")
	if btag == "-" {
		return nil
	}
	stag := sf.Tag.Get("jsonschema")
	if stag == "-" {
		return nil
	}
	name, _, _ := strings.Cut(btag, ",")
	if name == "" {
		name = strings.ToLower(sf.Name)
	}
	if prefix != "" {
		name = prefix + "." + name
	}

	key := KeyInfo{Key: name, Description: stag}

	sft := sf.Type
	sfkind := sf.Type.Kind()
	if sfkind == reflect.Array && reflect.TypeFor[bson.ObjectID]() == sft {
		key.Type = "oid"
		return []KeyInfo{key}
	}
	if sfkind == reflect.Slice || sfkind == reflect.Array {
		sft = sft.Elem()
	}
	if sft.Kind() == reflect.Pointer {
		sft = sft.Elem()
	}

	switch sft.Kind() {
	case reflect.String:
		key.Type = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		key.Type = "number"
	case reflect.Bool:
		key.Type = "boolean"
	case reflect.Struct:
		if reflect.TypeFor[time.Time]() == sft {
			key.Type = "datetime"
		}
	default:
	}

	if key.Type != "" {
		return []KeyInfo{key}
	}

	return parseStruct(sft, name)
}

func parseStruct(t reflect.Type, prefix string) []KeyInfo {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	keys := make([]KeyInfo, 0, 10)
	for field := range t.Fields() {
		subs := parseStructField(field, prefix)
		keys = append(keys, subs...)
	}

	return keys
}
