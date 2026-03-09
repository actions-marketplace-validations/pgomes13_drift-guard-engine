# npm SDK

`@pgomes13/drift-guard` is a thin npm wrapper around the drift-guard binary. On install, the correct pre-built binary for your platform is downloaded automatically — no Go toolchain required.

## Installation

```sh
npm install @pgomes13/drift-guard
```

Requires Node.js ≥ 16. Supported platforms: macOS arm64/amd64, Linux arm64/amd64, Windows amd64.

## Programmatic API

### OpenAPI

```ts
import { compareOpenAPI } from "@pgomes13/drift-guard";

const result = compareOpenAPI("old.yaml", "new.yaml");

console.log(result.summary);
// { total: 3, breaking: 1, non_breaking: 2, info: 0 }

for (const change of result.changes) {
  console.log(`[${change.severity}] ${change.description}`);
}
```

### GraphQL

```ts
import { compareGraphQL } from "@pgomes13/drift-guard";

const result = compareGraphQL("old.graphql", "new.graphql");
```

### gRPC / Protobuf

```ts
import { compareGRPC } from "@pgomes13/drift-guard";

const result = compareGRPC("old.proto", "new.proto");
```

### Impact analysis

Scan source code for references to each breaking change:

```ts
import { compareOpenAPI, impact } from "@pgomes13/drift-guard";

const result = compareOpenAPI("old.yaml", "new.yaml");

// Returns Hit[] — one entry per matching file:line
const hits = impact(result, "./src");

for (const hit of hits) {
  console.log(`${hit.file}:${hit.line_num}  (${hit.change_path})`);
}
```

Text or markdown report:

```ts
const report = impact(result, "./src", { format: "markdown" });
console.log(report);
```

## CLI via npx

The `drift-guard` binary is available as an npm bin after install:

```sh
npx drift-guard openapi --base old.yaml --head new.yaml
npx drift-guard graphql --base old.graphql --head new.graphql --format json
npx drift-guard impact --diff diff.json --scan ./src
```

## CommonJS

```js
const { compareOpenAPI, impact } = require("@pgomes13/drift-guard");
```

## TypeScript types

```ts
type Severity = "breaking" | "non-breaking" | "info";

interface Change {
  type: string;        // e.g. "endpoint_removed", "field_type_changed"
  severity: Severity;
  path: string;        // e.g. "/users/{id}"
  method: string;      // e.g. "DELETE"
  location: string;
  description: string;
  before?: string;
  after?: string;
}

interface Summary {
  total: number;
  breaking: number;
  non_breaking: number;
  info: number;
}

interface DiffResult {
  base_file: string;
  head_file: string;
  changes: Change[];
  summary: Summary;
}

interface Hit {
  file: string;
  line_num: number;
  line: string;
  change_type: string;  // e.g. "endpoint_removed"
  change_path: string;  // e.g. "DELETE /users/{id}"
}
```
