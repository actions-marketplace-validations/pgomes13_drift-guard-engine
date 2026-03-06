# CI Integration

drift-guard is designed to run in CI as a schema diff gate on pull requests.

## GitHub Actions

### Install and run

```yaml
- name: Set up Go
  uses: actions/setup-go@v5
  with:
    go-version: "1.22"

- name: Install drift-guard
  run: go install github.com/pgomes13/drift-guard-engine/cmd/driftengine@latest

- name: Check for breaking API changes
  run: |
    drift-guard openapi \
      --base api/base.yaml \
      --head api/head.yaml \
      --format github \
      --fail-on-breaking
```

### Full workflow example

```yaml
name: API Drift Check

on:
  pull_request:

permissions:
  contents: read

jobs:
  drift:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"

      - name: Install drift-guard
        run: go install github.com/pgomes13/drift-guard-engine/cmd/driftengine@latest

      - name: OpenAPI drift check
        run: |
          drift-guard openapi \
            --base api/base.yaml \
            --head api/head.yaml \
            --format github \
            --fail-on-breaking

      - name: GraphQL drift check
        run: |
          drift-guard graphql \
            --base schema/base.graphql \
            --head schema/head.graphql \
            --format github \
            --fail-on-breaking
```

## Key flags for CI

| Flag | Purpose |
|---|---|
| `--format github` | Renders inline PR annotations via workflow commands |
| `--fail-on-breaking` | Exits with code `1` to block merges on breaking changes |
| `--format json` | Use if you need to parse output in a subsequent step |

## Install via Homebrew (alternative)

```yaml
- name: Install drift-guard
  run: |
    brew tap pgomes13/tap
    brew install drift-guard
```

> Homebrew is pre-installed on `macos-*` runners. Works on `ubuntu-latest` but is slower than `go install`.
