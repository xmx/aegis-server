package serverd

import (
	"github.com/xmx/aegis-common/tunnel/tundial"
	"github.com/xmx/aegis-common/tunnel/tunutil"
	"github.com/xmx/aegis-control/datalayer/repository"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func New(opts ...options.Lister[option]) tunutil.Handler {
	return nil
}

type brokerHandler struct {
	repo repository.All
}

func (bh *brokerHandler) Handle(mux tundial.Muxer) error {
	//TODO implement me
	panic("implement me")
}

func (bh *brokerHandler) auth(mux tundial.Muxer) error {

}
