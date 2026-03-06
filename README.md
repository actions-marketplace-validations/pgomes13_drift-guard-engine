# drift-guard-engine

A schema diff engine that detects and classifies breaking vs. non-breaking API contract changes across **OpenAPI**, **GraphQL**, and **gRPC** schemas.

**[Full documentation](https://pgomes13.github.io/drift-guard-engine)**

## Quick install

```sh
brew tap pgomes13/tap
brew install drift-guard
```

## Quick start

```sh
drift-guard openapi --base api/base.yaml --head api/head.yaml --format github --fail-on-breaking
```

## Release

Tag a version to trigger GoReleaser — cross-compiles for macOS, Linux, and Windows and publishes to the Homebrew tap:

```sh
git tag v1.0.0
git push origin v1.0.0
```
