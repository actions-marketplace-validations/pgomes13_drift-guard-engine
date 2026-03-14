# api-drift-engine

API type safety across **OpenAPI**, **GraphQL**, and **gRPC** — catch breaking changes before they reach production.

**[Full documentation →](https://driftbot.github.io/api-drift-engine)**

## Quick install

```sh
# Homebrew
brew tap DriftBot/tap
brew install drift-bot

# npm
npm install @drift-bot/api-drift-engine
```

## Quick start

```sh
# Auto-generate and compare specs between branches
drift-bot compare

# GitHub Action — one line
- uses: DriftBot/api-drift-engine@v1
```

## npm / Node.js API

```ts
import { compareOpenAPI, impact } from "@drift-bot/api-drift-engine";

const result = compareOpenAPI("old.yaml", "new.yaml");
const hits = impact(result, "./src");
```
