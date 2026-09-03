package linkhub

import (
	"context"
	"net"
	"time"

	"github.com/xmx/aegis-common/muxutil"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/muxconn"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/time/rate"
)

type Info struct {
	ID         bson.ObjectID `json:"id"`
	MachineID  string        `json:"machine_id"`
	Inet       string        `json:"inet"`
	Semver     model.Semver  `json:"semver"`
	GOOS       string        `json:"goos,omitzero"`
	GOARCH     string        `json:"goarch,omitzero"`
	PID        int           `json:"pid,omitzero"`
	Args       []string      `json:"args,omitzero"`
	Hostname   string        `json:"hostname,omitzero"`
	Workdir    string        `json:"workdir,omitzero"`
	Executable string        `json:"executable,omitzero"`
	Environ    []string      `json:"environ,omitzero"`
}

type Muxer interface {
	muxconn.Muxer
	Info() Info
}

type clientMUX struct {
	mux muxconn.Muxer
	inf Info
}

func newClientMUX(mux muxconn.Muxer, inf *muxutil.Request, dat *model.Agent) *clientMUX {
	info := Info{
		ID:         dat.ID,
		MachineID:  inf.MachineID,
		Inet:       inf.Inet,
		Semver:     model.ParseSemver(inf.Version),
		GOOS:       inf.GOOS,
		GOARCH:     inf.GOARCH,
		PID:        inf.PID,
		Args:       inf.Args,
		Hostname:   inf.Hostname,
		Workdir:    inf.Workdir,
		Executable: inf.Executable,
		Environ:    inf.Environ,
	}

	return &clientMUX{mux: mux, inf: info}
}

func (c *clientMUX) Accept() (net.Conn, error) {
	return c.mux.Accept()
}

func (c *clientMUX) Close() error {
	return c.mux.Close()
}

func (c *clientMUX) Addr() net.Addr {
	return c.mux.Addr()
}

func (c *clientMUX) Open(ctx context.Context) (net.Conn, error) {
	return c.mux.Open(ctx)
}

func (c *clientMUX) RemoteAddr() net.Addr {
	return c.mux.RemoteAddr()
}

func (c *clientMUX) IsClosed() bool {
	return c.mux.IsClosed()
}

func (c *clientMUX) Limit() rate.Limit {
	return c.mux.Limit()
}

func (c *clientMUX) SetLimit(bps rate.Limit) bool {
	return c.mux.SetLimit(bps)
}

func (c *clientMUX) Streams() []muxconn.Streamer {
	return c.mux.Streams()
}

func (c *clientMUX) NumStreams() (int64, int64) {
	return c.mux.NumStreams()
}

func (c *clientMUX) Traffic() (rx, tx uint64) {
	return c.mux.Traffic()
}

func (c *clientMUX) Library() (name, module string) {
	return c.mux.Library()
}

func (c *clientMUX) ConnectedAt() time.Time {
	return c.mux.ConnectedAt()
}

func (c *clientMUX) Info() Info {
	return c.inf
}
