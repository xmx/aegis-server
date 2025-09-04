package broker

import (
	"context"
	"net"
	"net/url"
	"strings"

	"github.com/xmx/aegis-common/transport"
	"github.com/xmx/aegis-control/datalayer/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

const (
	peerHostSuffix   = ".aegis.internal"
	agentHostSuffix  = ".agent" + peerHostSuffix
	brokerHostSuffix = ".broker" + peerHostSuffix
)

func MakeAgentURL(id bson.ObjectID, path string) *url.URL {
	sid := id.Hex()
	return MakesAgentURL(sid, path)
}

func MakesAgentURL(sid string, path string) *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   sid + agentHostSuffix,
		Path:   path,
	}
}

func MakeBrokerURL(id bson.ObjectID, path string) *url.URL {
	sid := id.Hex()
	return MakesBrokerURL(sid, path)
}

func MakesBrokerURL(sid string, path string) *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   sid + brokerHostSuffix,
		Path:   path,
	}
}

func NewHubDialer(repo repository.All, hub transport.Huber[bson.ObjectID]) Dialer {
	return &hubDialer{
		repo: repo,
		hub:  hub,
	}
}

type hubDialer struct {
	repo repository.All
	hub  transport.Huber[bson.ObjectID]
	dial net.Dialer
}

func (hd *hubDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	if host, _, _ := net.SplitHostPort(address); strings.HasSuffix(host, peerHostSuffix) {
		if sid, found := strings.CutSuffix(host, brokerHostSuffix); found {
			id, err := bson.ObjectIDFromHex(sid)
			if err != nil {
				return nil, &net.AddrError{
					Addr: address,
					Err:  "broker 标识输入错误",
				}
			}
			peer := hd.hub.Get(id)
			if peer == nil {
				return nil, &net.AddrError{
					Err:  "broker 节点未上线",
					Addr: address,
				}
			}
			mux := peer.Mux()

			return mux.Open(ctx)
		}

		sid, found := strings.CutSuffix(host, agentHostSuffix)
		if !found {
			return nil, &net.AddrError{
				Addr: address,
				Err:  "内部地址输入错误",
			}
		}

		id, err := bson.ObjectIDFromHex(sid)
		if err != nil {
			return nil, &net.AddrError{
				Addr: address,
				Err:  "agent 标识输入错误",
			}
		}

		// 通过 agent ID 查询所在的 broker 节点
		repo := hd.repo.Agent()
		opt := options.FindOne().SetProjection(bson.M{"broker": 1, "status": 1})
		agt, err := repo.FindByID(ctx, id, opt)
		if err != nil {
			return nil, &net.AddrError{
				Addr: address,
				Err:  "agent 节点不存在",
			}
		}
		if !agt.Status.Online() ||
			agt.Broker == nil ||
			agt.Broker.ID.IsZero() {
			return nil, &net.AddrError{
				Addr: address,
				Err:  "agent 节点未上线",
			}
		}

		brokID := agt.Broker.ID
		peer := hd.hub.Get(brokID)
		if peer == nil {
			return nil, &net.AddrError{
				Err:  "agent 所在的 broker 节点未上线",
				Addr: address,
			}
		}
		mux := peer.Mux()

		return mux.Open(ctx)
	}

	return nil, &net.AddrError{
		Addr: address,
		Err:  "内部地址输入错误",
	}
}

func (hd *hubDialer) isBrokerID(address string) (bson.ObjectID, bool) {
	host, _, _ := net.SplitHostPort(address)
	if host == "" {
		return bson.NewObjectID(), false
	}
	sid, found := strings.CutSuffix(host, brokerHostSuffix)
	if !found {
		return bson.NewObjectID(), false
	}
	id, err := bson.ObjectIDFromHex(sid)
	if err != nil {
		return bson.NewObjectID(), false
	}

	return id, true
}
