package serverd

type AuthConfig struct {
	URI string `json:"uri"` // 数据库连接配置。
}

type authResponse struct {
	Code    int         `json:"code"` // 2xx: succeed
	Message string      `json:"message"`
	Config  *AuthConfig `json:"config,omitzero"`
}

type AuthRequest struct {
	Secret     string   `json:"secret"              validate:"required,gte=10,lte=100"`
	Inet       string   `json:"inet"                validate:"required,ip"`
	Goos       string   `json:"goos"                validate:"required,oneof=darwin dragonfly illumos ios js wasip1 linux android solaris freebsd nacl netbsd openbsd plan9 windows aix"`
	Goarch     string   `json:"goarch"              validate:"required,oneof=386 amd64 arm arm64 loong64 mips mipsle mips64 mips64le ppc64 ppc64le riscv64 s390x sparc64 wasm"`
	Semver     string   `json:"semver"              validate:"required,semver"`
	PID        int      `json:"pid,omitzero"`
	Args       []string `json:"args,omitzero"`
	Hostname   string   `json:"hostname,omitzero"`
	Workdir    string   `json:"workdir,omitzero"`
	Executable string   `json:"executable,omitzero"`
	Username   string   `json:"username,omitzero"`
}
