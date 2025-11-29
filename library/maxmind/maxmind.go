package maxmind

import (
	"net/netip"
	"sync"
	"sync/atomic"

	"github.com/oschwald/maxminddb-golang/v2"
)

func NewDB(mmdb string) *DB {
	return &DB{
		mmdb: mmdb,
	}
}

type DB struct {
	mmdb string
	mtx  sync.Mutex
	ptr  atomic.Pointer[maxmindHolder]
}

func (mdb *DB) Close() error {
	mdb.mtx.Lock()
	mhd := mdb.ptr.Swap(nil)
	mdb.mtx.Unlock()

	if mhd != nil {
		return mhd.close()
	}

	return nil
}

func (mdb *DB) Lookup(ip netip.Addr) (*Result, error) {
	mhd := mdb.ptr.Load()
	if mhd == nil {
		mhd = mdb.loadDB()
	}

	return mhd.lookup(ip)
}

func (mdb *DB) loadDB() *maxmindHolder {
	mdb.mtx.Lock()
	defer mdb.mtx.Unlock()
	if mhd := mdb.ptr.Load(); mhd != nil {
		return mhd
	}

	mrd, err := maxminddb.Open(mdb.mmdb)
	mhd := &maxmindHolder{mrd: mrd, err: err}
	mdb.ptr.Store(mhd)

	return mhd
}

type maxmindHolder struct {
	mrd *maxminddb.Reader
	err error
}

func (mh *maxmindHolder) lookup(ip netip.Addr) (*Result, error) {
	if mh.err != nil {
		return nil, mh.err
	}

	ret := new(Result)
	val := mh.mrd.Lookup(ip)
	if err := val.Decode(ret); err != nil {
		return nil, err
	}
	ret.Traits.IPAddress = ip
	ret.Traits.Network = val.Prefix()

	return ret, nil
}

func (mh *maxmindHolder) close() error {
	if mrd := mh.mrd; mrd != nil {
		return mrd.Close()
	}

	return nil
}
