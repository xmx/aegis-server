package crontab

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/xmx/aegis-common/library/cronv3"
	"github.com/xmx/metrics"
)

func NewMetrics(cfg func(ctx context.Context) (pushURL string, opts *metrics.PushOptions, err error)) cronv3.Tasker {
	hostname, _ := os.Hostname()

	return &metricsTask{
		cfg:   cfg,
		label: `instance="` + hostname + `",instance_type="server"`,
	}
}

type metricsTask struct {
	cfg   func(ctx context.Context) (pushURL string, opts *metrics.PushOptions, err error)
	label string
}

func (mt *metricsTask) Info() cronv3.TaskInfo {
	return cronv3.TaskInfo{
		Name:      "上报系统指标",
		Timeout:   5 * time.Second,
		CronSched: cronv3.NewInterval(5 * time.Second),
	}
}

func (mt *metricsTask) Call(ctx context.Context) error {
	pushURL, opts, err := mt.cfg(ctx)
	if err != nil {
		return err
	}
	opts.ExtraLabels = mt.label

	return metrics.PushMetricsExt(ctx, pushURL, mt.defaultWrite, opts)
}

func (*metricsTask) defaultWrite(w io.Writer) {
	metrics.WritePrometheus(w, true)
	metrics.WriteFDMetrics(w)
}
