package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var protocSpec = spec{
	name: "protoc",
	repo: "github.com/protocolbuffers/protobuf",
	os: func(goos goos) string {
		switch goos {
		case darwin:
			return "osx"
		case linux:
			return "linux"
		case windows:
			return "win"
		default:
			panic(fmt.Sprintf("unsupported OS: %v", goos))
		}
	},
	arch: func(goarch goarch) string {
		switch goarch {
		case amd64:
			return "x86_64"
		case arm64:
			return "aarch_64"
		default:
			panic(fmt.Sprintf("unsupported arch: %v", goarch))
		}
	},
	ext: func(string) string {
		return "zip"
	},
	url: func(ver, os, arch, ext string) string {
		// protoc has a uniquely inconsistent URL format.
		var suffix string
		switch os {
		case "osx":
			fallthrough
		case "linux":
			suffix = fmt.Sprintf("%s-%s", os, arch)
		case "windows":
			switch arch {
			case "x86_64":
				suffix = "win64"
			default:
				panic(fmt.Sprintf("unsupported windows arch: %v", arch))
			}
		default:
			panic("BUG")
		}

		return fmt.Sprintf("https://github.com/protocolbuffers/protobuf/releases/download/%s/protoc-%s-%s.zip", ver, ver[1:], suffix)
	},
	postDownload: func(dir string, _ string) error {
		if err := os.Chmod(filepath.Join(dir, "bin", "protoc"), 0755); err != nil {
			return err
		}

		return nil
	},
	executables: func(dir, ver, os, arch string) map[string]string {
		return map[string]string{"protoc": filepath.Join(dir, "bin", "protoc")}
	},
}

var protocGenGoSpec = spec{
	name: "protoc-gen-go",
	repo: "github.com/protocolbuffers/protobuf-go",
	arch: func(goarch) string {
		// Currently only amd64 is published, so just try it and either Rosetta or qemu may work.
		// https://github.com/golang/protobuf/issues/1466
		return "amd64"
	},
	url: func(ver, os, arch, ext string) string {
		return fmt.Sprintf("https://github.com/protocolbuffers/protobuf-go/releases/download/%s/protoc-gen-go.%s.%s.%s.%s", ver, ver, os, arch, ext)
	},
}

var protocGenGoGRPCSpec = spec{
	name: "protoc-gen-go-grpc",
	repo: "github.com/grpc/grpc-go",
	latestVer: func() (string, error) {
		// TODO: Use REST API? To find latest release as they have non-protoc tags in the same repo.
		return "v1.2.0", nil
	},
	arch: func(goarch) string {
		// Currently only amd64 is published, so just try it and either Rosetta or qemu may work.
		// https://github.com/golang/protobuf/issues/1466
		return "amd64"
	},
	ext: func(os string) string {
		return "tar.gz"
	},
	url: func(ver, os, arch, ext string) string {
		return fmt.Sprintf("https://github.com/grpc/grpc-go/releases/download/cmd%%2Fprotoc-gen-go-grpc%%2F%s/protoc-gen-go-grpc.%s.%s.%s.tar.gz", ver, ver, os, arch)
	},
}

var protocGenGRPCSpec = spec{
	name: "protoc-gen-grpc",
	repo: "github.com/grpc/grpc",
	latestVer: func() (string, error) {
		// TODO: Release binaries are not published in a consumable way yet
		// https://github.com/grpc/grpc/issues/30369
		return "0aba64fa077ee9e9c2762b883e1c8935c2d0b0a4-d4c05fc7-3960-48e0-aec0-8dfaa8f5016a", nil
	},
	os: func(goos goos) string {
		switch goos {
		case darwin:
			return "macos"
		case linux:
			return "linux"
		case windows:
			return "windows"
		default:
			panic(fmt.Sprintf("unsupported OS: %v", goos))
		}
	},
	arch: func(goarch) string {
		// Currently only amd64 is published, so just try it and either Rosetta or qemu may work.
		return "x64"
	},
	url: func(ver, os, arch, ext string) string {
		return fmt.Sprintf("https://packages.grpc.io/archive/2022/07/%s/protoc/grpc-protoc_%s_%s-1.49.0-dev.%s", ver, os, arch, ext)
	},
	postDownload: func(dir, _ string) error {
		// Remove protoc which is versioned separately from what we want.
		if err := os.Remove(filepath.Join(dir, "protoc")); err != nil {
			return err
		}

		return nil
	},
	executables: func(dir, ver, os, arch string) map[string]string {
		return map[string]string{
			"grpc_cpp_plugin":         filepath.Join(dir, "grpc_cpp_plugin"),
			"grpc_csharp_plugin":      filepath.Join(dir, "grpc_csharp_plugin"),
			"grpc_objective_c_plugin": filepath.Join(dir, "grpc_objective_c_plugin"),
			"grpc_node_plugin":        filepath.Join(dir, "grpc_node_plugin"),
			"grpc_php_plugin":         filepath.Join(dir, "grpc_php_plugin"),
			"grpc_python_plugin":      filepath.Join(dir, "grpc_python_plugin"),
			"grpc_ruby_plugin":        filepath.Join(dir, "grpc_ruby_plugin"),
		}
	},
}

