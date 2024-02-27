# Deprecated

This project is deprecated in favor of the approach used in wasilibs. It should provide a faster, more stable experience.

https://github.com/wasilibs/go-protoc-gen-builtins

# protog: a drop-in replacement for protoc that manages dependencies

[protoc](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation) is the official tool for compiling
protocol buffers. While it is simple enough to download it and compile a helloworld.proto, it quickly becomes difficult
to manage when the proto file becomes more complex and when generating code using protoc plugins, for example gRPC.

protog attempts to be a drop-in replacement for protoc that "just works". It has three main features:

- Automatically download any plugin specified in the command line
- Automatically download proto files used in imports
- Small quality of life improvements such as automatically creating directories for generated code, which if missing
otherwise cause `protoc` to crash

## Installation

You can grab a precompiled binary from the [releases](https://github.com/curioswitch/protog/releases).

We also offer docker images on [ghcr.io](https://github.com/curioswitch/protog/pkgs/container/protog).

protog can also be invoked [programatically](./protog.go). This is generally the most convenient option for
[mage](https://magefile.org) users.

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

By default, the latest version of the plugin is determined and fetched when missing. This can be overridden by specifying
the appropriate version environment variable for the plugin.

| Plugin                                                                                    | Command line flag       | Version environment variable         |
|-------------------------------------------------------------------------------------------|-------------------------|--------------------------------------|
| [Doc](https://github.com/pseudomuto/protoc-gen-doc)                                       | `--doc_out`             | `PROTOC_GEN_DOC_VERSION`             |
| [Docs (Istio)](https://github.com/istio/tools/tree/master/cmd/protoc-gen-docs)            | `--docs_out`            | `PROTOC_GEN_DOCS_VERSION`            |
| [Go](https://github.com/protocolbuffers/protobuf-go)                                      | `--go_out`              | `PROTOC_GEN_GO_VERSION`              |
| [GogoFast](https://github.com/gogo/protobuf)                                              | `--gogofast_out`        | `PROTOC_GEN_GOGO_FAST_VERSION`       |
| [Go Deep Copy](https://github.com/istio/tools/tree/master/cmd/protoc-gen-golang-deepcopy) | `--golang-deepcopy_out` | `PROTOC_GEN_GOLANG_DEEPCOPY_VERSION` |
| [Go JSON Shim](https://github.com/istio/tools/tree/master/cmd/protoc-gen-golang-jsonshim) | `--golang-jsonshim_out` | `PROTOC_GEN_GOLANG_JSONSHIM_VERSION` |
| [gRPC](https://github.com/grpc/grpc)[^1]                                                  | `--grpc_<lang>_out`     | `PROTOC_GEN_GRPC_VERSION`            |
| [gRPC Gateway](https://github.com/grpc-ecosystem/grpc-gateway)                            | `--grpc-gateway_out`    | `PROTOC_GEN_GRPC_GATEWAY_VERSION`    |
| [gRPC Go](https://github.com/grpc/grpc-go)                                                | `--go-grpc_out`         | `PROTOC_GEN_GO_GRPC_VERSION`         |
| [gRPC Java](https://github.com/grpc/grpc-java)                                            | `--grpc-java_out`       | `PROTOC_GEN_GRPC_JAVA_VERSION`       |
| [gRPC Web](https://github.com/grpc/grpc-web)                                              | `--grpc-web_out`        | `PROTOC_GEN_GRPC_WEB_VERSION`        |
| [JSON Schema](https://github.com/chrusty/protoc-gen-jsonschema)                           | `--jsonschema_out`      | `PROTOC_GEN_JSONSCHEMA_VERSION`      |
| [TypeScript](https://github.com/thesayyn/protoc-gen-ts)                                   | `--ts_out`              | `PROTOC_GEN_TS_VERSION`              |
| [TypeScript (Improbable ver)](https://github.com/improbable-eng/ts-protoc-gen)            | `--improbable_ts_out`   | `PROTOC_TS_GEN_VERSION`              |
| [Validate](https://github.com/envoyproxy/protoc-gen-validate)                             | `--validate_out`        | `PROTOC_GEN_VALIDATE_VERSION`        |

[^1]: Includes `cpp`, `csharp`, `js`, `objc`, `php`, `python`, and `ruby`. Does not currently work on linux-arm64.

Note that anytime an upstream artifact follows the protoc plugin naming convention `protoc-gen-<plugin>`, we use the same
for the command line flag for activating it. This means some inconsistent flags such as `grpc-java_out` vs `go-grpc_out`, but
we choose to remain consistent with the upstream naming where possible.

gRPC plugins do not follow this naming scheme, instead being named like `grpc_<lang>_plugin`. We have decided to use the
scheme of prefixing the `--<lang>_out` flag of protoc itself with `grpc_`, e.g. `--grpc_js_out`. Note that this means the
flag is significantly different for `objc` vs `grpc_objective_c_plugin` and `js` vs `grpc_node_plugin`.

There are two popular TypeScript plugins named `protoc-gen-ts`. We decided to give `--ts_out` to the one where the GitHub repository
name is also `protoc-gen-ts`, but this is not meant to indicate one is better than the other. If both repository names also matched,
we would have needed to flip a coin to decide which to use as the default TypeScript plugin.

Plugins are automatically downloaded to the [user cache dir](https://pkg.go.dev/os#UserCacheDir).

## Supported proto imports

The following protos can be fetched automatically when imported. In the future, protog will allow specifying additional
mappings when run. We would always be happy to add more entries to the built-in registry for open source protos when
needed.

| Prefix              | Repository                                        |
|---------------------|---------------------------------------------------|
| google/api          | https://github.com/googleapis/googleapis          |
| google/rpc          | https://github.com/googleapis/googleapis          |
| gogoproto           | https://github.com/gogo/protobuf                  |
| k8s.io/api          | https://github.com/kubernetes/api                 |
| k8s.io/apimachinery | https://github.com/kubernetes/apimachinery        |
| validate            | https://github.com/envoyproxy/protoc-gen-validate |

Imported protos are by default downloaded to `build/proto-includes` within the current working directory. The reason to
not use the same cache directory as plugins is by being relative to the project, IDEs can recognize the protos and load
them for completion. For example, in Jetbrains IDEs, an `alt-enter` on a missing import will automatically add this
folder to the search path.

The directory can be changed by providing the `PROTO_INCLUDES_DIR` environment variable.

## Additional Configuration

When needed, protog will download Golang or NodeJS for building missing plugins. The versions can be pinned using the
environment variables `GO_VERSION` and `NODEJS_VERSION` respectively.

## How it works

protog is not a reimplementation of protoc in Go, as cool as that would be :-) It is generally a package manager for
automatically downloading all the artifacts needed to execute a protobuf compilation. protog always downloads protoc
automatically from its published artifacts. It also parses the command line for known `<plugin>_out` options. For any it
finds, it checks against an included registry of [specs](internal/tools/specs.go) to determine how to download or build
the artifact and will do so. If building a plugin requires Go or NodeJS, it will download that too.

It also parses the command line for the proto files that are being built and scans them for `import` statements. It
compares the import statement to the included registry of [includes](internal/proto/includes.go) and if matches, downloads
the protos. Include protos do not support versioning because using version numbers is not common with protos. It would
also mean having different downloaded protos in different folders, but we've found IDEs behave most smoothly by having
them all in one, unversioned, folder.

## Alternatives

### bingo

[bingo](https://github.com/bwplotka/bingo) is a great tool for managing dependencies written in Go. Unfortunately, many
protoc plugins are not written in Go and bingo cannot manage them.

### buf

[buf](https://buf.build/) is a reimagining of protobuf compilation. It has many cool features like linting and checking
for backwards incompatible changes. Many will find it convenient and full-featured and will use it instead of protog. 
protog attempts to be an alternative to buf with these main differences:

- Drop-in replacement for protoc, so existing builds can switch to protog more easily
- Only needs to be online when fetching plugins, which happens when building for the first time or updating versions.
Buf remote plugins execute remotely at build time and prevent offline builds. Not using remote plugins requires tedious
plugin management.
- Only use the source of truth. protog downloads artifacts or source from official sources while Buf maintains a completely
separate registry. There are tradeoffs to both approach, but we believe ownership is easier to follow by using the
official sources. For example, when searching `protoc-gen-validate` on BSR, we can only find several uploads all from
different end-users, and no official source of truth.

If these are not issues for you, Buf is a great tool, and we recommend using it!
