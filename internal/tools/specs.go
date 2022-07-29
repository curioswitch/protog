package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
		case "win":
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
	executables: func(dir, ver, os, arch string) map[string]string {
		return map[string]string{"protoc": filepath.Join(dir, "bin", exe("protoc"))}
	},
}

var protocGenGoSpec = spec{
	name: "protoc-gen-go",
	repo: "github.com/protocolbuffers/protobuf-go",
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
	goFallbacks: []goFallback{
		{
			arch: arm64,
			spec: goSpec{
				name: "protoc-gen-go-grpc",
				repo: "github.com/grpc/grpc-go",
				latestVer: func() (string, error) {
					// TODO: Use REST API? To find latest release as they have non-protoc tags in the same repo.
					return "v1.2.0", nil
				},
				cmdPath: "google.golang.org/grpc/cmd/protoc-gen-go-grpc",
			},
		},
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
		return fmt.Sprintf("https://packages.grpc.io/archive/2022/07/%s/protoc/grpc-protoc_%s_%s-1.49.0-dev.%s", ver[1:], os, arch, ext)
	},
	executables: func(dir, ver, os, arch string) map[string]string {
		return map[string]string{
			"grpc_cpp_plugin":         filepath.Join(dir, exe("grpc_cpp_plugin")),
			"grpc_csharp_plugin":      filepath.Join(dir, exe("grpc_csharp_plugin")),
			"grpc_objective_c_plugin": filepath.Join(dir, exe("grpc_objective_c_plugin")),
			"grpc_node_plugin":        filepath.Join(dir, exe("grpc_node_plugin")),
			"grpc_php_plugin":         filepath.Join(dir, exe("grpc_php_plugin")),
			"grpc_python_plugin":      filepath.Join(dir, exe("grpc_python_plugin")),
			"grpc_ruby_plugin":        filepath.Join(dir, exe("grpc_ruby_plugin")),
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
		filename := exe("protoc-gen-grpc-gateway")

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
		return fmt.Sprintf("https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/%s/protoc-gen-grpc-java-%s-%s-%s.exe?filename=%s", ver[1:], ver[1:], os, arch, exe("protoc-gen-grpc-java"))
	},
	postDownload: func(dir, osStr string) error {
		if err := os.Chmod(filepath.Join(dir, exe("protoc-gen-grpc-java")), 0755); err != nil {
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

var protocGenGRPCWebSpec = spec{
	name: "protoc-gen-grpc-web",
	repo: "github.com/grpc/grpc-web",
	arch: func(goarch goarch) string {
		switch goarch {
		case amd64:
			return "x86_64"
		case arm64:
			return "aarch64"
		default:
			panic(fmt.Sprintf("unsupported arch: %v", goarch))
		}
	},
	url: func(ver, os, arch, ext string) string {
		ver = ver[1:]
		switch arch {
		case "x86_64":
			var suffix string
			if os == "windows" {
				suffix = ".exe"
			}
			return fmt.Sprintf("https://github.com/grpc/grpc-web/releases/download/%s/protoc-gen-grpc-web-%s-%s-%s%s?filename=%s", ver, ver, os, arch, suffix, exe("protoc-gen-grpc-web"))
		case "aarch64":
			// Using self-hosted binaries temporarily since the upstream ones are only in artifacts, requiring auth
			var suffix string
			if os == "windows" {
				suffix = ".exe"
			}
			return fmt.Sprintf("https://github.com/curioswitch/scratch/releases/download/grpc-web-test/protoc-gen-grpc-web-20220613-%s-%s%s?filename=%s", os, arch, suffix, exe("protoc-gen-grpc-web"))
		default:
			panic(fmt.Sprintf("unsupported arch: %v", arch))
		}
	},
	postDownload: func(dir, _ string) error {
		if err := os.Chmod(filepath.Join(dir, exe("protoc-gen-grpc-web")), 0755); err != nil {
			return err
		}

		return nil
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
		nodeDir := filepath.Join(dir, fmt.Sprintf("node-%s-%s-%s", ver, os, arch))
		if os == "win" {
			return []string{nodeDir}
		}
		return []string{filepath.Join(nodeDir, "bin")}
	},
	executables: func(dir, ver, os, arch string) map[string]string {
		nodeDir := filepath.Join(dir, fmt.Sprintf("node-%s-%s-%s", ver, os, arch))
		res := map[string]string{
			"node": filepath.Join(nodeDir, "bin", "node"),
		}
		if os == "win" {
			res["npm"] = filepath.Join(nodeDir, "npm.cmd")
		} else {
			// Workaround symlinks not being preserved by pointing at the lib file directly instead of bin.
			// https://github.com/hashicorp/go-getter/issues/60
			res["npm"] = filepath.Join(nodeDir, "lib", "node_modules", "npm", "bin", "npm-cli.js")
		}
		return res
	},
}

var protocGenTSSpec = nodeSpec{
	name: "protoc-gen-ts",
	repo: "github.com/thesayyn/protoc-gen-ts",
	path: func(dir, ver string) []string {
		return []string{filepath.Join(dir, "node_modules", ".bin")}
	},
}

var improbableTSProtocGenSpec = nodeSpec{
	name: "ts-protoc-gen",
	repo: "github.com/improbable-eng/ts-protoc-gen",
	path: func(dir, ver string) []string {
		// Don't add to path so ts_out corresponds to protoc-gen-ts
		return []string{}
	},
	executables: func(dir string) map[string]string {
		return map[string]string{
			"ts-protoc-gen": filepath.Join(dir, "node_modules", ".bin", cmd("protoc-gen-ts")),
		}
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
		// Strip off leading v
		ver = ver[1:]
		if !strings.HasPrefix(ver, "go") {
			ver = "go" + ver
		}
		return fmt.Sprintf("https://go.dev/dl/%s.%s-%s.%s", ver, os, arch, ext)
	},
	executables: func(dir, ver, os, arch string) map[string]string {
		return map[string]string{
			"go": filepath.Join(dir, "go", "bin", exe("go")),
		}
	},
}

var protocGenValidateSpec = goSpec{
	name:    "protoc-gen-validate",
	repo:    "github.com/envoyproxy/protoc-gen-validate",
	cmdPath: "github.com/envoyproxy/protoc-gen-validate",
}

var protocGenJSONSchemaSpec = goSpec{
	name:       "protoc-gen-jsonschema",
	repo:       "github.com/chrusty/protoc-gen-jsonschema",
	cmdPath:    "github.com/chrusty/protoc-gen-jsonschema/cmd/protoc-gen-jsonschema",
	versionNoV: true,
}

var protocGenDocsSpec = goSpec{
	name:    "protoc-gen-docs",
	repo:    "github.com/istio/tools",
	cmdPath: "istio.io/tools/cmd/protoc-gen-docs",
	latestVer: func() (string, error) {
		// TODO: Fetch tags to get version
		return "1.14.2", nil
	},
	versionNoV: true,
}

var protocGenGolangDeepCopySpec = goSpec{
	name:    "protoc-gen-golang-deepcopy",
	repo:    "github.com/istio/tools",
	cmdPath: "istio.io/tools/cmd/protoc-gen-golang-deepcopy",
	latestVer: func() (string, error) {
		// TODO: Fetch tags to get version
		return "1.14.2", nil
	},
	versionNoV: true,
}

var protocGenGolangJSONShimSpec = goSpec{
	name:    "protoc-gen-golang-jsonshim",
	repo:    "github.com/istio/tools",
	cmdPath: "istio.io/tools/cmd/protoc-gen-golang-jsonshim",
	latestVer: func() (string, error) {
		// TODO: Fetch tags to get version
		return "1.14.2", nil
	},
	versionNoV: true,
}

var protocGenGogoFastSpec = goSpec{
	name:    "protoc-gen-gogofast",
	repo:    "github.com/gogo/protobuf",
	cmdPath: "github.com/gogo/protobuf/protoc-gen-gogofast",
}

func exe(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func cmd(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".cmd"
	}
	return name
}
