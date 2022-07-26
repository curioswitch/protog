package protog

import (
	"github.com/curioswitch/protog/internal/cmd"
	"github.com/curioswitch/protog/internal/tools"
)

type Versions = tools.Versions

func Run(args []string, versions Versions) error {
	// round-tripping versions through env and back is a bit weird but keeps things simplest since we need to
	// parse args even for programmatic invocation.
	env := map[string]string{}
	env["GO_VERSION"] = versions.Go
	env["NODEJS_VERESION"] = versions.NodeJS
	env["PROTOC_VERSION"] = versions.Protoc
	env["PROTOC_GEN_DOC_VERSION"] = versions.ProtocGenDoc
	env["PROTOC_GEN_GO_VERSION"] = versions.ProtocGenGo
	env["PROTOC_GEN_GO_GRPC_VERSION"] = versions.ProtocGenGoGRPC
	env["PROTOC_GEN_GRPC_VERSION"] = versions.ProtocGenGRPC
	env["PROTOC_GEN_GRPC_GATEWAY_VERSION"] = versions.ProtocGenGRPCGateway
	env["PROTOC_GEN_GRPC_JAVA_VERSION"] = versions.ProtocGenGRPCJava
	env["PROTOC_GEN_TS_VERSION"] = versions.ProtocGenTS
	env["PROTOC_GEN_VALIDATE_VERSION"] = versions.ProtocGenValidate

	return cmd.Run(args, env)
}
