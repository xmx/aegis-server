package serverd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/xmx/aegis-common/muxlink/muxproto"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-control/linkhub"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMixedDialer(hub linkhub.Huber, agt AgentOpener, back muxproto.Dialer) muxproto.Dialer {
	return &mixedDialer{
		hub:  hub,
		agt:  agt,
		back: back,
	}
}

type mixedDialer struct {
	hub  linkhub.Huber
	agt  AgentOpener
	back muxproto.Dialer
}

func (m *mixedDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	if sid, domain, found := strings.Cut(host, "."); found {
		if m.hub != nil && domain == m.hub.Domain() {
			if peer := m.hub.Get(host); peer != nil {
				mux := peer.Muxer()
				return mux.Open(ctx)
			}

			return nil, &net.OpError{
				Op:   "find",
				Net:  "broker",
				Addr: &net.UnixAddr{Net: network, Name: address},
				Err:  errors.New("节点已离线或未注册"),
			}
		}

		if m.agt != nil && domain == m.agt.Domain() {
			id, _ := bson.ObjectIDFromHex(sid)
			if id.IsZero() {
				return nil, &net.OpError{
					Op:   "parse",
					Net:  "agent",
					Addr: &net.UnixAddr{Net: network, Name: address},
					Err:  errors.New("节点主机标识格式错误"),
				}
			}

			return m.agt.Open(ctx, id)
		}
	}

	if m.back != nil {
		return m.back.DialContext(ctx, network, address)
	}

	return nil, &net.OpError{
		Op:   "dial",
		Net:  network,
		Addr: &net.UnixAddr{Net: network, Name: address},
		Err:  net.UnknownNetworkError("没有找到适配的拨号器"),
	}
}

type AgentOpener interface {
	Open(ctx context.Context, agentID bson.ObjectID) (net.Conn, error)
	Domain() string
}

func NewAgentOpener(repo repository.All, hub linkhub.Huber, domain string) AgentOpener {
	return &agentDialer{
		domain: domain,
		repo:   repo,
		hub:    hub,
	}
}

type agentDialer struct {
	domain string
	repo   repository.All
	hub    linkhub.Huber
}

func (a *agentDialer) Open(ctx context.Context, agentID bson.ObjectID) (net.Conn, error) {
	// 1. 通过 agent id 查询所在的 broker。
	repo := a.repo.Agent()
	opt := options.FindOne().SetProjection(bson.M{"broker": 1, "status": 1})
	agt, err := repo.FindByID(ctx, agentID, opt)
	if err != nil {
		return nil, a.agentNotExists(agentID, err)
	} else if !agt.Status || agt.Broker == nil {
		return nil, a.agentOffline(agentID)
	}

	brok := agt.Broker
	if brok == nil || brok.ID.IsZero() {
		return nil, a.agentException(agentID)
	}

	peer := a.hub.GetID(brok.ID)
	if peer != nil {
		mux := peer.Muxer()
		return mux.Open(ctx)
	}

	return nil, a.brokerOffline(brok.ID, agentID)
}

func (a *agentDialer) Domain() string {
	return a.domain
}

func (*agentDialer) agentNotExists(agentID bson.ObjectID, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = errors.New("节点不存在")
	} else {
		err = fmt.Errorf("查询节点信息错误：%w", err)
	}

	return &net.OpError{
		Op:   "find",
		Net:  "agent",
		Addr: &net.UnixAddr{Net: "tunnel", Name: agentID.Hex()},
		Err:  err,
	}
}

func (*agentDialer) agentException(agentID bson.ObjectID) error {
	return &net.OpError{
		Op:   "find",
		Net:  "agent",
		Addr: &net.UnixAddr{Net: "tunnel", Name: agentID.Hex()},
		Err:  errors.New("节点状态异常"),
	}
}

func (*agentDialer) agentOffline(agentID bson.ObjectID) error {
	return &net.OpError{
		Op:   "lookup",
		Net:  "agent",
		Addr: &net.UnixAddr{Net: "tunnel", Name: agentID.Hex()},
		Err:  errors.New("节点已离线"),
	}
}

func (*agentDialer) brokerOffline(brokerID, agentID bson.ObjectID) error {
	return &net.OpError{
		Op:     "find",
		Net:    "broker",
		Source: &net.UnixAddr{Net: "tunnel", Name: brokerID.Hex()},
		Addr:   &net.UnixAddr{Net: "tunnel", Name: agentID.Hex()},
		Err:    errors.New("节点已离线"),
	}
}
