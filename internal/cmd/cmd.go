package cmd

import (
	"os"

	"github.com/curioswitch/protog/internal/tools"
	"github.com/spf13/cobra"
)

func Run(args []string, env map[string]string) error {
	var cppOut string
	var cppGRPCOut string
	var cSharpOut string
	var cSharpGRPCOut string
	var docOut string
	var docsOut string
	var goOut string
	var gogoFastOut string
	var golangDeepCopyOut string
	var golangJsonShimOut string
	var goGrpcOut string
	var grpcGatewayOut string
	var grpcWebOut string
	var javaOut string
	var javaGrpcOut string
	var jsOut string
	var jsonSchemaOut string
	var nodeGRPCOut string
	var objcOut string
	var objcGRPCOut string
	var phpOut string
	var phpGRPCOut string
	var pythonOut string
	var pythonGRPCOut string
	var rubyOut string
	var rubyGRPCOut string
	var tsOut string
	var improbableTSOut string
	var validateOut string

	cmd := &cobra.Command{
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		RunE: func(c *cobra.Command, protos []string) error {
			for _, path := range []string{
				cppOut,
				cppGRPCOut,
				cSharpOut,
				cSharpGRPCOut,
				docOut,
				docsOut,
				goOut,
				gogoFastOut,
				golangDeepCopyOut,
				golangJsonShimOut,
				goGrpcOut,
				grpcGatewayOut,
				grpcWebOut,
				improbableTSOut,
				javaOut,
				javaGrpcOut,
				jsOut,
				jsonSchemaOut,
				nodeGRPCOut,
				objcOut,
				objcGRPCOut,
				phpOut,
				phpGRPCOut,
				pythonOut,
				pythonGRPCOut,
				rubyOut,
				rubyGRPCOut,
				tsOut,
				validateOut,
			} {
				if err := mkdir(path); err != nil {
					return err
				}
			}

			m, err := tools.NewToolManager(
				tools.Config{
					Versions: tools.Versions{
						Go:                      env["GO_VERSION"],
						NodeJS:                  env["NODEJS_VERSION"],
						Protoc:                  env["PROTOC_VERSION"],
						ProtocGenDoc:            env["PROTOC_GEN_DOC_VERSION"],
						ProtocGenDocs:           env["PROTOC_GEN_DOCS_VERSION"],
						ProtocGenGo:             env["PROTOC_GEN_GO_VERSION"],
						ProtocGenGolangDeepCopy: env["PROTOC_GEN_GOLANG_DEEPCOPY_VERSION"],
						ProtocGenGolangJSONShim: env["PROTOC_GEN_GOLANG_JSONSHIM_VERSION"],
						ProtocGenGogoFast:       env["PROTOC_GEN_GOGO_FAST_VERSION"],
						ProtocGenGoGRPC:         env["PROTOC_GEN_GO_GRPC_VERSION"],
						ProtocGenGRPC:           env["PROTOC_GEN_GRPC_VERSION"],
						ProtocGenGRPCGateway:    env["PROTOC_GEN_GRPC_GATEWAY_VERSION"],
						ProtocGenGRPCWeb:        env["PROTOC_GEN_GRPC_WEB_VERSION"],
						ProtocGenGRPCJava:       env["PROTOC_GEN_GRPC_JAVA_VERSION"],
						ProtocGenJSONSchema:     env["PROTOC_GEN_JSONSCHEMA_VERSION"],
						ProtocGenTS:             env["PROTOC_GEN_TS_VERSION"],
						ProtocGenValidate:       env["PROTOC_GEN_VALIDATE_VERSION"],
						ProtocTSGen:             env["PROTOC_TS_GEN_VERSION"],
					},
					Protoc: tools.ProtocConfig{
						CppGRPC:        cppGRPCOut != "",
						CSharpGRPC:     cSharpGRPCOut != "",
						Doc:            docOut != "",
						Docs:           docsOut != "",
						Go:             goOut != "",
						GogoFast:       gogoFastOut != "",
						GolangDeepCopy: golangDeepCopyOut != "",
						GolangJSONShim: golangJsonShimOut != "",
						GoGRPC:         goGrpcOut != "",
						GRPCGateway:    grpcGatewayOut != "",
						GRPCWeb:        grpcWebOut != "",
						ImprobableTS:   improbableTSOut != "",
						JavaGRPC:       javaGrpcOut != "",
						JSONSchema:     jsonSchemaOut != "",
						NodeGRPC:       nodeGRPCOut != "",
						ObjectiveCGRPC: objcGRPCOut != "",
						PHPGRPC:        phpGRPCOut != "",
						PythonGRPC:     pythonGRPCOut != "",
						RubyGRPC:       rubyGRPCOut != "",
						TS:             tsOut != "",
						Validate:       validateOut != "",
					},
				},
			)
			if err != nil {
				return err
			}

			if err := m.RunProtoc(args, protos); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.SetArgs(args)

	cmd.Flags().StringVar(&cppOut, "cpp_out", "", "Generate C++ header and source.")
	cmd.Flags().StringVar(&cppGRPCOut, "grpc_cpp_out", "", "Generate C++ gRPC header and source.")

	cmd.Flags().StringVar(&cSharpOut, "csharp_out", "", "Generate C# source file.")
	cmd.Flags().StringVar(&cSharpGRPCOut, "grpc_csharp_out", "", "Generate C# gRPC source file.")

	cmd.Flags().StringVar(&jsOut, "js_out", "", "Generate JS source file.")
	cmd.Flags().StringVar(&nodeGRPCOut, "grpc_node_out", "", "Generate NodeJS gRPC source file.")

	cmd.Flags().StringVar(&javaOut, "java_out", "", "Generate Java source file.")
	cmd.Flags().StringVar(&javaGrpcOut, "grpc-java_out", "", "Generate Java gRPC source file.")

	cmd.Flags().StringVar(&goOut, "go_out", "", "Generate Go source file.")
	cmd.Flags().StringVar(&goGrpcOut, "go-grpc_out", "", "Generate Go gRPC source file.")

	cmd.Flags().StringVar(&grpcGatewayOut, "grpc-gateway_out", "", "Generate Go gRPC gateway source file.")

	cmd.Flags().StringVar(&grpcWebOut, "grpc-web_out", "", "Generate gRPC Web source file.")

	cmd.Flags().StringVar(&objcOut, "objc_out", "", "Generate Objective-C header and source.")
	cmd.Flags().StringVar(&objcGRPCOut, "grpc_objc_out", "", "Generate Objective-C header and source.")

	cmd.Flags().StringVar(&phpOut, "php_out", "", "Generate PHP source file.")
	cmd.Flags().StringVar(&phpGRPCOut, "grpc_php_out", "", "Generate PHP gRPC source file.")

	cmd.Flags().StringVar(&pythonOut, "python_out", "", "Generate Python source file.")
	cmd.Flags().StringVar(&pythonGRPCOut, "grpc_python_out", "", "Generate Python gRPC source file.")

	cmd.Flags().StringVar(&rubyOut, "ruby_out", "", "Generate Ruby source file.")
	cmd.Flags().StringVar(&rubyGRPCOut, "grpc_ruby_out", "", "Generate Ruby gRPC source file.")

	cmd.Flags().StringVar(&docOut, "doc_out", "", "Generate docs.")
	cmd.Flags().StringVar(&docsOut, "docs_out", "", "Generate docs (istio tools).")

	cmd.Flags().StringVar(&tsOut, "ts_out", "", "Generate TypeScript source file.")
	cmd.Flags().StringVar(&improbableTSOut, "improbable_ts_out", "", "Generate TypeScript source file using improbable-eng compiler.")

	cmd.Flags().StringVar(&jsonSchemaOut, "jsonschema_out", "", "Generate JSON Schema file.")

	cmd.Flags().StringVar(&validateOut, "validate_out", "", "Generate validate source file.")

	cmd.Flags().StringVar(&golangDeepCopyOut, "golang-deepcopy_out", "", "Generate Go deepcopy file.")
	cmd.Flags().StringVar(&golangJsonShimOut, "golang-jsonshim_out", "", "Generate Go JSON shim file.")

	cmd.Flags().StringVar(&gogoFastOut, "gogofast_out", "", "Generate Go source file using gogofast.")

	return cmd.Execute()
}

func mkdir(path string) error {
	if path == "" {
		return nil
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	return nil
}
