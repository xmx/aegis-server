package bservice

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-server/channel/transport"
	"github.com/xmx/aegis-server/contract/brequest"
	"github.com/xmx/aegis-server/datalayer/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewSystem(repo repository.All, log *slog.Logger) *System {
	return &System{
		repo: repo,
		log:  log,
	}
}

type System struct {
	repo repository.All
	log  *slog.Logger
}

func (stm *System) NetworkCard(ctx context.Context, req *brequest.SystemNetworkCard, peer transport.Peer) error {
	id, _ := bson.ObjectIDFromHex(peer.ID())
	repo := stm.repo.Broker()
	update := bson.M{"$set": bson.M{"network_cards": req.NetworkCards}}
	_, err := repo.UpdateByID(ctx, id, update)

	return err
}
