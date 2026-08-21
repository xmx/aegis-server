package crontab

import (
	"slices"
	"time"

	"github.com/robfig/cron/v3"
)

// NewFixedTimes 定点任务，在指定的几个时间点执行。
//
// 例如仅在：
//
//	2025-01-01 00:00:00
//	2025-03-01 00:00:00
//	2025-03-15 00:00:00
//
// 执行三次就完事的任务。
func NewFixedTimes(times []time.Time) cron.Schedule {
	slices.SortFunc(times, func(a, b time.Time) int {
		return a.Compare(b)
	})

	return &fixedTimes{
		times: times,
	}
}

type fixedTimes struct {
	times []time.Time
	index int
}

func (st *fixedTimes) Next(now time.Time) time.Time {
	for st.index < len(st.times) {
		at := st.times[st.index]
		st.index++

		if at.After(now) {
			return at
		}
	}

	return time.Time{}
}
