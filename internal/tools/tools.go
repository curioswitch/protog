package tools

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/curioswitch/protog/internal/proto"
	"github.com/hashicorp/go-getter"
)

type Versions struct {
	Go                      string
	NodeJS                  string
	Protoc                  string
	ProtocGenDoc            string
	ProtocGenDocs           string
	ProtocGenGo             string
	ProtocGenGogoFast       string
	ProtocGenGoGRPC         string
	ProtocGenGolangDeepCopy string
	ProtocGenGolangJSONShim string
	ProtocGenGRPC           string
	ProtocGenGRPCGateway    string
	ProtocGenGRPCJava       string
	ProtocGenJSONSchema     string
	ProtocGenTS             string
	ProtocGenValidate       string
}

type ProtocConfig struct {
	CppGRPC        bool
	CSharpGRPC     bool
	Doc            bool
	Docs           bool
	Go             bool
	GogoFast       bool
	GolangDeepCopy bool
	GolangJSONShim bool
	GoGRPC         bool
	GRPCGateway    bool
	JavaGRPC       bool
	JSONSchema     bool
	NodeGRPC       bool
	ObjectiveCGRPC bool
	PHPGRPC        bool
	PythonGRPC     bool
	RubyGRPC       bool
	TS             bool
	Validate       bool
}

type Config struct {
	Versions Versions
	Protoc   ProtocConfig
}

type ToolManager struct {
	config Config

	dir string

	path        []string
	executables map[string]string
}

func NewToolManager(config Config) (*ToolManager, error) {
	rootDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("could not determine cache dir: %w", err)
	}

	return &ToolManager{
		config: config,

		dir:         filepath.Join(rootDir, "org.curioswitch.protog"),
		executables: map[string]string{},
	}, nil
}

