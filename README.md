# protog: a drop-in replacement for protoc that manages dependencies

[protoc](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation) is the official tool for compiling
protocol buffers. While it is simple enough to download it and compile a helloworld.proto, it quickly becomes difficult
to manage when the proto file becomes more complex and when generating code using protoc plugins, for example gRPC.

protog attempts to be a drop-in replacement for protoc that "just works". It has three main features:

- Automatically download any plugin specified in the command line
- Automatically download proto files used in imports
- Small quality of life improvements such as automatically creating directories for generated code, which if missing
otherwise cause `protoc` to crash

## Supported Platforms

protog is supported on the following platforms

- linux-amd64
- linux-arm64
  - grpc/grpc plugins such as `grpc_python_plugin` are not currently supported as they are not published for arm64
- darwin-amd64
- darwin-arm64 (currently only verified locally and not on CI)
  - grpc/grpc plugins such as `grpc_python_plugin` are only available for amd64 but work through Rosetta
- windows-amd64

windows-arm64 will come in the future when it becomes a viable platform.

## Supported Plugins

The following protoc plugins are supported and will be automatically installed if used. Support can be added simply
for any plugin that provides prebuilt binaries or can be built with Go or NodeJS, and we would be happy to add support
when needed.

- [Doc](https://github.com/pseudomuto/protoc-gen-doc) - `--doc_out`
- [Docs (Istio)](https://github.com/istio/tools/tree/master/cmd/protoc-gen-docs) - `--docs_out`
- [Go](https://github.com/protocolbuffers/protobuf-go) - `--go_out`
- [GogoFast](https://github.com/gogo/protobuf) -- `gogofast_out`
- [Go Deep Copy](https://github.com/istio/tools/tree/master/cmd/protoc-gen-golang-deepcopy) - `--golang-deepcopy_out`
- [Go JSON Shim](https://github.com/istio/tools/tree/master/cmd/protoc-gen-golang-deepcopy) - `--golang-jsonshim_out`
- [gRPC](https://github.com/grpc/grpc) - `--grpc_<lang>_out`
  - Includes `cpp`, `csharp`, `js`, `objc`, `php`, `python`, and `ruby`
  - Does not currently work on linux-arm64
- [gRPC Gateway](https://github.com/grpc-ecosystem/grpc-gateway) - `--grpc-gateway_out`
- [gRPC Go](https://github.com/grpc/grpc-go) - `--go-grpc_out`
- [gRPC Java](https://github.com/grpc/grpc-java) - `--grpc-java_out`
- [gRPC Web](https://github.com/grpc/grpc-web) -- `--grpc-web_out`
- [JSON Schema](https://github.com/chrusty/protoc-gen-jsonschema) - `--jsonschema_out`
- [TypeScript](https://github.com/thesayyn/protoc-gen-ts) - `--ts_out`
- [TypeScript (Improbable ver)](https://github.com/improbable-eng/ts-protoc-gen) - `--improbable_ts_out`
- [Validate](https://github.com/envoyproxy/protoc-gen-validate) - `--validate_out`
