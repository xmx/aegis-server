package launch

import (
	"context"
	"os"

	"github.com/xmx/aegis-server/application/config"
)

// Run 读取配置文件并启动应用，如果配置文件不存在则启动初始化向导。
func Run(ctx context.Context, path string) error {
	cfg, err := config.JSONC(path)
	if err != nil && os.IsNotExist(err) {
		cfg, err = runOOBE(ctx, path)
	}
	if err != nil {
		return err
	}

	return runApp(ctx, cfg)
}
