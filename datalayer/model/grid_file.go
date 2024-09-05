package model

import "time"

type GridFile struct {
	ID        int64     `json:"id,string"  gorm:"column:id;primaryKey;autoIncrement"`                       // 文件 ID。
	Filename  string    `json:"filename"   gorm:"column:filename;varchar(255);not null;index:idx_filename"` // 文件名。
	Length    int64     `json:"length"     gorm:"column:length"`                                            // 文件总大小。
	Burst     uint16    `json:"burst"      gorm:"column:burst"`                                             // 文件分片大小。
	SHA1      string    `json:"sha1"       gorm:"column:sha1;type:char(40)"`                                // SHA-1
	SHA256    string    `json:"sha256"     gorm:"column:sha256;type:char(64)"`                              // SHA-256
	Done      bool      `json:"done"       gorm:"column:done"`                                              // 文件是否已就绪。
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null;default:now(3)"`                // 文件更新时间。
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;default:now(3)"`                // 创建时间。
}

func (GridFile) TableName() string {
	return "grid_file"
}
