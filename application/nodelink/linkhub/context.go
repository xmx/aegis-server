package linkhub

import "context"

type contextKey struct {
	key string
}

var muxContextKey = &contextKey{key: "mux-context-key"}

func withContext(cli *clientMUX) context.Context {
	return context.WithValue(context.Background(), muxContextKey, cli)
}

func FromContext(ctx context.Context) Muxer {
	cli, ok := ctx.Value(muxContextKey).(*clientMUX)
	if !ok || cli == nil {
		return nil
	}

	return cli
}
