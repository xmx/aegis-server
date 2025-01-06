package launch

import (
	"github.com/xmx/aegis-server/datalayer/model"
	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(model.All()...)
}
