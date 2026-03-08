# CI Integration

drift-guard is designed to run in CI as a schema diff gate on pull requests.

## GitHub Action (recommended)

Available on the [GitHub Marketplace](https://github.com/marketplace/actions/drift-guard). Add API drift detection to any pull request in one line:

```yaml
- uses: pgomes13/drift-guard-engine@v1
```

The action automatically:

- Detects your framework and runs comparisons for all API types found — REST (OpenAPI), GraphQL, and gRPC
- Posts a PR comment with the full diff report when drift is detected
- Updates a drift log on your GitHub Pages site, with one entry per PR showing timestamps in your local time

See [Supported](/supported) for all supported languages and frameworks.

Full workflow example:

```yaml
name: API Drift Check

on:
  pull_request:

permissions:
  contents: write     # required for updating the drift log on GitHub Pages
  pull-requests: write
  issues: write       # required for posting PR comments

jobs:
  drift:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: pgomes13/drift-guard-engine@v1
```

> **Note:** `contents: write` is required to update the drift log on your GitHub Pages branch. `issues: write` is required because GitHub's PR comment API uses the Issues REST endpoint.

## Live example

See drift-guard in action on a real pull request: [pgomes13/nest-coffee#8](https://github.com/pgomes13/nest-coffee/pull/8)

## Claude Code Plugin

drift-guard is available as a [Claude Code](https://claude.ai/claude-code) plugin, letting you run schema diffing directly from your AI coding session.

### Installation

**1. Install Claude Code**

```bash
npm install -g @anthropic-ai/claude-code
```

**2. Build and install the MCP server binary**

```bash
git clone https://github.com/pgomes13/drift-guard-engine
cd drift-guard-engine
make build-mcp
sudo mv drift-guard-mcp /usr/local/bin/
```

**3. Install the plugin**

```bash
claude plugin install https://github.com/pgomes13/drift-guard-engine
```

### Available tools

Once installed, Claude Code gains four tools:

| Tool | Description |
|------|-------------|
| `diff_openapi` | Compare two OpenAPI 3.x schema files |
| `diff_graphql` | Compare two GraphQL SDL files |
| `diff_grpc` | Compare two Protobuf `.proto` files |
| `detect_project` | Detect project type, framework, and available schema types |

All diff tools accept an optional `format` parameter: `text` (default), `json`, or `markdown`.

### Example usage

Ask Claude Code:

> "Compare `docs/base.yaml` and `docs/head.yaml` for breaking API changes"

> "What schema types does this project use?"

## Key flags

| Flag | Purpose |
|---|---|
| `--format markdown` | Renders a Markdown table — used by the action for PR comments |
| `--format github` | Renders inline PR annotations via workflow commands |
| `--fail-on-breaking` | Exits with code `1` to block merges on breaking changes |
| `--format json` | Use if you need to parse output in a subsequent step |
