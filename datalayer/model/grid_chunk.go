package model

import "time"

type GridChunk struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	FileID    int64     `gorm:"column:file_id;index:idx_file_id_sequence"`
	Sequence  int64     `gorm:"column:sequence;index:idx_file_id_sequence"`
	Data      []byte    `gorm:"column:data;blob"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:now(3)"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now(3)"`
}

func (GridChunk) TableName() string {
	return "grid_chunk"
}
