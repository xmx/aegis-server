package launch

import (
	"github.com/xmx/aegis-server/datalayer/model"
	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		model.ConfigCertificate{},
		model.ConfigServer{},
		model.GridChunk{},
		model.GridFile{},
		model.Pressure{},
	)
}
