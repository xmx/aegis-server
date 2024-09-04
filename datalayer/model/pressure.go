package model

import "time"

type Pressure struct {
	ID        int       `json:"id,string"  gorm:"column:id;primaryKey;autoIncrement"`
	PID       int       `json:"pid"        gorm:"column:pid"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;default:now(3)"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null;default:now(3)"`
}

func (Pressure) TableName() string {
	return "pressure"
}
