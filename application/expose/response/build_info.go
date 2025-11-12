package response

type ExecutableMetadata struct {
	Goos    string `json:"goos"`
	Goarch  string `json:"goarch"`
	Version string `json:"version"`
}
