BIN              := drift-guard
CMD              := ./cmd/drift-guard
HOMEBREW_TAP     := pgomes13/homebrew-tap
FORMULA          := drift-guard

.PHONY: build test vet lint clean run-openapi run-graphql run-grpc release major minor patch gha commit

build:
	go build -o $(BIN) $(CMD)

test:
	go test ./...

vet:
	go vet ./...

lint: vet
	staticcheck ./...

clean:
	rm -f $(BIN)

## Quick smoke runs against the bundled fixtures
run-openapi: build
	./$(BIN) openapi --base internal/testdata/base.yaml --head internal/testdata/head.yaml

run-graphql: build
	./$(BIN) graphql --base internal/testdata/base.graphql --head internal/testdata/head.graphql

run-grpc: build
	./$(BIN) grpc --base internal/testdata/base.proto --head internal/testdata/head.proto

## Commit: stage all changes, commit with a message, and push to the current branch.
##
## Usage:
##   make commit   # prompts for commit message
##
commit:
	@read -p "Commit message: " msg; \
	git add .; \
	git commit -m "$$msg"; \
	git push origin $$(git rev-parse --abbrev-ref HEAD)

## Release targets
##
## Usage:
##   make release homebrew          # bump patch → tag → push (triggers goreleaser + Homebrew update)
##   make release homebrew minor    # bump minor → tag → push
##   make release homebrew major    # bump major → tag → push
##   make release gha               # force-update floating v1 tag for GitHub Action users
##
## Requires: gh CLI (https://cli.github.com) authenticated with repo access.
ifneq (,$(filter major,$(MAKECMDGOALS)))
  _bump := major
else ifneq (,$(filter minor,$(MAKECMDGOALS)))
  _bump := minor
else
  _bump := patch
endif

major minor patch homebrew:
	@true

release:
ifneq (,$(filter gha,$(MAKECMDGOALS)))
	@set -e; \
	LATEST=$$(git describe --tags --abbrev=0 --match "v*.*.*" 2>/dev/null); \
	if [ -z "$$LATEST" ]; then echo "error: no version tag found"; exit 1; fi; \
	MAJOR=$$(echo "$$LATEST" | grep -oE '^v[0-9]+'); \
	echo "Updating floating tag $$MAJOR → $$LATEST"; \
	git tag -f "$$MAJOR"; \
	git push origin "$$MAJOR" --force
else
	@command -v gh >/dev/null 2>&1 || { echo "error: gh CLI not found — install from https://cli.github.com"; exit 1; }
	@set -e; \
	echo "Fetching current version from $(HOMEBREW_TAP)..."; \
	RAW=$$(gh api "repos/$(HOMEBREW_TAP)/contents/$(FORMULA).rb" --jq '.content' | base64 -d); \
	CURRENT=$$(echo "$$RAW" | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1); \
	if [ -z "$$CURRENT" ]; then \
		echo "error: could not parse version from $(FORMULA).rb in $(HOMEBREW_TAP)"; exit 1; \
	fi; \
	MAJOR=$$(echo "$$CURRENT" | cut -d. -f1); \
	MINOR=$$(echo "$$CURRENT" | cut -d. -f2); \
	PATCH=$$(echo "$$CURRENT" | cut -d. -f3); \
	case "$(_bump)" in \
		major) NEXT="v$$((MAJOR + 1)).0.0" ;; \
		minor) NEXT="v$$MAJOR.$$((MINOR + 1)).0" ;; \
		patch) NEXT="v$$MAJOR.$$MINOR.$$((PATCH + 1))" ;; \
	esac; \
	echo "Current: v$$CURRENT  →  Next: $$NEXT  ($(_bump) bump)"; \
	git tag "$$NEXT"; \
	git push origin "$$NEXT"
endif

gha:
	@true
