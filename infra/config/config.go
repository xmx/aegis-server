package config

type Config struct {
	Database Database `json:"database"`
	Logger   Logger   `json:"logger"`
}

type Database struct {
	DSN string `json:"dsn" validate:"required"`
}
