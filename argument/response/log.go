package response

import "log/slog"

type LogLevel struct {
	Level slog.Level `json:"level"`
}
