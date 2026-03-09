.PHONY: run-api run-openapi run-graphql run-grpc local playground

run-api: build-api
	./$(API_BIN)

## local: build Go API + run Next.js playground side-by-side.
## Go API → http://localhost:9000   Next.js → http://localhost:3000
##
## Usage:
##   make local
##
local: build-api
	@echo "Starting Go API on :9000 and Next.js on :3000 (Ctrl+C to stop both)"
	@trap 'kill 0' SIGINT SIGTERM; \
	./$(API_BIN) & \
	cd playground && npm run dev; \
	wait

## playground: run Go API with embedded static playground (no Node.js required).
## Playground → http://localhost:9000
##
## Usage:
##   make playground
##
playground: build-api
	@echo "Starting standalone playground on http://localhost:9000"
	./$(API_BIN)

## Quick smoke runs against the bundled fixtures
run-openapi: build
	./$(BIN) openapi --base internal/testdata/base.yaml --head internal/testdata/head.yaml

run-graphql: build
	./$(BIN) graphql --base internal/testdata/base.graphql --head internal/testdata/head.graphql

run-grpc: build
	./$(BIN) grpc --base internal/testdata/base.proto --head internal/testdata/head.proto
