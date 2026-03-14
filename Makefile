BIN              := driftabot
MCP_BIN          := driftabot-mcp
API_BIN          := driftabot-api
CMD              := ./cmd/driftabot
MCP_CMD          := ./cmd/mcp-server
API_CMD          := ./cmd/playground
HOMEBREW_TAP     := DriftaBot/homebrew-tap
FORMULA          := driftabot

include make/build.mk
include make/test.mk
include make/run.mk
include make/deploy.mk
include make/git.mk
include make/release.mk
