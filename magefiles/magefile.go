package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

func BuildAll() error {
	return sh.Run("go", "run", fmt.Sprintf("github.com/goreleaser/goreleaser@%s", goReleaserVer), "build")
}

func Build() error {
	return sh.Run("go", "build", "-o", "build/protog", "./cmd")
}

func Test() error {
	return sh.Run("go", "test", "./...")
}

func Coverage() error {
	if err := sh.Run("go", "test", "-race", "-coverprofile=build/coverage.txt", "-covermode=atomic", "-coverpkg=github.com/curioswitch/protog/...", "./..."); err != nil {
		return err
	}

	if err := sh.Run("go", "tool", "cover", "-html=build/coverage.txt", "-o", "build/coverage.html"); err != nil {
		return err
	}

	return nil
}

func Format() error {
	return sh.Run("go", "run", fmt.Sprintf("github.com/rinchsan/gosimports/cmd/gosimports@%s", gosImportsVer), "-w", ".")
}

var Default = Build
