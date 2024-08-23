package banner

import (
	"os"
	"os/user"
	"runtime/debug"
	"sync"
	"time"
)

const ansiLogo = "\u001B[1;33m" +
	"   ______________  ___ \n" +
	"  / __/  _/ __/  |/  / \n" +
	" _\\ \\_/ // _// /|// /  \n" +
	"/___/___/___/_/  /_/   \u001B[0m\u001B[1;31m%s\u001B[0m\n" +
	"Powered By: https://github.com/xmx\n\n" +
	"    操作系统: \u001B[1;1m%s\u001B[0m\n" +
	"    系统架构: \u001B[1;1m%s\u001B[0m\n" +
	"    主机名称: \u001B[1;1m%s\u001B[0m\n" +
	"    当前用户: \u001B[1;1m%s\u001B[0m\n" +
	"    工作目录: \u001B[1;1m%s\u001B[0m\n" +
	"    进程 PID: \u001B[1;1m%d\u001B[0m\n" +
	"    编译时间: \u001B[1;1m%s\u001B[0m\n" +
	"    提交时间: \u001B[1;1m%s\u001B[0m\n" +
	"    修订版本: \u001B[1;1mhttps://%s/tree/%s\u001B[0m\n\n\n"

var (
	// version 项目发布版本号
	// 项目每次发布版本后会打一个 tag, 这个版本号就来自 git 最新的 tag
	version string

	// revision 修订版本, 代码最近一次的提交 ID
	revision string

	// compileAt 编译时间, 由编译脚本在编译时 -X 写入。
	compileTime string

	compileAt time.Time

	// commitAt 代码最近一次提交时间
	commitAt time.Time

	path string

	// pid 进程 ID
	pid int

	// username 当前系统用户名
	username string

	workdir string

	// hostname 主机名
	hostname string

	onceParseArgs = sync.OnceFunc(parseArgs)
)

// parseArgs 处理编译与运行时参数
func parseArgs() {
	pid = os.Getpid() // 获取 PID
	if cu, _ := user.Current(); cu != nil {
		username = cu.Username
	}
	hostname, _ = os.Hostname()
	workdir, _ = os.Getwd()
	compileAt = parseLocalTime(compileTime)

	info, _ := debug.ReadBuildInfo()
	if info == nil {
		return
	}

	path = info.Main.Path
	if version == "" {
		version = info.Main.Version
	}

	settings := info.Settings
	for _, set := range settings {
		key, val := set.Key, set.Value
		switch key {
		case "vcs.revision":
			revision = val
		case "vcs.time":
			commitAt = parseLocalTime(val)
		}
	}
}

// parseLocalTime 给定一个字符串格式化为当前地区的时间。
//
// - time.RFC1123Z Linux `date -R` 输出的时间格式。
// - time.UnixDate macOS `date` 输出的时间格式。
func parseLocalTime(str string) time.Time {
	for _, layout := range []string{
		time.RFC1123Z, time.UnixDate, time.Layout, time.ANSIC,
		time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850,
		time.RFC1123, time.RFC3339, time.RFC3339Nano, time.Kitchen,
		time.Stamp, time.StampMilli, time.StampMicro, time.StampNano,
		time.DateTime, time.DateOnly,
	} {
		if at, err := time.Parse(layout, str); err == nil {
			return at.Local()
		}
	}
	epoch := time.Date(0, 1, 1, 0, 0, 0, 0, time.Local)

	return epoch
}
