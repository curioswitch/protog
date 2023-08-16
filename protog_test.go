package protog

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	cacheDir, _ := os.UserCacheDir()
	expectedFilesNoJS := []string{
		filepath.Join("cpp", "helloworld.pb.h"),
		filepath.Join("cpp", "helloworld.pb.cc"),
		filepath.Join("csharp", "Helloworld.cs"),
		filepath.Join("es", "helloworld_pb.ts"),
		filepath.Join("es", "helloworld_connect.ts"),
		filepath.Join("go", "testing", "helloworld", "helloworld.pb.go"),
		filepath.Join("go", "testing", "helloworld", "helloworldconnect", "helloworld.connect.go"),
		filepath.Join("go", "testing", "helloworld", "helloworld.pb.validate.go"),
		filepath.Join("go", "testing", "helloworld", "helloworld_grpc.pb.go"),
		filepath.Join("go", "testing", "helloworld", "helloworld_deepcopy.gen.go"),
		filepath.Join("go", "testing", "helloworld", "helloworld_json.gen.go"),
		filepath.Join("gogofast", "testing", "helloworld", "helloworld.pb.go"),
		filepath.Join("java", "io", "grpc", "examples", "helloworld", "HelloWorldProto.java"),
		filepath.Join("java", "io", "grpc", "examples", "helloworld", "GreeterGrpc.java"),
		filepath.Join("jsonschema", "HelloReply.json"),
		filepath.Join("jsonschema", "HelloRequest.json"),
		filepath.Join("objc", "Helloworld.pbobjc.h"),
		filepath.Join("objc", "Helloworld.pbobjc.m"),
		filepath.Join("php", "Helloworld", "HelloReply.php"),
		filepath.Join("php", "Helloworld", "HelloRequest.php"),
		filepath.Join("python", "helloworld_pb2.py"),
		filepath.Join("ruby", "helloworld_pb.rb"),
		filepath.Join("ts", "helloworld.ts"),
		filepath.Join("doc", "index.html"),
		filepath.Join("docs", "helloworld.pb.html"),
	}

	if runtime.GOARCH == "amd64" || runtime.GOOS == "darwin" {
		// grpc plugins are only available for amd64, though we do run them on darwin because they work fine with
		// Rosetta.
		expectedFilesNoJS = append(expectedFilesNoJS,
			filepath.Join("cpp", "helloworld.grpc.pb.h"),
			filepath.Join("cpp", "helloworld.grpc.pb.cc"),
			filepath.Join("csharp", "HelloworldGrpc.cs"),
			filepath.Join("php", "Helloworld", "GreeterClient.php"),
			filepath.Join("python", "helloworld_pb2_grpc.py"),
			filepath.Join("ruby", "helloworld_services_pb.rb"),
		)

		if runtime.GOOS != "windows" {
			// grpc_objective_c_plugin crashes on Windows
			expectedFilesNoJS = append(expectedFilesNoJS,
				filepath.Join("objc", "Helloworld.pbrpc.h"),
				filepath.Join("objc", "Helloworld.pbrpc.m"),
			)
		}
	}

	expectedFilesJS := append(expectedFilesNoJS,
		filepath.Join("js", "helloreply.js"),
		filepath.Join("js", "hellorequest.js"),
		filepath.Join("js", "helloworld_grpc_web.js"),
		filepath.Join("js", "helloworld_pb.d.ts"),
		filepath.Join("js", "helloworld_pb.js"),
	)
	if runtime.GOARCH == "amd64" || runtime.GOOS == "darwin" {
		expectedFilesJS = append(expectedFilesJS,
			filepath.Join("js", "helloworld_grpc.js"),
		)
	}

	argsNoJS := func(dir string) []string {
		args := []string{
			"--cpp_out=" + filepath.Join(dir, "cpp"),
			"--csharp_out=" + filepath.Join(dir, "csharp"),
			"--connect-es_out=" + filepath.Join(dir, "es"),
			"--connect-es_opt=target=ts",
			"--es_out=" + filepath.Join(dir, "es"),
			"--es_opt=target=ts",
			"--java_out=" + filepath.Join(dir, "java"),
			"--grpc-java_out=" + filepath.Join(dir, "java"),
			"--go_out=" + filepath.Join(dir, "go"),
			"--go-grpc_out=" + filepath.Join(dir, "go"),
			"--golang-deepcopy_out=" + filepath.Join(dir, "go"),
			"--golang-jsonshim_out=" + filepath.Join(dir, "go"),
			"--grpc-gateway_out=" + filepath.Join(dir, "go"),
			"--gogofast_out=" + filepath.Join(dir, "gogofast"),
			"--php_out=" + filepath.Join(dir, "php"),
			"--python_out=" + filepath.Join(dir, "python"),
			"--objc_out=" + filepath.Join(dir, "objc"),
			"--ruby_out=" + filepath.Join(dir, "ruby"),
			"--ts_out=" + filepath.Join(dir, "ts"),
			"--doc_out=" + filepath.Join(dir, "doc"),
			"--docs_out=" + filepath.Join(dir, "docs"),
			"--connect-go_out=" + filepath.Join(dir, "go"),
			"--validate_out=" + filepath.Join(dir, "go"),
			"--validate_opt=lang=go",
			"--jsonschema_out=" + filepath.Join(dir, "jsonschema"),
			"--proto_path=testdata",
			filepath.Join("testdata", "helloworld.proto"),
		}

		if runtime.GOARCH == "amd64" || runtime.GOOS == "darwin" {
			args = append(args,
				"--grpc_cpp_out="+filepath.Join(dir, "cpp"),
				"--grpc_csharp_out="+filepath.Join(dir, "csharp"),
				"--grpc_php_out="+filepath.Join(dir, "php"),
				"--grpc_python_out="+filepath.Join(dir, "python"),
				"--grpc_ruby_out="+filepath.Join(dir, "ruby"),
			)

			if runtime.GOOS != "windows" {
				args = append(args,
					"--grpc_objc_out="+filepath.Join(dir, "objc"),
				)
			}
		}

		return args
	}

	tests := []struct {
		name          string
		versions      Versions
		args          func(dir string) []string
		expectedFiles []string
		includesDir   string
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
				Go:                      "1.18.3",
				NodeJS:                  "16.16.0",
				Protoc:                  "3.20.1",
				ProtocGenConnectES:      "0.8.5",
				ProtocGenConnectGo:      "1.5.0",
				ProtocGenDoc:            "1.5.1",
				ProtocGenDocs:           "1.13.6",
				ProtocGenES:             "1.1.1",
				ProtocGenGo:             "1.28.1",
				ProtocGenGolangDeepCopy: "1.13.6",
				ProtocGenGolangJSONShim: "1.13.6",
				ProtocGenGogoFast:       "1.3.1",
				ProtocGenGoGRPC:         "1.2.0",
				// Version override not supported yet
				ProtocGenGRPC:        "",
				ProtocGenGRPCGateway: "2.10.3",
				ProtocGenGRPCWeb:     "1.3.1",
				ProtocGenGRPCJava:    "1.47.0",
				ProtocGenJSONSchema:  "1.3.8",
				ProtocGenTS:          "0.8.6",
				ProtocGenValidate:    "0.6.6",
				ProtocTSGen:          "0.14.0",
			},
			// No js_out which isn't in latest protoc or any plugin yet
			args: func(dir string) []string {
				args := append(argsNoJS(dir),
					"--js_out="+filepath.Join(dir, "js"),
					"--js_opt=import_style=commonjs",
					"--improbable_ts_out="+filepath.Join(dir, "js"),
					"--grpc-web_out="+filepath.Join(dir, "js"),
					"--grpc-web_opt=import_style=commonjs,mode=grpcwebtext",
				)
				if runtime.GOARCH == "amd64" || runtime.GOOS == "darwin" {
					args = append(args, "--grpc_js_out="+filepath.Join(dir, "js"))
				}
				return args
			},
			expectedFiles: expectedFilesJS,
			includesDir:   filepath.Join(cacheDir, "org.curioswitch.protog", "proto-includes"),
		},
	}

	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			args := tt.args(dir)
			err := Run(args, Config{
				ProtoIncludesDir: tt.includesDir,
				Versions:         tt.versions,
			})
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
