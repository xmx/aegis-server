package request

type ReleaseDownload struct {
	Goos   string `json:"goos"   form:"goos"   query:"goos"   validate:"required,oneof=darwin dragonfly illumos ios js wasip1 linux android solaris freebsd nacl netbsd openbsd plan9 windows aix"`
	Goarch string `json:"goarch" form:"goarch" query:"goarch" validate:"required,oneof=386 amd64 arm arm64 loong64 mips mipsle mips64 mips64le ppc64 ppc64le riscv64 s390x sparc64 wasm"`
}

type ReleaseBrokerDownload struct {
	Name string `json:"name" form:"name" query:"name" validate:"required"`
	ReleaseDownload
}
