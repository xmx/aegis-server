package serverd

import (
	"log/slog"
	"net/http"

	"github.com/xmx/aegis-control/linkhub"
)

type option struct {
	logger *slog.Logger
	server *http.Server
	huber  linkhub.Huber
}

type OptionBuilder struct {
	opts []func(option) option
}

func (ob OptionBuilder) List() []func(option) option {
	return ob.opts
}

func (ob OptionBuilder) Logger(l *slog.Logger) OptionBuilder {
	ob.opts = append(ob.opts, func(o option) option {
		o.logger = l
		return o
	})

	return ob
}

func (ob OptionBuilder) Server(s *http.Server) OptionBuilder {
	ob.opts = append(ob.opts, func(o option) option {
		o.server = s
		return o
	})

	return ob
}

func (ob OptionBuilder) Handler(h http.Handler) OptionBuilder {
	ob.opts = append(ob.opts, func(o option) option {
		if o.server == nil {
			o.server = &http.Server{}
		}
		o.server.Handler = h

		return o
	})

	return ob
}

func (ob OptionBuilder) Huber(h linkhub.Huber) OptionBuilder {
	ob.opts = append(ob.opts, func(o option) option {
		o.huber = h
		return o
	})

	return ob
}
