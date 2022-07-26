package proto

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/go-getter"
)

type includeSpec struct {
	prefix  string
	repo    string
	repoDir string
	dir     string
}

var importRe = regexp.MustCompile(`.*import "([^"]+)";.*`)

func FetchIncludes(protos []string, dir string) error {
	for _, proto := range protos {
		if err := fetchInclude(proto, dir); err != nil {
			return err
		}
	}

	return nil
}

func fetchInclude(proto string, dir string) error {
	f, err := os.Open(proto)
	if err != nil {
		return err
	}
	defer f.Close()
	// It would be simpler to use a structured parse, but protoc does not seem to allow it with missing imports.
	// This regex should work well enough.
	// https://github.com/protocolbuffers/protobuf/issues/10310
	s := bufio.NewScanner(f)
	for s.Scan() {
		if m := importRe.FindStringSubmatch(s.Text()); len(m) > 0 {
			for _, includeSpec := range includeSpecs {
				if strings.HasPrefix(m[1], includeSpec.prefix) {
					dst := filepath.Join(dir, includeSpec.dir)

					if _, err := os.Stat(dst); err == nil {
						continue
					}

					url := fmt.Sprintf("%s//%s?depth=1", includeSpec.repo, includeSpec.repoDir)
					if err := getter.Get(dst, url, getter.WithUmask(0022), getter.WithMode(getter.ClientModeAny)); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
