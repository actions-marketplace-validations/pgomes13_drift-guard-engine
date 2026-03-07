# CLI

Run drift-guard locally to check for API drift before adding it to CI.

## Install

See [Installation](/install) for all install options.

## Run on your repository

From the root of your project, run:

```sh
drift-guard compare
```

drift-guard will auto-detect your framework, generate schemas for your current branch and `origin/main`, and print a diff. This is a good way to verify it works with your project before wiring up the GitHub Action.

### Check for breaking changes only

```sh
drift-guard compare --fail-on-breaking
```

Exits with code `1` if any breaking changes are found — same behavior as in CI.

### Markdown output

```sh
drift-guard compare --format markdown
```

Renders the same table that gets posted as a PR comment in CI.
