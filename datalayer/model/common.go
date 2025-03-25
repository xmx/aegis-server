package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Duration struct {
	Nanosecond int64  `json:"nanosecond" bson:"nanosecond"` // 纳秒
	Formatted  string `json:"formatted"  bson:"formatted"`  // 格式化后的时间
}

type Operator struct {
	ID   bson.ObjectID `json:"id"   bson:"id"`   // 用户 ID
	Name string        `json:"name" bson:"name"` // 用户名
}

type FrontendDuration struct {
	Days        int           `json:"days,omitempty"        bson:"days,omitempty"`
	Hours       int           `json:"hours,omitempty"       bson:"hours,omitempty"`
	Minutes     int           `json:"minutes,omitempty"     bson:"minutes,omitempty"`
	Seconds     int           `json:"seconds,omitempty"     bson:"seconds,omitempty"`
	Millis      int           `json:"millis,omitempty"      bson:"millis,omitempty"`
	Nanos       int           `json:"nanos,omitempty"       bson:"nanos,omitempty"`
	Nanoseconds time.Duration `json:"nanoseconds,omitempty" bson:"nanoseconds,omitempty"`
	Formatted   string        `json:"formatted,omitempty"   bson:"formatted,omitempty"`
}

func FromDuration(du time.Duration) FrontendDuration {
	totalNanoseconds := du.Nanoseconds()
	totalMilliseconds := totalNanoseconds / int64(time.Millisecond)
	totalSeconds := totalNanoseconds / int64(time.Second)
	totalMinutes := totalNanoseconds / int64(time.Minute)
	totalHours := totalNanoseconds / int64(time.Hour)
	totalDays := totalNanoseconds / int64(24*time.Hour)

	return FrontendDuration{
		Days:        int(totalDays),
		Hours:       int(totalHours % 24),
		Minutes:     int(totalMinutes % 60),
		Seconds:     int(totalSeconds % 60),
		Millis:      int(totalMilliseconds % 1000),
		Nanos:       int(totalNanoseconds % int64(time.Millisecond)),
		Nanoseconds: du,
		Formatted:   du.String(),
	}
}
