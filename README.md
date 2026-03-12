# ms-365-mcp-server (Go Edition)

This repository has been migrated from TypeScript/Node.js to Go.

## Build

```bash
go build ./cmd/ms365-mcp-server
```

## Run

```bash
go run ./cmd/ms365-mcp-server --help
```

## Features retained in the Go rewrite

- CLI flag support for login/logout/account management and mode flags.
- Local account cache persisted in `$HOME/.ms365-mcp/accounts.json`.
- Stdio mode emits MCP-style server metadata JSON.
- HTTP mode includes health and OAuth discovery placeholders.

## Notes

This rewrite intentionally provides a lightweight Go-native baseline so the project no longer depends on Node.js.
