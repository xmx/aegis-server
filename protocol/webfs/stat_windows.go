package webfs

import (
	"syscall"
	"time"
)

func readStat(v any) sysStat {
	d, ok := v.(*syscall.Win32FileAttributeData)
	if !ok || d == nil {
		return sysStat{}
	}

	stat := sysStat{
		AccessedAt: formatTime(d.LastAccessTime),
		CreatedAt:  formatTime(d.CreationTime),
		UpdatedAt:  formatTime(d.LastWriteTime),
	}

	return stat
}

func formatTime(at syscall.Filetime) time.Time {
	nano := at.Nanoseconds()
	return time.Unix(0, nano)
}
