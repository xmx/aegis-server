package main

import (
	"github.com/xmx/aegis-server/datalayer/model"
	"gorm.io/gen"
)

func main() {
	tables := []any{
		model.ConfigCertificate{},
		model.ConfigServer{},
	}

	g := gen.NewGenerator(gen.Config{
		Mode:    gen.WithDefaultQuery,
		OutPath: "datalayer/query",
	})
	g.ApplyBasic(tables...)
	g.Execute()
}
