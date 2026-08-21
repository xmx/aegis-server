package crontab

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
)

type TaskInfo struct {
	ID        string        `json:"id"                  jsonschema:"定时任务唯一标识"` // 任务唯一标识，如果为空则会通过反射获取 FQDN 作为唯一标识，此时同一个 struct 不同实例 FQDN 是一样的。
	Name      string        `json:"name"                jsonschema:"任务名，简要描述任务功能"`
	Immediate bool          `json:"immediate,omitzero"  jsonschema:"添加到定时任务管理器时是否立即执行"`
	Timeout   time.Duration `json:"timeout"             jsonschema:"每次执行超时时间，小于等于零表示不限制超时时间"`
	CronSched cron.Schedule `json:"cron_sched,omitzero" jsonschema:"触发周期，等同于 cron_spec，如果 cron_sched 和 cron_spec 同时存在时，则以 cron_sched 为准"`
	CronSpec  string        `json:"cron_spec,omitzero"  jsonschema:"触发周期，等同于 cron_sched，但是当 cron_sched 为空时才会生效。"`
}

func (t TaskInfo) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("id", t.ID),
		slog.String("name", t.Name),
		slog.Bool("immediate", t.Immediate),
		slog.Duration("timeout", t.Timeout),
	}
	if str := t.CronSpec; str != "" {
		attrs = append(attrs, slog.String("cron_spec", str))
	}

	return slog.GroupValue(attrs...)
}

type EntryInfo struct {
	Entry cron.Entry
	Info  TaskInfo
}

// Tasker 定时任务。
type Tasker interface {
	// Info 定时任务信息。
	Info() TaskInfo

	// Call 定时触发的实际函数。
	Call(ctx context.Context) error
}

type scheduledTask struct {
	log     *slog.Logger
	entryID cron.EntryID
	work    atomic.Bool // 是否正在执行中
	info    TaskInfo    // 任务信息
	exec    func(context.Context) error
}

func (s *scheduledTask) Run() {
	_ = s.Call(context.Background())
}

func (s *scheduledTask) Info() TaskInfo {
	return s.info
}

func (s *scheduledTask) Call(ctx context.Context) error {
	if !s.work.CompareAndSwap(false, true) {
		return nil
	}
	defer s.work.Store(false)

	startedAt := time.Now()
	attrs := []any{"info", s.info, "started_at", startedAt}
	s.log.Debug("定时任务开始执行", attrs...)
	panicked, err := s.safeCall(ctx)
	finishedAt := time.Now()
	elapsed := finishedAt.Sub(startedAt)
	attrs = append(attrs, "finished_at", finishedAt, "elapsed", elapsed)
	if err == nil {
		s.log.Debug("定时任务执行完毕", attrs...)
		return nil
	}
	attrs = append(attrs, "err", err)
	if panicked {
		s.log.Error("定时任务执行发生 panic", attrs...)
	} else {
		s.log.Warn("定时任务执行出错", attrs...)
	}

	return err
}

func (s *scheduledTask) safeCall(ctx context.Context) (panicked bool, err error) {
	if du := s.info.Timeout; du > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, du)
		defer cancel()
	}

	defer func() {
		if v := recover(); v != nil {
			panicked = true
			stack := debug.Stack()
			err = fmt.Errorf("panic: %v\n%s", v, stack)
			_, _ = fmt.Fprint(os.Stderr, err.Error())
		}
	}()

	err = s.exec(ctx)

	return
}

func wrapFunc(cmd func()) func(context.Context) error {
	return func(context.Context) error {
		cmd()
		return nil
	}
}

// QualifiedID 通过 FQDN 生成任务唯一 ID。
//
// 只要对象的名称或路径未改变，则该方法生成的 ID 能保证唯一、稳定。
// 由于是通过 FQDN 生成的名字，所以严禁一个 struct 生成多个实例。
// 起名困难症的福音。
func QualifiedID(v any) string {
	s := fqdn(v)
	sum := sha1.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func fqdn(v any) string {
	tof := reflect.TypeOf(v)
	kind := tof.Kind()
	if kind == reflect.Pointer {
		tof = tof.Elem()
	} else if kind == reflect.Func {
		vof := reflect.ValueOf(v).Pointer()
		if pc := runtime.FuncForPC(vof); pc != nil {
			return pc.Name()
		}

		return ""
	}

	pkg, name := tof.PkgPath(), tof.Name()
	if pkg == "" {
		pkg = "main"
	}

	return pkg + "." + name
}
