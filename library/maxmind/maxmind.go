package maxmind

import (
	"context"
	"net/netip"
	"sync"
	"sync/atomic"

	"github.com/oschwald/maxminddb-golang/v2"
)

type File string

func (m File) LoadFile(context.Context) (string, error) {
	return string(m), nil
}

type FileLoader interface {
	LoadFile(ctx context.Context) (file string, err error)
}

func NewDB(mfl FileLoader) *DB {
	return &DB{
		mfl: mfl,
	}
}

type DB struct {
	mfl FileLoader
	mtx sync.Mutex
	ptr atomic.Pointer[maxmindHolder]
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

func (mdb *DB) Lookup(ctx context.Context, ip netip.Addr) (*Result, error) {
	mhd := mdb.ptr.Load()
	if mhd == nil {
		mhd = mdb.loadDB(ctx)
	}

	return mhd.lookup(ip)
}

func (mdb *DB) loadDB(ctx context.Context) *maxmindHolder {
	mdb.mtx.Lock()
	defer mdb.mtx.Unlock()
	if mhd := mdb.ptr.Load(); mhd != nil {
		return mhd
	}

	mhd := new(maxmindHolder)
	file, err := mdb.mfl.LoadFile(ctx)
	if err != nil {
		mhd.err = err
		if te, ok := err.(interface{ Timeout() bool }); !ok || !te.Timeout() {
			mdb.ptr.Store(mhd)
		}

		return mhd
	}

	mrd, err := maxminddb.Open(file)
	if err != nil {
		mhd.err = err
		mdb.ptr.Store(mhd)
		return mhd
	}

	mhd.mrd = mrd
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
