package request

import "log/slog"

type LogWatch struct {
	Format string `json:"format" query:"format" validate:"omitempty,oneof=json text"`             // 默认 JSON
	Level  string `json:"level"  query:"level"  validate:"omitempty,oneof=DEBUG INFO WARN ERROR"` // 默认 INFO
}

func (lw LogWatch) JSONFormat() bool {
	return lw.Format != "text"
}

func (lw LogWatch) LevelVar() *slog.LevelVar {
	lvl := new(slog.LevelVar)
	_ = lvl.UnmarshalText([]byte(lw.Level))

	return lvl
}
