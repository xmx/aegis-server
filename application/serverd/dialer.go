package serverd

import (
	"context"
	"net"
	"strings"

	"github.com/xmx/aegis-common/muxlink/muxproto"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-control/linkhub"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewFindAgentDialer(suffix string, hub linkhub.Huber, repo repository.All) muxproto.Dialer {
	return &findAgentDialer{
		repo:   repo,
		huber:  hub,
		suffix: suffix,
	}
}

type findAgentDialer struct {
	repo   repository.All
	huber  linkhub.Huber
	suffix string
}

func (fad *findAgentDialer) Dial(network, address string) (net.Conn, error) {
	return fad.DialContext(context.Background(), network, address)
}

func (fad *findAgentDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, _, _ := net.SplitHostPort(address)
	if host == "" {
		return nil, nil
	}
	agentID, found := strings.CutSuffix(host, fad.suffix)
	if !found {
		return nil, nil
	}

	aid, err := bson.ObjectIDFromHex(agentID)
	if err != nil {
		return nil, net.InvalidAddrError("agent id 无效：" + agentID)
	}

	opt := options.FindOne().SetProjection(bson.M{"broker": 1})
	repo := fad.repo.Agent()
	agt, err := repo.FindByID(ctx, aid, opt)
	if err != nil || agt.Broker == nil {
		return nil, &net.OpError{
			Op:   "lookup",
			Net:  network,
			Addr: &net.UnixAddr{Net: network, Name: address},
			Err:  err,
		}
	}

	brok := agt.Broker
	if peer := fad.huber.GetID(brok.ID); peer != nil {
		mux := peer.Muxer()
		return mux.Open(ctx)
	}

	// 此时是 broker 节点不通，改写错误提示。
	brokID := brok.ID.Hex()
	brokHost := muxproto.ServerToBrokerURL(brokID, "").Host
	return nil, &net.OpError{
		Op:   "dial",
		Net:  network,
		Addr: &net.UnixAddr{Net: network, Name: brokHost},
		Err:  net.UnknownNetworkError("no route to host"),
	}
}
