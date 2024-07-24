package alpnx

type Negotiation interface {
	Lookup(proto string) Application
	Protocols() []string
}

func NewNegotiation(apps []Application) Negotiation {
	size := len(apps)
	hm := make(map[string]Application, size)
	ps := make([]string, 0, size)
	for _, app := range apps {
		proto := app.Proto()
		if proto == "" || app == nil {
			continue
		}
		if _, ok := hm[proto]; ok {
			continue
		}
		ps = append(ps, proto)
		hm[proto] = app
	}

	return &negotiationLayer{
		hm: hm,
		ps: ps,
	}
}

type negotiationLayer struct {
	hm map[string]Application
	ps []string
}

func (nl *negotiationLayer) Lookup(proto string) Application {
	return nl.hm[proto]
}

func (nl *negotiationLayer) Protocols() []string {
	return nl.ps
}
