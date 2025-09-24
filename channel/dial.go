package broker

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/xmx/aegis-common/transport"
	"github.com/xmx/aegis-control/contract/linkhub"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-server/library/memoize"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewDialer(repo repository.All, hub linkhub.Huber, dial ...*net.Dialer) transport.Dialer {
	md := &multiDialer{hub: hub, repo: repo}
	if len(dial) != 0 && dial[0] != nil {
		md.dia = dial[0]
	} else {
		md.dia = &net.Dialer{Timeout: 10 * time.Second}
	}
	md.agent = memoize.NewMap2(md.slowLoad)

	return md
}

type multiDialer struct {
	hub   linkhub.Huber
	dia   *net.Dialer
	repo  repository.All
	agent memoize.Map2[string, string, error]
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
	if !strings.HasSuffix(host, transport.BrokerHostSuffix) {
		return nil, false, nil
	}

	if agentID, found := strings.CutSuffix(host, transport.AgentBrokerHostSuffix); found {
		if conn, exx := md.lookupAgentBroker(ctx, agentID); exx != nil {
			return nil, true, exx
		} else if conn == nil {
			return nil, true, &net.AddrError{Err: "(server) no route to broker host", Addr: address}
		} else {
			return conn, true, nil
		}
	}

	brokerID, _ := strings.CutSuffix(host, transport.BrokerHostSuffix)
	if conn, exx := md.open(ctx, brokerID); exx != nil {
		return nil, true, exx
	} else if conn == nil {
		return nil, true, &net.AddrError{Err: "(server) no route to broker host", Addr: address}
	} else {
		return conn, true, nil
	}
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
			Err:  "(server) no route to agent host",
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
			Err:  "(server) no route to agent host",
		}
	}

	peer := md.hub.GetByObjectID(brk.ID)
	if peer == nil {
		return nil, true, &net.AddrError{
			Addr: address,
			Err:  "(server) no route to broker host",
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
			Err:  "(server) no route to broker host",
		}
	}

	mux := peer.Muxer()
	conn, err := mux.Open(ctx)

	return conn, true, err
}

func (md *multiDialer) slowLoad(ctx context.Context, address string) (string, error) {
	oid, err := bson.ObjectIDFromHex(address)
	if err != nil {
		return "", err
	}

	opt := options.FindOne().SetProjection(bson.M{"broker": 1, "status": 1})
	repo := md.repo.Agent()
	agt, err := repo.FindByID(ctx, oid, opt)
	if err != nil {
		return "", err
	}
	brk := agt.Broker
	if brk == nil || agt.ID.IsZero() {
		return "", mongo.ErrNoDocuments
	}

	return brk.ID.Hex(), nil
}

func (md *multiDialer) lookupAgentBroker(ctx context.Context, agentID string) (net.Conn, error) {
	oid, err := bson.ObjectIDFromHex(agentID)
	if err != nil {
		return nil, err
	}

	opt := options.FindOne().SetProjection(bson.M{"broker": 1, "status": 1})
	repo := md.repo.Agent()
	agt, err := repo.FindByID(ctx, oid, opt)
	if err != nil {
		return nil, err
	}
	brk := agt.Broker
	if brk == nil || agt.ID.IsZero() {
		return nil, mongo.ErrNoDocuments
	}

	return md.open(ctx, brk.ID.Hex())
}

func (md *multiDialer) open(ctx context.Context, brokerID string) (net.Conn, error) {
	peer := md.hub.Get(brokerID)
	if peer == nil {
		return nil, nil
	}

	return peer.Muxer().Open(ctx)
}
