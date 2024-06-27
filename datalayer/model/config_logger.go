package model

import "time"

type ConfigLogger struct {
	ID        int64     `json:"id,string"  gorm:"column:id;primaryKey;autoIncrement"`
	Enabled   bool      `json:"enabled"    gorm:"column:enabled"`
	Level     string    `json:"level"      gorm:"column:level;type:varchar(10)"`
	Terminal  bool      `json:"terminal"   gorm:"column:terminal"`
	Filename  string    `json:"filename"   gorm:"column:filename;type:varchar(255)"`
	MaxAge    int64     `json:"max_age"    gorm:"column:max_age;type:bigint"`
	MaxBackup int64     `json:"max_backup" gorm:"column:max_backup;type:bigint"`
	Localtime bool      `json:"localtime"  gorm:"column:localtime"`
	Compress  bool      `json:"compress"   gorm:"column:compress"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null;default:now(3)"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;default:now(3)"`
}

func (ConfigLogger) TableName() string {
	return "config_logger"
}
