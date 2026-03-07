BIN              := drift-guard
CMD              := ./cmd/drift-guard
HOMEBREW_TAP     := pgomes13/homebrew-tap
FORMULA          := drift-guard

.PHONY: build test vet lint clean run-openapi run-graphql run-grpc release major minor patch gha homebrew commit override

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
##   make release          # bump patch → tag → push → update floating major tag
##   make release minor    # bump minor → tag → push → update floating major tag
##   make release major    # bump major → tag → push → update floating major tag
##   make release override # re-tag current version (force) → push → update floating major tag
##   make release gha      # force-update floating major tag only (no version bump)
##
## Pushing the semver tag triggers the release.yml workflow (goreleaser + Homebrew update).
ifneq (,$(filter major,$(MAKECMDGOALS)))
  _bump := major
else ifneq (,$(filter minor,$(MAKECMDGOALS)))
  _bump := minor
else
  _bump := patch
endif

major minor patch homebrew override:
	@true

release:
ifneq (,$(filter gha,$(MAKECMDGOALS)))
	@set -e; \
	LATEST=$$(git tag --list 'v*.*.*' --points-at HEAD --sort=-version:refname | head -1); \
	if [ -z "$$LATEST" ]; then LATEST=$$(git describe --tags --abbrev=0 --match "v*.*.*" 2>/dev/null); fi; \
	if [ -z "$$LATEST" ]; then echo "error: no version tag found"; exit 1; fi; \
	FLOAT=$$(echo "$$LATEST" | grep -oE '^v[0-9]+'); \
	echo "Updating floating tag $$FLOAT → $$LATEST"; \
	git tag -f "$$FLOAT"; \
	git push origin "$$FLOAT" --force
else ifneq (,$(filter override,$(MAKECMDGOALS)))
	@set -e; \
	CURRENT=$$(git tag --list 'v*.*.*' --points-at HEAD --sort=-version:refname | head -1); \
	if [ -z "$$CURRENT" ]; then CURRENT=$$(git describe --tags --abbrev=0 --match "v*.*.*" 2>/dev/null); fi; \
	if [ -z "$$CURRENT" ]; then \
		echo "error: no semver tag found in repo (expected v<major>.<minor>.<patch>)"; exit 1; \
	fi; \
	echo "Re-tagging: $$CURRENT (force)"; \
	git tag -f "$$CURRENT"; \
	git push origin "$$CURRENT" --force; \
	FLOAT=$$(echo "$$CURRENT" | grep -oE '^v[0-9]+'); \
	echo "Updating floating tag $$FLOAT → $$CURRENT"; \
	git tag -f "$$FLOAT"; \
	git push origin "$$FLOAT" --force
else
	@set -e; \
	_TAG=$$(git tag --list 'v*.*.*' --points-at HEAD --sort=-version:refname | head -1); \
	if [ -z "$$_TAG" ]; then _TAG=$$(git describe --tags --abbrev=0 --match "v*.*.*" 2>/dev/null); fi; \
	CURRENT=$$(echo "$$_TAG" | sed 's/^v//'); \
	if [ -z "$$CURRENT" ]; then \
		echo "error: no semver tag found in repo (expected v<major>.<minor>.<patch>)"; exit 1; \
	fi; \
	echo "Current: v$$CURRENT"; \
	MAJOR=$$(echo "$$CURRENT" | cut -d. -f1); \
	MINOR=$$(echo "$$CURRENT" | cut -d. -f2); \
	PATCH=$$(echo "$$CURRENT" | cut -d. -f3); \
	case "$(_bump)" in \
		major) NEXT="v$$((MAJOR + 1)).0.0" ;; \
		minor) NEXT="v$$MAJOR.$$((MINOR + 1)).0" ;; \
		patch) NEXT="v$$MAJOR.$$MINOR.$$((PATCH + 1))" ;; \
	esac; \
	echo "Next:    $$NEXT  ($(_bump) bump)"; \
	git tag -f "$$NEXT"; \
	git push origin "$$NEXT" --force; \
	FLOAT=$$(echo "$$NEXT" | grep -oE '^v[0-9]+'); \
	echo "Updating floating tag $$FLOAT → $$NEXT"; \
	git tag -f "$$FLOAT"; \
	git push origin "$$FLOAT" --force
endif

gha:
	@true
