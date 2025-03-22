package cronv3

import (
	"slices"
	"time"

	"github.com/robfig/cron/v3"
)

// NewSpecificTimes 定点任务，在指定的几个时间点执行。
//
// 例如仅在：
//
//	2025-01-01 00:00:00
//	2025-03-01 00:00:00
//	2025-03-15 00:00:00
//
// 执行三次就完事的任务。
func NewSpecificTimes(times []time.Time) cron.Schedule {
	slices.SortFunc(times, func(a, b time.Time) int {
		return int(a.Sub(b))
	})

	return &specificTimes{
		times: times,
	}
}

type specificTimes struct {
	times []time.Time
}

func (st *specificTimes) Next(now time.Time) time.Time {
	for idx, at := range st.times {
		if at.After(now) {
			st.times = st.times[idx:]
			return at
		}
	}

	return time.Time{}
}

// NewPeriodicallyTimes 根据时间间隔定期执行。
//
// 假设用户有这样的一个需求：
// 有个定时任务要每小时执行一次。
// 那 cron 表达式要写：0 0 * * * * 还是 0 * * * * 呢？
// 这要取决于 cron/v3 在初始化时允许秒级任务还是分钟级任务。
func NewPeriodicallyTimes(du time.Duration) cron.Schedule {
	if du < time.Minute {
		du = time.Minute
	}

	return &periodicallyTimes{
		du: du,
	}
}

type periodicallyTimes struct {
	du time.Duration
}

func (pt *periodicallyTimes) Next(now time.Time) time.Time {
	return now.Add(pt.du)
}