func (m *ToolManager) RunProtoc(args []string, protos []string) error {
	if err := m.fetch(protocSpec, m.config.Versions.Protoc); err != nil {
		return err
	}

	if err := m.fetch(golangSpec, m.config.Versions.Go); err != nil {
		return err
	}
	if err := m.fetch(nodeJSSpec, m.config.Versions.NodeJS); err != nil {
		return err
	}

	if m.config.Protoc.Go {
		if err := m.fetch(protocGenGoSpec, m.config.Versions.ProtocGenGo); err != nil {
			return err
		}
	}

	if m.config.Protoc.GoGRPC {
		if err := m.fetch(protocGenGoGRPCSpec, m.config.Versions.ProtocGenGoGRPC); err != nil {
			return err
		}
	}

	if m.config.Protoc.JavaGRPC {
		if err := m.fetch(protocGenGRPCJavaSpec, m.config.Versions.ProtocGenGRPCJava); err != nil {
			return err
		}
	}

	if m.config.Protoc.PythonGRPC || m.config.Protoc.RubyGRPC {
		if err := m.fetch(protocGenGRPCSpec, m.config.Versions.ProtocGenGRPC); err != nil {
			return err
		}
	}

	// protoc does not source --plugin executables from PATH despite a deceptive error message
	// https://github.com/protocolbuffers/protobuf/issues/10302
	if m.config.Protoc.CppGRPC {
		args = append(args, fmt.Sprintf("--plugin=protoc-gen-grpc_cpp=%s", m.executables["grpc_cpp_plugin"]))
	}
	if m.config.Protoc.CSharpGRPC {
		args = append(args, fmt.Sprintf("--plugin=protoc-gen-grpc_csharp=%s", m.executables["grpc_csharp_plugin"]))
	}
	if m.config.Protoc.ObjectiveCGRPC {
		args = append(args, fmt.Sprintf("--plugin=protoc-gen-grpc_objc=%s", m.executables["grpc_objective_c_plugin"]))
	}
	if m.config.Protoc.NodeGRPC {
		args = append(args, fmt.Sprintf("--plugin=protoc-gen-grpc_node=%s", m.executables["grpc_node_plugin"]))
	}
	if m.config.Protoc.PHPGRPC {
		args = append(args, fmt.Sprintf("--plugin=protoc-gen-grpc_php=%s", m.executables["grpc_php_plugin"]))
	}
	if m.config.Protoc.PythonGRPC {
		args = append(args, fmt.Sprintf("--plugin=protoc-gen-grpc_python=%s", m.executables["grpc_python_plugin"]))
	}
	if m.config.Protoc.RubyGRPC {
		args = append(args, fmt.Sprintf("--plugin=protoc-gen-grpc_ruby=%s", m.executables["grpc_ruby_plugin"]))
	}

	if m.config.Protoc.GogoFast {
		if err := m.fetchGoSpec(protocGenGogoFastSpec, m.config.Versions.ProtocGenGogoFast); err != nil {
			return err
		}
	}

	if m.config.Protoc.Doc {
		if err := m.fetch(protocGenDocSpec, m.config.Versions.ProtocGenDoc); err != nil {
			return err
		}
	}

	if m.config.Protoc.Docs {
		if err := m.fetchGoSpec(protocGenDocsSpec, m.config.Versions.ProtocGenDocs); err != nil {
			return err
		}
	}

	if m.config.Protoc.GRPCGateway {
		if err := m.fetch(protocGenGRPCGatewaySpec, m.config.Versions.ProtocGenGRPCGateway); err != nil {
			return err
		}
	}

	if m.config.Protoc.TS {
		if err := m.fetchNodeSpec(protocGenTSSpec, m.config.Versions.ProtocGenTS); err != nil {
			return err
		}
	}

	if m.config.Protoc.GolangDeepCopy {
		if err := m.fetchGoSpec(protocGenGolangDeepCopySpec, m.config.Versions.ProtocGenGolangDeepCopy); err != nil {
			return err
		}
	}

	if m.config.Protoc.JSONSchema {
		if err := m.fetchGoSpec(protocGenJSONSchemaSpec, m.config.Versions.ProtocGenJSONSchema); err != nil {
			return err
		}
	}

	if m.config.Protoc.GolangJSONShim {
		if err := m.fetchGoSpec(protocGenGolangJSONShimSpec, m.config.Versions.ProtocGenGolangJSONShim); err != nil {
			return err
		}
	}

	if m.config.Protoc.Validate {
		if err := m.fetchGoSpec(protocGenValidateSpec, m.config.Versions.ProtocGenValidate); err != nil {
			return err
		}
	}

	includesDir := filepath.Join("build", "proto-includes")
	if err := os.MkdirAll(includesDir, 0755); err != nil {
		return err
	}
	if err := proto.FetchIncludes(protos, includesDir); err != nil {
		return err
	}
	args = append(args, fmt.Sprintf("--proto_path=%s", includesDir))
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	args = append(args, fmt.Sprintf("--proto_path=%s", cwd))

	cmd := exec.Command(m.executables["protoc"], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = []string{fmt.Sprintf("PATH=%s", mergePath(m.path))}
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func determineLatestVersionForGitHubRepo(repo string) (string, error) {
	latestURL := fmt.Sprintf("https://%s/releases/latest", repo)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(latestURL)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 302 {
		return "", fmt.Errorf("invalid status code: %v", resp.StatusCode)
	}
	location := resp.Header.Get("Location")
	if idx := strings.LastIndexByte(location, '/'); idx != -1 {
		return location[idx+1:], nil
	} else {
		return "", fmt.Errorf("invalid location: %v", location)
	}
}

func (m *ToolManager) fetch(s spec, ver string) error {
	var goos goos
	switch runtime.GOOS {
	case "darwin":
		goos = darwin
	case "linux":
		goos = linux
	case "windows":
		goos = windows
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	var goarch goarch
	switch runtime.GOARCH {
	case "amd64":
		goarch = amd64
	case "arm64":
		goarch = arm64
	default:
		return fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}

	for _, f := range s.goFallbacks {
		if f.arch == goarch {
			return m.fetchGoSpec(f.spec, ver)
		}
	}

	if ver == "" {
		if s.latestVer != nil {
			v, err := s.latestVer()
			if err != nil {
				return err
			}
			ver = v
		} else {
			v, err := determineLatestVersionForGitHubRepo(s.repo)
			if err != nil {
				return err
			}
			ver = v
		}
	} else if ver[0] != 'v' {
		ver = "v" + ver
	}

	var osStr string
	if s.os != nil {
		osStr = s.os(goos)
	} else {
		osStr = runtime.GOOS
	}

	var archStr string
	if s.arch != nil {
		archStr = s.arch(goarch)
	} else {
		archStr = runtime.GOARCH
	}

	var ext string
	if s.ext != nil {
		ext = s.ext(osStr)
	} else {
		switch goos {
		case darwin:
			fallthrough
		case linux:
			ext = "tar.gz"
		case windows:
			ext = "zip"
		}
	}

	dir := filepath.Join(m.dir, s.name, ver)
	if s.path != nil {
		m.path = append(s.path(dir, ver, osStr, archStr), m.path...)
	} else {
		m.path = append([]string{dir}, m.path...)
	}

	if s.executables != nil {
		for k, v := range s.executables(dir, ver, osStr, archStr) {
			m.executables[k] = v
		}
	}

	if _, err := os.Stat(dir); err == nil {
		return nil
	}

	url := s.url(ver, osStr, archStr, ext)

	if err := getter.Get(dir, url, getter.WithUmask(0022), getter.WithMode(getter.ClientModeAny), getter.WithGetters(map[string]getter.Getter{
		"https": &getter.HttpGetter{
			XTerraformGetDisabled: true,
		},
	})); err != nil {
		return fmt.Errorf("fetching %s from %s: %w", s.name, url, err)
	}

	if s.postDownload != nil {
		if err := s.postDownload(dir, osStr); err != nil {
			return err
		}
	}

	return nil
}

func (m *ToolManager) fetchNodeSpec(s nodeSpec, ver string) error {
	if ver == "" {
		if s.latestVer != nil {
			ver = s.latestVer()
		} else {
			v, err := determineLatestVersionForGitHubRepo(s.repo)
			if err != nil {
				return err
			}
			ver = v
		}
	}

	dir := filepath.Join(m.dir, s.name, ver)
	if s.path != nil {
		m.path = append(s.path(dir, ver), m.path...)
	} else {
		m.path = append([]string{dir}, m.path...)
	}

	if _, err := os.Stat(dir); err == nil {
		return nil
	}

	cmd := exec.Command(m.executables["npm"], "install", "--prefix", dir, fmt.Sprintf("%s@%s", s.name, ver))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = []string{fmt.Sprintf("PATH=%s", mergePath(m.path))}
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (m *ToolManager) fetchGoSpec(s goSpec, ver string) error {
	if ver == "" {
		if s.latestVer != nil {
			v, err := s.latestVer()
			if err != nil {
				return err
			}
			ver = v
		} else {
			v, err := determineLatestVersionForGitHubRepo(s.repo)
			if err != nil {
				return err
			}
			ver = v
		}
	}

	dir := filepath.Join(m.dir, s.name, ver)
	m.path = append([]string{filepath.Join(dir, "bin")}, m.path...)

	if _, err := os.Stat(dir); err == nil {
		return nil
	}

	cmd := exec.Command(m.executables["go"], "install", fmt.Sprintf("%s@%s", s.cmdPath, ver))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = []string{fmt.Sprintf("GOPATH=%s", dir), fmt.Sprintf("GOCACHE=%s", filepath.Join(m.dir, "gocache")), "CGO_ENABLED=0"}
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func mergePath(path []string) string {
	sep := ":"
	if runtime.GOOS == "windows" {
		sep = ";"
	}
	return strings.Join(path, sep)
}
