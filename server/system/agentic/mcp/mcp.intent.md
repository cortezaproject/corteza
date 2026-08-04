---
kind: folder
covers: recursive
owner: be
depends-on:
  - server/system/agentic/runtime
  - server/system/agentic/policy
  - server/system/agentic/toolkit
touched-by:
  - server/compose/agentic
  - server/automation/agentic
  - server/system/agentic
  - server/app/boot_levels.go
tests:
  - server/tests/mcp
---

# server/system/agentic/mcp

## Responsibility

The single registry of tools Human exposes to LLM agents, and the HTTP
(streamable MCP) transport in front of it.

Everything an agent can do to Human passes through here. Tool *implementations*
live beside the services they call — `compose/agentic`, `automation/agentic`,
`system/agentic` — and register into this package at boot. This package owns the
registry, the tagging vocabulary, and what a given caller is allowed to see.

`CONVENTIONS.md` alongside this file is the authoring spec every tool is written
against; `RESOURCES.md` fixes resource naming; `TOOLS.md` is generated coverage.
This document is the contract, those are the procedure.

## Key types & services

- `Registry` — name → tool, handler, tags. Populated once at boot from
  `app/boot_levels.go`. Duplicate registration panics: two handlers claiming one
  name previously produced a green build with one tool silently missing.
- `Group` / `Risk` (`tool_options.go`) — ride on `mcp.Tool.Meta` via
  `InGroup`/`WithRisk`. `WithRisk` is the **sole** writer of the protocol's four
  annotation hints; authors never set them by hand.
- `Scope` (`scope.go`) — per-request narrowing resolved from the URL.
- `disclosure` (`disclosure.go`) — per-session set of pulled-in tools.
- `MCPServer` — wraps mcp-go, applies the filter, the risk ceiling and the two
  meta-tools.

## Two consumers, one set of handlers

| Surface | Entry | Consumer |
|---|---|---|
| In-process | `Registry.GetTools` / `ExecuteTool` | Human's own agentic runtime |
| HTTP | `/api/mcp` | External MCP clients, primarily Claude Code |

They share handler functions. Where they legitimately differ is recorded in
`CONVENTIONS.md` §2.3; where they differed by accident, that was a defect.

## RBAC / permissions surface

**This package performs no authorization of its own, by design.** Handlers pass
the caller's context to the service layer and RBAC decides. Two consequences
that are easy to get wrong:

- **Group is presentation, not a boundary.** It shrinks what `tools/list`
  returns. A caller naming a tool outside its group has done nothing RBAC would
  not already permit.
- **Risk is a ceiling and *is* enforced** at dispatch, but it is self-selected
  via the URL — a seatbelt against accidents, not a lock against a hostile
  caller.

The premise has documented exceptions (`CONVENTIONS.md` §2.2), and §8.6 gives
tool authors the rule for a service that authorizes nothing.

`/api/mcp` needs **both** `HttpTokenValidator` and `HttpAuthenticatedOnly`: the
former checks a present token's scope but does not require one, and tool
discovery hands back every description before RBAC runs.

## Store / DAL usage

None. This package never touches the store. Handlers reach services, services
reach the store.

## When changing this

- Adding a tool: read `CONVENTIONS.md` first. Every tool needs a group and a
  risk; the structural test in `server/tests/mcp` fails without them.
- Adding a handler: wire it in `app/boot_levels.go` **and** in `buildRegistry`
  in the structural test, or its tools are invisible to every assertion.
- Changing what a session sees: `scope.go` for group and risk, `disclosure.go`
  for progressive disclosure. The default listing is five tools; a session
  searches for the rest.
- Changing anything on the wire: run `dev/agent/mcp-verify.py`. Unit tests
  assert what the code declares, and that has already been insufficient twice —
  the authentication fix passed its tests and did nothing, and the paging cursor
  was declared correctly and unusable.
