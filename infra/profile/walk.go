package profile

import (
	"io"
	"os"

	"github.com/xmx/aegis-server/library/jsonc"
)

func JSONC(path string) (*Config, error) {
	cfg := new(Config)
	for name, err := range readdir(path, "*.jsonc") {
		if err != nil {
			return nil, err
		}
		if err = unmarshalJSONC(name, cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func unmarshalJSONC(name string, v any) error {
	fd, err := os.Open(name)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer fd.Close()

	// 我们使用的 jsonc 不支持流式解析，为了防止恶意大文件造成的 OOM，此处
	// 限制读取文件 8MiB，按照常理和经验，该大小已经足够容纳正常的 jsonc 配置了。
	lr := io.LimitReader(fd, 2<<22) // 8MiB
	data, err := io.ReadAll(lr)
	if err != nil {
		return err
	}

	return jsonc.Unmarshal(data, v)
}
