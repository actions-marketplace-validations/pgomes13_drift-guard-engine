# Development

## Make commands

```sh
make build       # compile binary
make test        # run all tests
make vet         # run go vet
make lint        # run go vet + staticcheck
make clean       # remove binary

make run-openapi  # build and diff bundled OpenAPI fixtures
make run-graphql  # build and diff bundled GraphQL fixtures
make run-grpc     # build and diff bundled gRPC fixtures
```

## Architecture

```
cmd/drift-guard/          # CLI entry point (drift-guard binary)
cmd/server/               # gRPC server entry point
api/drift-guard/v1/       # Protobuf service definition & generated Go code
internal/
  parser/
    openapi/             # OpenAPI YAML/JSON → schema.Schema
    graphql/             # GraphQL SDL → schema.GQLSchema
    grpc/                # Protobuf .proto → schema.GRPCSchema
  differ/
    openapi/             # Diffs two schema.Schema values
    graphql/             # Diffs two schema.GQLSchema values
    grpc/                # Diffs two schema.GRPCSchema values
  classifier/
    openapi/             # Assigns severity to OpenAPI changes
    graphql/             # Assigns severity to GraphQL changes
    grpc/                # Assigns severity to gRPC changes
  reporter/              # Renders DiffResult as text / JSON / GitHub annotations
pkg/schema/
  types.schema.go        # Shared types: Change, DiffResult, Severity
  openapi.schema.go      # OpenAPI types and change type constants
  graphql.schema.go      # GraphQL types and change type constants
  grpc.schema.go         # gRPC types and change type constants
```

## Releasing a new version

```sh
# Homebrew release — bump patch / minor / major
make release homebrew
make release homebrew minor
make release homebrew major

# GitHub Action release — update floating v1 tag
make release gha
```

`make release homebrew` fetches the current version from the Homebrew tap, bumps it, tags the commit, and pushes the tag. This triggers the `release.yml` workflow which cross-compiles binaries for macOS, Linux, and Windows via [GoReleaser](https://goreleaser.com), publishes a GitHub Release, pushes the Homebrew formula to [`pgomes13/homebrew-tap`](https://github.com/pgomes13/homebrew-tap), and updates the floating `v1` major tag used by the GitHub Action.
