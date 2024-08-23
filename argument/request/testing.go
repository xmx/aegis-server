package request

type TestingListen struct {
	Addr string `json:"addr" validate:"required"`
}

type TestingTiDB struct {
	DSN string `json:"dsn" validate:"required"`
}

type TestingCert struct {
	Cert string `json:"cert" validate:"required"`
	Pkey string `json:"pkey" validate:"required"`
}

type TestingConfig struct {
	TiDB TestingTiDB `json:"tidb"`
}
