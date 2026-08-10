---
kind: folder
covers: recursive
owner: be
depends-on: []
touched-by:
  - server/system/agentic/mcp
  - server/app/boot_levels.go
  - dev/mcp
tests:
  - server/pkg/mcpkit
  - server/tests/mcp
---

# server/pkg/mcpkit

## Responsibility

The transport-level machinery of an MCP server, with no knowledge of Human.

Two servers are built on it: the configurator MCP (`server/system/agentic/mcp`),
which is product code exposing Human to agents, and the developer MCP
(`dev/mcp`), a separate Go module that operates on this repository. They share
the registry, the tagging vocabulary, scope, progressive disclosure and both
transports; they share nothing about what a tool does.

## Locked contracts

- **No Human domain imports.** `pkg/` is allowed; anything under `compose/`,
  `system/`, `automation/` or `app/` is not. `boundary_test.go` enforces it, and
  the enforcement is verified to fail when violated rather than assumed to work.

  This is what keeps `dev/mcp`'s vendor tree at two packages instead of the
  server's entire object graph. When something here needs a Human type, invert
  it: expose a neutral seam and let the domain side project onto it, the way
  `Registry.Select` returns `ToolDef` and `system/agentic/mcp` turns those into
  the agentic runtime's `rt.Tool`.

- **`WithRisk` is the sole writer of the four protocol annotation hints.** Tool
  authors set a risk; they never set `readOnlyHint`, `destructiveHint`,
  `idempotentHint` or `openWorldHint` by hand.

- **Group is presentation; risk is enforced.** Group shrinks what `tools/list`
  returns and nothing more. The risk ceiling refuses at dispatch — a filter
  alone would not, because mcp-go filters apply to listing, not to calls.

## Key types

- `Registry` — name → tool, handler, tags. Duplicate registration panics: two
  handlers claiming one name previously produced a green build with one tool
  silently missing.
- `Group` / `Risk` (`tool_options.go`) — carried on `mcp.Tool.Meta` under the
  `human.dev/` keys, via `InGroup` / `WithRisk`.
- `Scope` (`scope.go`) — per-request narrowing parsed from the URL. There is no
  URL on stdio, and the zero Scope permits everything, which is deliberate: a
  server the developer launched in their own checkout has no caller to narrow.
- `disclosure` (`disclosure.go`) — per-session set of pulled-in tools, keyed by
  MCP session ID and dropped on teardown.
- `MCPServer` (`server.go`) — wraps mcp-go. `MountRoutes` for HTTP,
  `ServeStdio` for a client that launches the server as a child process.
- `toolkit/` — the argument, result and error plumbing every handler repeats.
  `JSONResult` is the only sanctioned path from a Go value to a tool result, so
  the size ceiling lives in one place.

## Progressive disclosure

A session sees five tools until it searches. That is the difference between
~1,400 tokens and ~39,000 per request, and it stays flat as the surface grows.

`searchTools` ranks by how many of the caller's terms matched, and requires at
least one. It used to require **all** of them, which made a natural multi-word
question the worst possible input — "workflow create tool" returned nothing
while "workflow" returned the family — and a model reading "no tool matches"
concludes the capability does not exist.

Search returns full definitions inline **and** fires `list_changed`, so a client
that ignores the notification still works.

## Store / DAL usage

None, and there cannot be any: the store is Human's, and this package may not
import it.

## When changing this

- Both servers are affected. `dev/mcp` is a separate module but compiles these
  packages **from source** through its `replace` directive, so a change here
  reaches it as soon as its binary is rebuilt — `dev/mcp/run.sh` watches this
  folder alongside its own sources, which is the only reason a change here is
  not silently absent from the developer MCP.

  It used to vendor them instead. That duplicated this package into git, and
  because the staleness check compared mtimes, any checkout that merely rewrote
  these files — a rebase, a branch switch — re-vendored and left an untracked
  copy behind.
- Changing the wire: run `dev/agent/mcp-verify.py`. Unit tests assert what the
  code declares, and that has been insufficient more than once.
- Adding a tool option or tag: `server/tests/mcp` asserts every tool carries a
  group and a risk, and `TOOLS.md` is generated from the same source.
