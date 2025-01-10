package response

import "log/slog"

type LogLevel struct {
	Log  slog.Level `json:"log"`
	Gorm slog.Level `json:"gorm"`
}
