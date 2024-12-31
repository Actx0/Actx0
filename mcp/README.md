# Actx0 MCP

Streamable HTTP MCP proxy for the Actx0 management API (session memories).

## Run

```bash
go run . mcp
```

Defaults:

| Flag | Default | Description |
|------|---------|-------------|
| `--mgt` | `http://127.0.0.1:8000` | Management server URL |
| `--port` | `8001` | MCP HTTP listen port |
| `--path` | `/mcp` | MCP HTTP path |
| `--timeout` | `50` | Request timeout in seconds |

Endpoint: `http://127.0.0.1:8001/mcp`

## Configure Cursor

1. Start the MCP server (`go run . mcp`).
2. Add a project config at `.cursor/mcp.json` (or global `~/.cursor/mcp.json`):

```json
{
  "mcpServers": {
    "actx0": {
      "url": "http://127.0.0.1:8001/mcp",
      "headers": {
        "X-Workspace-Id": "<workspace-id>",
        "X-Agent-Id": "<agent-id>",
        "X-Session-Id": "<external-session-id>",
        "X-Access-Key": "<access-key>"
      }
    }
  }
}
```

3. Reload MCP in **Settings → Tools & MCP**, or restart Cursor.

Do not commit real access keys. Prefer placeholders in shared configs, or keep secrets in a local untracked `.cursor/mcp.json`.

### Headers

| Header | Purpose |
|--------|---------|
| `X-Workspace-Id` | Workspace UUID |
| `X-Agent-Id` | Agent UUID |
| `X-Session-Id` | External session id (created if missing) |
| `X-Access-Key` | Management API access key |

### Access key permissions

The `access key` must include:

- `Get session memories`
- `List session memories`
- `Delete session memories`
- `Update session memories`
- `Create session memories`
- `Query session memories`
- `Create agent sessions`
- `Update agent sessions`
- `Delete agent sessions`
- `List agent sessions`
- `Get agent sessions`
