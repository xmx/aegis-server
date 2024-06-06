package request

type TestingListen struct {
	Addr string `json:"addr" validate:"required"`
}
