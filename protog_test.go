package protog

import (
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	expectedFilesNoJS := []string{
		filepath.Join("cpp", "helloworld.pb.h"),
		filepath.Join("cpp", "helloworld.pb.cc"),
		filepath.Join("cpp", "helloworld.grpc.pb.h"),
		filepath.Join("cpp", "helloworld.grpc.pb.cc"),
		filepath.Join("csharp", "Helloworld.cs"),
		filepath.Join("csharp", "HelloworldGrpc.cs"),
		filepath.Join("go", "testing", "helloworld", "helloworld.pb.go"),
		filepath.Join("go", "testing", "helloworld", "helloworld.pb.validate.go"),
		filepath.Join("go", "testing", "helloworld", "helloworld_grpc.pb.go"),
		filepath.Join("java", "io", "grpc", "examples", "helloworld", "HelloWorldProto.java"),
		filepath.Join("java", "io", "grpc", "examples", "helloworld", "GreeterGrpc.java"),
		filepath.Join("objc", "Helloworld.pbobjc.h"),
		filepath.Join("objc", "Helloworld.pbobjc.m"),
		filepath.Join("php", "Helloworld", "GreeterClient.php"),
		filepath.Join("php", "Helloworld", "HelloReply.php"),
		filepath.Join("php", "Helloworld", "HelloRequest.php"),
		filepath.Join("python", "helloworld_pb2.py"),
		filepath.Join("python", "helloworld_pb2_grpc.py"),
		filepath.Join("ruby", "helloworld_pb.rb"),
		filepath.Join("ruby", "helloworld_services_pb.rb"),
		filepath.Join("ts", "helloworld.ts"),
		filepath.Join("doc", "index.html"),
	}

	if runtime.GOOS != "windows" {
		// grpc_objective_c_plugin crashes on Windows
		expectedFilesNoJS = append(expectedFilesNoJS,
			filepath.Join("objc", "Helloworld.pbrpc.h"),
			filepath.Join("objc", "Helloworld.pbrpc.m"),
		)
	}

	argsNoJS := func(dir string) []string {
		args := []string{
			"--cpp_out=" + filepath.Join(dir, "cpp"),
			"--grpc_cpp_out=" + filepath.Join(dir, "cpp"),
			"--csharp_out=" + filepath.Join(dir, "csharp"),
			"--grpc_csharp_out=" + filepath.Join(dir, "csharp"),
			"--java_out=" + filepath.Join(dir, "java"),
			"--grpc-java_out=" + filepath.Join(dir, "java"),
			"--go_out=" + filepath.Join(dir, "go"),
			"--go-grpc_out=" + filepath.Join(dir, "go"),
			"--grpc-gateway_out=" + filepath.Join(dir, "go"),
			"--php_out=" + filepath.Join(dir, "php"),
			"--grpc_php_out=" + filepath.Join(dir, "php"),
			"--python_out=" + filepath.Join(dir, "python"),
			"--grpc_python_out=" + filepath.Join(dir, "python"),
			"--objc_out=" + filepath.Join(dir, "objc"),
			"--ruby_out=" + filepath.Join(dir, "ruby"),
			"--grpc_ruby_out=" + filepath.Join(dir, "ruby"),
			"--ts_out=" + filepath.Join(dir, "ts"),
			"--doc_out=" + filepath.Join(dir, "doc"),
			"--validate_out=" + filepath.Join(dir, "go"),
			"--validate_opt=lang=go",
			"--proto_path=testdata",
			filepath.Join("testdata", "helloworld.proto"),
		}

		if runtime.GOOS != "windows" {
			args = append(args,
				"--grpc_objc_out="+filepath.Join(dir, "objc"),
			)
		}

		return args
	}

	tests := []struct {
		name          string
		versions      Versions
		args          func(dir string) []string
		expectedFiles []string
	}{
		{
			name:     "latest versions",
			versions: Versions{},
			// No js_out which isn't in latest protoc or any plugin yet
			args: func(dir string) []string {
				return argsNoJS(dir)
			},
			expectedFiles: expectedFilesNoJS,
		},
		{
			name: "fixed versions",
			versions: Versions{
				Go:           "1.18.3",
				NodeJS:       "16.16.0",
				Protoc:       "3.20.1",
				ProtocGenDoc: "1.5.1",
				ProtocGenGo:  "1.2.0",
				// Version override not supported yet
				ProtocGenGRPC:        "",
				ProtocGenGRPCGateway: "2.10.3",
				ProtocGenGRPCJava:    "1.47.0",
				ProtocGenTS:          "0.8.4",
				ProtocGenValidate:    "0.6.6",
			},
			// No js_out which isn't in latest protoc or any plugin yet
			args: func(dir string) []string {
				return append([]string{"--js_out=" + filepath.Join(dir, "js"), "--grpc_node_out=" + filepath.Join(dir, "js")}, argsNoJS(dir)...)
			},
			expectedFiles: append([]string{
				filepath.Join("js", "helloreply.js"),
				filepath.Join("js", "hellorequest.js"),
				filepath.Join("js", "helloworld_grpc_pb.js"),
			}, expectedFilesNoJS...),
		},
	}

	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			args := tt.args(dir)
			err := Run(args, tt.versions)
			var files []string
			_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				files = append(files, path)
				return nil
			})
			require.NoErrorf(t, err, "args: %v\nfiles\n%s", args, strings.Join(files, "\n"))

			for _, f := range expectedFilesNoJS {
				require.FileExistsf(t, filepath.Join(dir, f), "found files:\n%s", strings.Join(files, "\n"))
			}
		})
	}
}
