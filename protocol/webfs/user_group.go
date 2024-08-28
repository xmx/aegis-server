package webfs

import (
	"os/user"
	"strconv"
	"sync"
)

var userGroup = &userGroupMapper{
	groups: make(map[uint32]*user.Group, 16),
	users:  make(map[uint32]*user.User, 16),
}

type userGroupMapper struct {
	ulock  sync.RWMutex
	glock  sync.RWMutex
	users  map[uint32]*user.User
	groups map[uint32]*user.Group
}

func (m *userGroupMapper) Lookup(uid, gid uint32) (*user.User, *user.Group) {
	m.ulock.RLock()
	u := m.users[uid]
	m.ulock.RUnlock()
	if u == nil {
		u = m.slowLookupUser(uid)
	}

	m.glock.RLock()
	g := m.groups[gid]
	m.glock.RUnlock()
	if g == nil {
		g = m.slowLookupGroup(gid)
	}

	return u, g
}

func (m *userGroupMapper) slowLookupUser(uid uint32) *user.User {
	m.ulock.Lock()
	defer m.ulock.Unlock()

	if u := m.users[uid]; u != nil {
		return u
	}

	sid := strconv.FormatUint(uint64(uid), 10)
	u, err := user.LookupId(sid)
	if err != nil && u != nil {
		m.users[uid] = u
	}

	return u
}

func (m *userGroupMapper) slowLookupGroup(gid uint32) *user.Group {
	m.glock.Lock()
	defer m.glock.Unlock()

	if g := m.groups[gid]; g != nil {
		return g
	}

	sid := strconv.FormatUint(uint64(gid), 10)
	g, err := user.LookupGroupId(sid)
	if err != nil && g != nil {
		m.groups[gid] = g
	}

	return g
}
