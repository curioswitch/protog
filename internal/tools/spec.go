package tools

type goos int64

const (
	darwin  goos = 0
	linux   goos = 1
	windows goos = 2
)

type goarch int64

const (
	amd64 goarch = 0
	arm64 goarch = 1
)

type spec struct {
	name         string
	repo         string
	latestVer    func() (string, error)
	os           func(goos goos) string
	arch         func(goarch goarch) string
	ext          func(os string) string
	url          func(ver, os, arch, ext string) string
	postDownload func(dir, os string) error
	path         func(dir, ver, os, arch string) []string
	executables  func(dir, ver, os, arch string) map[string]string
}

type nodeSpec struct {
	name      string
	repo      string
	latestVer func() string
	path      func(dir, ver string) []string
}

type goSpec struct {
	name      string
	repo      string
	latestVer func() string
	cmdPath   string
}
