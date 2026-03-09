.PHONY: build build-mcp build-api clean

build:
	go build -o $(BIN) $(CMD)

build-mcp:
	go build -o $(MCP_BIN) $(MCP_CMD)

build-api:
	go build -o $(API_BIN) $(API_CMD)

clean:
	rm -f $(BIN) $(MCP_BIN) $(API_BIN)
