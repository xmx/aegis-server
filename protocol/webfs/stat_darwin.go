package webfs

import (
	"syscall"
	"time"
)

func readStat(v any) sysStat {
	st, ok := v.(*syscall.Stat_t)
	if !ok || st == nil {
		return sysStat{}
	}

	stat := sysStat{
		AccessedAt: formatTime(st.Atimespec),
		CreatedAt:  formatTime(st.Ctimespec),
		UpdatedAt:  formatTime(st.Mtimespec),
	}
	u, g := userGroup.Lookup(st.Uid, st.Gid)
	if u != nil {
		stat.User = u.Name
	}
	if g != nil {
		stat.Group = g.Name
	}

	return stat
}

func formatTime(at syscall.Timespec) time.Time {
	sec, nano := at.Unix()
	return time.Unix(sec, nano)
}
