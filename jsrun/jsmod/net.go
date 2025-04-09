package jsmod

import (
	"net"

	"github.com/xmx/aegis-server/jsrun/jsvm"
)

func NewNet() jsvm.GlobalRegister {
	return new(stdNet)
}

type stdNet struct {
	eng jsvm.Engineer
}

func (s *stdNet) RegisterGlobal(eng jsvm.Engineer) error {
	s.eng = eng
	fns := map[string]any{
		"listen": s.listen,
	}
	return eng.Runtime().Set("net", fns)
}

func (s *stdNet) listen(network, address string) (net.Listener, error) {
	lis, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	s.eng.AddFinalizer(lis.Close)

	return lis, nil
}
