package broker

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/xmx/aegis-common/transport"
	"github.com/xmx/aegis-control/contract/linkhub"
	"github.com/xmx/aegis-control/datalayer/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

func NewDialer(repo repository.All, hub linkhub.Huber, dial ...*net.Dialer) Dialer {
	md := &multiDialer{hub: hub, repo: repo}
	if len(dial) != 0 && dial[0] != nil {
		md.dia = dial[0]
	} else {
		md.dia = &net.Dialer{Timeout: 10 * time.Second}
	}

	return md
}

type multiDialer struct {
	hub  linkhub.Huber
	dia  *net.Dialer
	repo repository.All
}

func (md *multiDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	if conn, match, err := md.matchTunnel(ctx, address); match {
		return conn, err
	}

	return md.dia.DialContext(ctx, network, address)
}

func (md *multiDialer) matchTunnel(ctx context.Context, address string) (net.Conn, bool, error) {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return nil, false, err
	}

	if conn, match, exx := md.matchBroker(ctx, host); match {
		return conn, true, exx
	}

	return md.matchAgent(ctx, host)
}

func (md *multiDialer) matchAgent(ctx context.Context, address string) (net.Conn, bool, error) {
	host, found := strings.CutSuffix(address, transport.AgentHostSuffix)
	if !found {
		return nil, false, nil
	}

	id, err := bson.ObjectIDFromHex(host)
	if err != nil {
		return nil, true, &net.AddrError{
			Addr: address,
			Err:  "no route to agent host",
		}
	}

	repo := md.repo.Agent()
	opt := options.FindOne().SetProjection(bson.M{"broker": 1, "status": 1})
	agt, err := repo.FindByID(ctx, id, opt)
	if err != nil {
		return nil, true, &net.AddrError{
			Addr: address,
			Err:  err.Error(),
		}
	}
	brk := agt.Broker
	if brk == nil || agt.ID.IsZero() {
		return nil, true, &net.AddrError{
			Addr: address,
			Err:  "no route to agent host",
		}
	}

	peer := md.hub.GetByObjectID(brk.ID)
	if peer == nil {
		return nil, true, &net.AddrError{
			Addr: address,
			Err:  "no route to broker host",
		}
	}

	mux := peer.Muxer()
	conn, err := mux.Open(ctx)

	return conn, true, err
}

func (md *multiDialer) matchBroker(ctx context.Context, address string) (net.Conn, bool, error) {
	host, found := strings.CutSuffix(address, transport.BrokerHostSuffix)
	if !found {
		return nil, false, nil
	}

	peer := md.hub.Get(host)
	if peer == nil {
		return nil, true, &net.AddrError{
			Addr: address,
			Err:  "no route to broker host",
		}
	}

	mux := peer.Muxer()
	conn, err := mux.Open(ctx)

	return conn, true, err
}
