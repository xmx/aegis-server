package model

import "time"

type GridFile struct {
	ID        int64     `json:"id,string"  gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	Filename  string    `json:"filename"   gorm:"column:filename;varchar(255);not null;index:idx_filename;comment:文件名"`
	Length    int64     `json:"length"     gorm:"column:length;comment:文件大小"`
	Burst     uint16    `json:"-"          gorm:"column:burst;comment:分片大小"`
	SHA1      string    `json:"sha1"       gorm:"column:sha1;type:char(40);comment:SHA-1"`
	SHA256    string    `json:"sha256"     gorm:"column:sha256;type:char(64);comment:SHA-256"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null;default:now(3);comment:更新时间"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;default:now(3);comment:创建时间"`
}

func (GridFile) TableName() string {
	return "grid_file"
}
