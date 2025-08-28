package transport

import "sync"

func NewHub() Huber {
	return &mapHub{
		peers: make(map[string]Peer, 32),
	}
}

type Huber interface {
	// Get 通过 ID 获取 peer。
	Get(id string) Peer

	// Del 通过删除节点，并返回删除的数据，nil 代表之前无数据。
	Del(id string) Peer

	// Put 存放节点并返回原来的数据（如果存在的话）。
	Put(p Peer) Peer

	// PutIfAbsent 当此 ID 没有时才存放节点，并返回是否放入成功。
	PutIfAbsent(p Peer) bool
}

type mapHub struct {
	mutex sync.RWMutex
	peers map[string]Peer
}

func (mh *mapHub) Get(id string) Peer {
	return mh.peers[id]
}

func (mh *mapHub) Del(id string) Peer {
	mh.mutex.Lock()
	peer := mh.peers[id]
	delete(mh.peers, id)
	mh.mutex.Unlock()

	return peer
}

func (mh *mapHub) Put(p Peer) Peer {
	id := p.ID()

	mh.mutex.Lock()
	last := mh.peers[id]
	mh.peers[id] = p
	mh.mutex.Unlock()

	return last
}

func (mh *mapHub) PutIfAbsent(p Peer) bool {
	id := p.ID()

	mh.mutex.Lock()
	last := mh.peers[id]
	absent := last == nil
	if absent {
		mh.peers[id] = p
	}
	mh.mutex.Unlock()

	return absent
}
