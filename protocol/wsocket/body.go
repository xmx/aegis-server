package wsocket

const (
	KindStdout Kind = "stdout"
	KindStderr Kind = "stderr"
	KindError  Kind = "error"
)

type Kind string

type Body struct {
	Kind Kind `json:"kind"`
	Data any  `json:"data"`
}

func StdoutBody(data any) *Body {
	return &Body{Kind: KindStdout, Data: data}
}

func StderrBody(data any) *Body {
	return &Body{Kind: KindStderr, Data: data}
}

func ErrorBody(err error) *Body {
	return &Body{Kind: KindError, Data: err.Error()}
}

func SErrorBody(msg string) *Body {
	return &Body{Kind: KindError, Data: msg}
}
