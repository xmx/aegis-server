package service

import (
	"context"

	"github.com/oschwald/geoip2-golang/v2"
	"github.com/xmx/aegis-control/library/memoize"
)

func NewMaxmind() *Maxmind {
	mdb := new(Maxmind)
	mdb.cache = memoize.NewCache2(mdb.slowLoad)

	return mdb
}

type Maxmind struct {
	cache memoize.Cache2[*geoip2.Reader, error]
}

func (m *Maxmind) ReadMaxmind(ctx context.Context) (*geoip2.Reader, error) {
	return m.cache.Load(ctx)
}

func (m *Maxmind) slowLoad(context.Context) (*geoip2.Reader, error) {
	return geoip2.Open("resources/mmdb/GeoLite2-City.mmdb")
}
