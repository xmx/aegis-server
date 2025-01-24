package main

import (
	"github.com/xmx/aegis-server/datalayer/model"
	"gorm.io/gen"
)

func main() {
	c := gen.Config{OutPath: "datalayer/query"}
	g := gen.NewGenerator(c)
	g.ApplyBasic(model.All()...)
	g.Execute()
}
