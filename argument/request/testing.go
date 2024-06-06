package request

type TestingListen struct {
	Addr string `json:"addr" validate:"required"`
}

type TestingTiDB struct {
	DSN string `json:"dsn" validate:"required"`
}
