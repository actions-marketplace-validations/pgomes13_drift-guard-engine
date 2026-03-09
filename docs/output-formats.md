# Output Formats

Formats apply to both schema diff commands (`openapi`, `graphql`, `grpc`) and the `impact` command. Not all formats are available for both — see the tables below.

## Schema diff formats

| Format | `openapi` / `graphql` / `grpc` | `impact` |
|--------|:---:|:---:|
| `text` | ✓ | ✓ |
| `json` | ✓ | ✓ |
| `markdown` | ✓ | ✓ |
| `github` | ✓ | ✓ |

---

## `text` (default)

```
Total: 4  Breaking: 2  Non-Breaking: 1  Info: 1

SEVERITY        TYPE                    PATH            DESCRIPTION
----------------------------------------------------------------------------------------------------
[BREAKING]      endpoint_removed        /users/{id}     Endpoint '/users/{id}' method DELETE was removed
[BREAKING]      param_type_changed      /users/{id}     Param 'id' type changed from 'string' to 'integer'
[non-breaking]  endpoint_added          /posts          Endpoint '/posts' was added
[info]          field_added             /users          Field 'role' was added
```

**Impact (`text`):**

```
Breaking change: DELETE /users/{id} (endpoint_removed)
  services/client.go:42    client.Delete("/users/" + id)
  apps/routes.go:17        r.DELETE("/users/:id", handler)
```

---

## `json`

```json
{
  "base_file": "base.yaml",
  "head_file": "head.yaml",
  "changes": [
    {
      "type": "endpoint_removed",
      "severity": "breaking",
      "path": "/users/{id}",
      "method": "DELETE",
      "location": "",
      "description": "Endpoint '/users/{id}' method DELETE was removed"
    }
  ],
  "summary": {
    "total": 4,
    "breaking": 2,
    "non_breaking": 1,
    "info": 1
  }
}
```

**Impact (`json`)** — returns a flat array of hits:

```json
[
  {
    "file": "services/client.go",
    "line_num": 42,
    "line": "client.Delete(\"/users/\" + id)",
    "change_type": "endpoint_removed",
    "change_path": "DELETE /users/{id} (endpoint_removed)"
  }
]
```

---

## `github`

Emits [GitHub Actions workflow commands](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/workflow-commands-for-github-actions) that render as inline annotations on the PR diff.

**Schema diff (`github`):**

```
::error title=Breaking Change::Endpoint '/users/{id}' method DELETE was removed
::warning title=Non-Breaking Change::Endpoint '/posts' was added
```

**Impact (`github`)** — one annotation per hit, pointing to the exact file and line:

```
::error file=services/client.go,line=42,title=Breaking API change%3A DELETE /users/{id}::client.Delete("/users/" + id)
::error file=apps/routes.go,line=17,title=Breaking API change%3A DELETE /users/{id}::r.DELETE("/users/:id", handler)
```

Each line appears as a red underline on the specific file and line in the GitHub PR diff view. Use `--format github` in CI workflows to surface impact hits directly in the code review.

---

## `markdown`

Renders GitHub-flavored Markdown — used by the [GitHub Action](./ci.md) for automatic PR comments.

**Schema diff (`markdown`):**

```markdown
| Severity | Type | Path | Description |
|----------|------|------|-------------|
| [BREAKING] | endpoint_removed | /users/{id} | Endpoint '/users/{id}' method DELETE was removed |
| [non-breaking] | endpoint_added | /posts | Endpoint '/posts' was added |
```

**Impact (`markdown`)** — summary line + collapsible section per breaking change:

```markdown
> **3** reference(s) to breaking changes across **2** file(s)

<details>
<summary>🔴 DELETE /users/{id} (endpoint_removed) — 2 reference(s)</summary>

| File | Line | Code |
|------|------|------|
| `services/client.go` | 42 | `client.Delete("/users/" + id)` |
| `apps/routes.go` | 17 | `r.DELETE("/users/:id", handler)` |

</details>
```
