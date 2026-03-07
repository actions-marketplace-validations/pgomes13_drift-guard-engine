BIN              := drift-guard
CMD              := ./cmd/drift-guard
HOMEBREW_TAP     := pgomes13/homebrew-tap
FORMULA          := drift-guard

.PHONY: build test vet lint clean run-openapi run-graphql run-grpc release

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

## Release: bump patch version based on current homebrew tap formula, tag, and push.
## Requires: gh CLI (https://cli.github.com) authenticated with repo access.
release:
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
	NEXT="v$$MAJOR.$$MINOR.$$((PATCH + 1))"; \
	echo "Current: v$$CURRENT  →  Next: $$NEXT"; \
	git tag "$$NEXT"; \
	git push origin "$$NEXT"
