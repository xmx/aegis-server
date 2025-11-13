package request

type ReleaseDownload struct {
	Goos   string `json:"goos"   form:"goos"   query:"goos"   validate:"required,oneof=windows linux darwin"`
	Goarch string `json:"goarch" form:"goarch" query:"goarch" validate:"required,oneof=amd64 arm64"`
}

type ReleaseBrokerDownload struct {
	Name string `json:"name" form:"name" query:"name" validate:"required"`
	ReleaseDownload
}
