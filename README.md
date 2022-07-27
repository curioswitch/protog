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
- darwin-amd64
- darwin-arm64 (currently only verified locally and not on CI)
- windows-amd64

linux-arm64 support will come soon. windows-arm64 will come in the future when it becomes a viable platform.
