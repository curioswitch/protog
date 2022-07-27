package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/sh"
)

func Snapshot() error {
	return sh.RunV("go", "run", fmt.Sprintf("github.com/goreleaser/goreleaser@%s", goReleaserVer), "release", "--snapshot", "--rm-dist")
}

func Release() error {
	return sh.RunV("go", "run", fmt.Sprintf("github.com/goreleaser/goreleaser@%s", goReleaserVer), "release", "--rm-dist")
}

func Build() error {
	return sh.RunV("go", "build", "-o", "build/protog", "./cmd")
}

func Test() error {
	return sh.RunV("go", "test", "./...")
}

func Coverage() error {
	if err := sh.RunV("go", "test", "-race", "-coverprofile=build/coverage.txt", "-covermode=atomic", "-coverpkg=github.com/curioswitch/protog/...", "./..."); err != nil {
		return err
	}

	if err := sh.RunV("go", "tool", "cover", "-html=build/coverage.txt", "-o", "build/coverage.html"); err != nil {
		return err
	}

	return nil
}

func Format() error {
	return sh.RunV("go", "run", fmt.Sprintf("github.com/rinchsan/gosimports/cmd/gosimports@%s", gosImportsVer), "-w", ".")
}

func CacheDir() error {
	dir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	fmt.Print(filepath.Join(dir, "org.curioswitch.protog"))
	return nil
}

var Default = Build
