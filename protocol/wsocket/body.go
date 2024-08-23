package wsocket

import "encoding/json"

const (
	KindStdin  Kind = "stdin"
	KindStdout Kind = "stdout"
	KindStderr Kind = "stderr"
	KindError  Kind = "error"
	KindConfig Kind = "config"
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

type Recv struct {
	Kind Kind   `json:"kind"`
	Data []byte `json:"data"`
}

func (r Recv) JSON(v any) error {
	return json.Unmarshal(r.Data, v)
}
