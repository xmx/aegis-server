package config

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

const (
	// EnvironOOBEAddr 初始化向导服务监听地址的环境变量。
	EnvironOOBEAddr = "SIEM_OOBE_ADDR"

	// EnvironOOBEDist 自定义初始化 webui 文件路径。
	EnvironOOBEDist = "SIEM_OOBE_DIST"

	// DefaultLogFilename 默认日志输出位置。
	DefaultLogFilename = "resources/log/application.jsonl"
)

type Config struct {
	Database Database `json:"database"`
}

type Database struct {
	URI string `json:"uri" validate:"required"`
}

// JSONC 方式读取配置文件，由于这种方式不是流式读取，为了防止误读大文件
// 导致 OOM，可以选择一个合适的值限制最大读取量。
func JSONC(filename string, maxsize ...int64) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var size int64
	if len(maxsize) > 0 {
		size = maxsize[0]
	}
	lr := wrapLimitReader(f, size)

	raw, err := io.ReadAll(lr)
	if err != nil {
		return nil, err
	}

	bs := toJSON(raw, nil)
	cfg := new(Config)
	dec := json.NewDecoder(bytes.NewReader(bs))
	dec.DisallowUnknownFields() // 严格模式
	if err = dec.Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// WriteFile 将配置写入到文件，支持在第一行添加注释。
func WriteFile(name string, cfg *Config, comment string) error {
	dir := filepath.Dir(name)
	if _, err := os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err = os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	buf := new(bytes.Buffer)
	if comment != "" {
		// 注释中不能包含换行等特殊符号，如有请调用者自行处理。
		buf.WriteString("// " + comment + "\n")
	}

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil {
		return err
	}

	return os.WriteFile(name, buf.Bytes(), 0o600)
}

func wrapLimitReader(r io.Reader, size int64) io.Reader {
	if size <= 0 {
		return r
	}

	return io.LimitReader(r, size)
}
