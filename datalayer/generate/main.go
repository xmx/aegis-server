package main

import (
	"github.com/xmx/aegis-server/datalayer/model"
	"gorm.io/gen"
)

func main() {
	tables := []any{
		model.ConfigCertificate{},
		model.ConfigServer{},
		model.GridChunk{},
		model.GridFile{},
		model.OAuthConfig{},
		model.Menu{},
		model.Oplog{},
		model.Role{},
		model.RoleMenu{},
		model.User{},
	}

	g := gen.NewGenerator(gen.Config{
		Mode:    gen.WithDefaultQuery,
		OutPath: "datalayer/query",
	})
	g.ApplyBasic(tables...)
	g.Execute()
}
