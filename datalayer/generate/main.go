package main

import (
	"github.com/xmx/aegis-server/datalayer/model"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		Mode:    gen.WithDefaultQuery,
		OutPath: "datalayer/query",
	})
	g.ApplyBasic(model.All()...)
	g.Execute()
}
