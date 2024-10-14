package model

import "time"

type GridChunk struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	FileID    int64     `gorm:"column:file_id;index:idx_file_id_sequence;comment:文件ID"`
	Sequence  int64     `gorm:"column:sequence;index:idx_file_id_sequence;comment:分片序号"`
	Data      []byte    `gorm:"column:data;type:blob;comment:分片内容"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:now(3);comment:更新时间"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now(3);comment:创建时间"`
}

func (GridChunk) TableName() string {
	return "grid_chunk"
}
