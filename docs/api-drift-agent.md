# API Drift Agent

[![View on GitHub Marketplace](https://img.shields.io/badge/GitHub%20Marketplace-api--drift--agent-blue?logo=github)](https://github.com/marketplace/actions/api-drift-agent)

> **Recommended integration.** The API Drift Agent is the recommended way to solve API drift at scale. Rather than wiring up the engine manually, the agent handles discovery, analysis, and consumer notification automatically.

`api-drift-agent` is a LangGraph-powered agentic workflow that detects breaking API changes in provider PRs and automatically opens GitHub Issues in affected consumer repos — no changes required in consumer repos, no explicit consumer list to maintain.

## How it works

```
Provider PR opened
       │
       ▼
┌─────────────────────────────────────┐
│  Download drift-guard-engine binary │
│  Auto-detect schema type & compare  │
│  (OpenAPI, GraphQL, or gRPC/proto)  │
└─────────────────────────────────────┘
       │ breaking changes found
       ▼
┌─────────────────────────────────────┐
│  Search org for repos that          │
│  reference affected endpoints       │
└─────────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────┐
│  Clone each consumer repo           │
│  Scan for affected files            │
│  Open (or update) a GitHub Issue    │
│  Post DriftAgent Report on PR       │
└─────────────────────────────────────┘
       │ PR re-run / changes fixed
       ▼
┌─────────────────────────────────────┐
│  Close resolved consumer issues     │
│  Update PR comment → all clear ✅   │
└─────────────────────────────────────┘
```

## Prerequisites

- Create a GitHub Personal Access Token (PAT) with `repo` (or `public_repo` for public-only orgs) and `read:org` scopes. This is required to search, clone, and open issues in consumer repos. Add it as a repository secret named `ORG_READ_TOKEN` (**Settings → Secrets and variables → Actions → New repository secret**).
- Optionally, add an `ANTHROPIC_API_KEY` secret to enable Claude-powered risk analysis in the issues the agent opens.

## Usage

Add to your **provider** repo's workflow:

```yaml
name: API Drift Check

on:
  pull_request:

permissions:
  contents: read
  issues: write

jobs:
  drift:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: DriftAgent/api-drift-agent@v1
        with:
          org-read-token: ${{ secrets.ORG_READ_TOKEN }}
          # anthropic-api-key: ${{ secrets.ANTHROPIC_API_KEY }}  # optional: enables AI risk analysis
```

## Inputs

| Input | Required | Description |
|---|---|---|
| `base-schema` | No | Path to schema file (auto-detected if omitted). Supports OpenAPI (`.yaml`/`.yml`/`.json`), GraphQL (`.graphql`/`.gql`), and Protobuf (`.proto`). |
| `head-schema` | No | Path on PR branch (defaults to `base-schema`) |
| `org-read-token` | No | PAT with `repo` (or `public_repo`) + `read:org` scopes. Required to search, clone, and open issues in consumer repos. Falls back to `GITHUB_TOKEN` (cannot open issues in other repos). |
| `anthropic-api-key` | No | Enables Claude risk analysis in opened issues |

## Re-run behaviour

The agent is fully idempotent across CI rebuilds:

| Scenario | PR comment | Consumer issues |
|---|---|---|
| Re-run, same breaking changes | Updated in-place | Updated in-place — no duplicates |
| Re-run, more breaking changes | Updated in-place | Updated in-place |
| PR fixed — breaking changes gone | Updated → ✅ all clear | Closed with "Breaking changes resolved" |
| Clean PR, no previous activity | Nothing posted | Nothing touched |

## Troubleshooting

| Symptom | Cause | Fix |
|---|---|---|
| Action fails: "No API schema found" | Schema file not at a standard path, or generated at runtime and not committed | Set the `base-schema` input explicitly |
| Action fails: "drift-guard-engine failed to diff schemas" | Schema file is invalid or malformed | Validate locally: `drift-guard openapi --base ... --head ...` (or `graphql`/`grpc`) |
| Issues created but no AI explanations | `ANTHROPIC_API_KEY` not set | Set the secret in your repo — the agent runs without it but skips Claude risk analysis |
| No issues created in consumer repos | `org-read-token` not set, or PAT has read-only scope | Set `org-read-token` to a PAT with `repo` (or `public_repo`) + `read:org` scopes |
| No consumers found (public org) | Breaking change path is too generic (e.g. `/v1`) | The agent searches for the first stable path segment — very short or version-only segments may not yield useful results |

## Python CLI

Use this if you want to run the agent locally or integrate it into a non-GitHub CI system. You'll need a diff JSON file produced by `drift-guard-engine` first.

```sh
pip install drift-guard-agent

drift-guard-agent \
  --diff diff.json \
  --org my-org \
  --token $ORG_READ_TOKEN \
  --github-token $ORG_READ_TOKEN \
  --pr 42
```
