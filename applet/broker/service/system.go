package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-server/channel/broker"
	"github.com/xmx/aegis-server/contract/brequest"
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

func (stm *System) NetworkCard(ctx context.Context, req *brequest.SystemNetworkCard, peer broker.Peer) error {
	id := peer.ID()
	repo := stm.repo.Broker()
	update := bson.M{"$set": bson.M{"network_cards": req.NetworkCards}}
	_, err := repo.UpdateByID(ctx, id, update)

	return err
}
