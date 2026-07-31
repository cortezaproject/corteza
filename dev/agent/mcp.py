#!/usr/bin/env python3
"""Call the local Human MCP server (streamable HTTP at /api/mcp).

Usage:
  mcp.py tools                     # list tool names + one-line descriptions
  mcp.py schema <tool>             # full input schema of one tool
  mcp.py call <tool> ['<json>']    # call a tool with JSON arguments

Auth: bearer token from token.sh (agent-dev identity). Local-only.
Sessions are per-invocation (initialize → call), which is fine for CLI use.
"""

import json
import os
import subprocess
import sys
import urllib.parse
import urllib.request

AGENT_DIR = os.path.dirname(os.path.abspath(__file__))
API = os.environ.get("HUMAN_API", "http://localhost:1043/api")
MCP = API + "/mcp"

if urllib.parse.urlparse(API).hostname not in ("localhost", "127.0.0.1", "::1"):
    sys.exit(f"mcp.py is local-only; refusing to touch {API}")

TOKEN = subprocess.run(
    [os.path.join(AGENT_DIR, "token.sh")], capture_output=True, text=True, check=True
).stdout.strip()


def post(body, sid=None):
    h = {
        "Content-Type": "application/json",
        "Accept": "application/json, text/event-stream",
        "Authorization": f"Bearer {TOKEN}",
    }
    if sid:
        h["Mcp-Session-Id"] = sid
    req = urllib.request.Request(MCP, data=json.dumps(body).encode(), headers=h, method="POST")
    with urllib.request.urlopen(req, timeout=60) as r:
        raw = r.read().decode()
        new_sid = r.headers.get("Mcp-Session-Id")
    # unwrap SSE framing if present
    if raw.startswith(("event:", "data:")) or "\ndata:" in raw:
        raw = next((l[5:].strip() for l in raw.splitlines() if l.startswith("data:")), "")
    return (json.loads(raw) if raw.strip() else None), new_sid


def rpc(method, params, sid=None, rid=1):
    resp, new_sid = post({"jsonrpc": "2.0", "id": rid, "method": method, "params": params}, sid)
    if resp and "error" in resp:
        sys.exit(f"RPC error: {json.dumps(resp['error'])}")
    return (resp or {}).get("result"), new_sid


def session():
    _, sid = rpc(
        "initialize",
        {"protocolVersion": "2025-03-26", "capabilities": {}, "clientInfo": {"name": "dev-agent-mcp", "version": "1"}},
    )
    post({"jsonrpc": "2.0", "method": "notifications/initialized"}, sid)
    return sid


def main():
    if len(sys.argv) < 2 or sys.argv[1] not in ("tools", "schema", "call"):
        sys.exit(__doc__)
    cmd = sys.argv[1]
    sid = session()
    tools, _ = rpc("tools/list", {}, sid, rid=2)

    if cmd == "tools":
        for t in sorted(tools["tools"], key=lambda t: t["name"]):
            desc = (t.get("description") or "").split("\n")[0]
            print(f"{t['name']:36} {desc[:90]}")
        return

    name = sys.argv[2] if len(sys.argv) > 2 else sys.exit("tool name required")

    if cmd == "schema":
        for t in tools["tools"]:
            if t["name"] == name:
                print(json.dumps(t, indent=1))
                return
        sys.exit(f"unknown tool: {name}")

    args = json.loads(sys.argv[3]) if len(sys.argv) > 3 else {}
    result, _ = rpc("tools/call", {"name": name, "arguments": args}, sid, rid=3)
    if result.get("isError"):
        print("TOOL ERROR:", file=sys.stderr)
    for c in result.get("content", []):
        text = c.get("text", "")
        try:
            print(json.dumps(json.loads(text), indent=1))
        except (json.JSONDecodeError, TypeError):
            print(text)
    if result.get("isError"):
        sys.exit(1)


if __name__ == "__main__":
    main()
