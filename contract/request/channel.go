package request

type ChannelOpen struct {
	ID     string `json:"id"     validate:"required"`
	Goos   string `json:"goos"   validate:"required,oneof=darwin dragonfly illumos ios js wasip1 linux android solaris freebsd nacl netbsd openbsd plan9 windows aix"`
	Arch   string `json:"arch"   validate:"required,oneof=386 amd64 arm arm64 loong64 mips mipsle mips64 mips64le ppc64 ppc64le riscv64 s390x sparc64 wasm"`
	Secret string `json:"secret" validate:"required,lte=500"`
}