var protocGenGRPCGatewaySpec = spec{
	name: "protoc-gen-grpc-gateway",
	repo: "github.com/grpc-ecosystem/grpc-gateway",
	arch: func(goarch goarch) string {
		switch goarch {
		case amd64:
			return "x86_64"
		case arm64:
			return "arm64"
		default:
			panic(fmt.Sprintf("unsupported arch: %v", goarch))
		}
	},
	url: func(ver, os, arch, ext string) string {
		filename := "protoc-gen-grpc-gateway"
		var suffix string
		if os == "windows" {
			suffix = ".exe"
			filename += suffix
		}

		return fmt.Sprintf("https://github.com/grpc-ecosystem/grpc-gateway/releases/download/%s/protoc-gen-grpc-gateway-%s-%s-%s%s?filename=%s", ver, ver, os, arch, suffix, filename)
	},
	postDownload: func(dir, osStr string) error {
		filename := "protoc-gen-grpc-gateway"
		if osStr == "windows" {
			filename += ".exe"
		}

		if err := os.Chmod(filepath.Join(dir, filename), 0755); err != nil {
			return err
		}

		return nil
	},
}

var protocGenGRPCJavaSpec = spec{
	name: "protoc-gen-grpc-java",
	repo: "github.com/grpc/grpc-java",
	os: func(goos goos) string {
		switch goos {
		case darwin:
			return "osx"
		case linux:
			return "linux"
		case windows:
			return "windows"
		default:
			panic(fmt.Sprintf("unsupported OS: %v", goos))
		}
	},
	arch: func(goarch goarch) string {
		switch goarch {
		case amd64:
			return "x86_64"
		case arm64:
			return "aarch_64"
		default:
			panic(fmt.Sprintf("unsupported arch: %v", goarch))
		}
	},
	url: func(ver, os, arch, ext string) string {
		filename := "protoc-gen-grpc-java"
		if os == "windows" {
			filename += ".exe"
		}

		return fmt.Sprintf("https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/%s/protoc-gen-grpc-java-%s-%s-%s.exe?filename=%s", ver[1:], ver[1:], os, arch, filename)
	},
	postDownload: func(dir, osStr string) error {
		filename := "protoc-gen-grpc-java"
		if osStr == "windows" {
			filename += ".exe"
		}

		if err := os.Chmod(filepath.Join(dir, filename), 0755); err != nil {
			return err
		}

		return nil
	},
}

var protocGenDocSpec = spec{
	name: "protoc-gen-doc",
	repo: "github.com/pseudomuto/protoc-gen-doc",
	ext: func(string) string {
		return "tar.gz"
	},
	url: func(ver, os, arch, ext string) string {
		return fmt.Sprintf("https://github.com/pseudomuto/protoc-gen-doc/releases/download/%s/protoc-gen-doc_%s_%s_%s.%s", ver, ver[1:], os, arch, ext)
	},
}

var nodeJSSpec = spec{
	name: "nodejs",
	repo: "github.com/nodejs/node",
	os: func(goos goos) string {
		switch goos {
		case darwin:
			return "darwin"
		case linux:
			return "linux"
		case windows:
			return "win"
		default:
			panic(fmt.Sprintf("unsupported OS: %v", goos))
		}
	},
	arch: func(goarch goarch) string {
		switch goarch {
		case amd64:
			return "x64"
		case arm64:
			return "arm64"
		default:
			panic(fmt.Sprintf("unsupported arch: %v", goarch))
		}
	},
	url: func(ver, os, arch, ext string) string {
		return fmt.Sprintf("https://nodejs.org/dist/%s/node-%s-%s-%s.%s", ver, ver, os, arch, ext)
	},
	path: func(dir, ver, os, arch string) []string {
		return []string{filepath.Join(dir, fmt.Sprintf("node-%s-%s-%s", ver, os, arch), "bin")}
	},
	executables: func(dir, ver, os, arch string) map[string]string {
		bin := filepath.Join(dir, fmt.Sprintf("node-%s-%s-%s", ver, os, arch), "bin")
		return map[string]string{
			"node": filepath.Join(bin, "node"),
			"npm":  filepath.Join(bin, "npm"),
		}
	},
}

var protocGenTSSpec = nodeSpec{
	name: "protoc-gen-ts",
	repo: "github.com/thesayyn/protoc-gen-ts",
	path: func(dir, ver string) []string {
		return []string{filepath.Join(dir, "node_modules", ".bin")}
	},
}

var golangSpec = spec{
	name: "golang",
	repo: "github.com/golang/go",
	latestVer: func() (string, error) {
		resp, err := http.Get("https://go.dev/VERSION?m=text")
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		ver, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(ver), nil
	},
	url: func(ver, os, arch, ext string) string {
		return fmt.Sprintf("https://go.dev/dl/%s.%s-%s.%s", ver, os, arch, ext)
	},
	executables: func(dir, ver, os, arch string) map[string]string {
		return map[string]string{
			"go": filepath.Join(dir, "go", "bin", "go"),
		}
	},
}

var protocGenValidateSpec = goSpec{
	name:    "protoc-gen-validate",
	repo:    "github.com/envoyproxy/protoc-gen-validate",
	cmdPath: "github.com/envoyproxy/protoc-gen-validate",
}
