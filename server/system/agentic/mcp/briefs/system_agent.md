# Brief: `system_agent`

Tool file: `server/system/agentic/agent_tools.go`
Handler:   `server/system/agentic/agent_handler.go`
Group: `configuring`

## 1. Service

`system/service/agent.go` + `agent.gen.go` — `sysService.DefaultAgent`.

| Method | Signature | Becomes |
|---|---|---|
| `FindByID` | `(ctx, ID uint64) (*types.Agent, error)` | single-fetch in `lookup` |
| `Search` | `(ctx, types.AgentFilter) (AgentSet, AgentFilter, error)` | list mode, and handle resolution |
| `Create` / `Update` | `(ctx, *types.Agent) (*types.Agent, error)` | `create` / `update` |
| `DeleteByID` | `(ctx, ID uint64) error` | `delete` |
| `UndeleteByID` | `(ctx, ID uint64) error` | `undelete` — hand-written, `agent.go:145` |
| `Get` | alias for `FindByID` | not a tool |

Risks: `lookup` read; `create`/`update`/`undelete` write; `delete` destructive.

## 2. Authorization (§8.6) — PASSES

`ac.Can*` throughout: `CanReadAgent` (`agent.go:66`, and again in `onSearch`'s
`filter.Check` at `:173`), `CanUpdateAgent` (`:114`), `CanDeleteAgent` (`:152`).
`CanCreateAgent` is in the generated `Create` scaffold. `onSearch` also narrows
per row, so a listing cannot leak an agent the caller may not read.

Note the delete/undelete pairing is non-standard: `UndeleteByID` checks
`CanDeleteAgent`, not a create/update permission. That is deliberate and is why
the generated undelete is disabled.

## 3. Identifier strategy

**No `FindByAny`.** Resolve a ref by trying, in order:

1. all-digits → `FindByID`
2. otherwise → `Search` with `AgentFilter{Handle: ref}`, exactly one hit

An agent has **no `Name` field** — the human-readable label is `Meta.Short`,
and `Meta.Description` is the long one. Do not resolve by `Short`: it is not
unique and nothing indexes it.

`undelete` takes a required `agentID`. `AgentFilter.Deleted` defaults to
excluded, so a deleted agent is not findable by handle.

## 4. Filter fields → params

`AgentFilter` (`system/types/agent.go`):

| Field | Param? | Note |
|---|---|---|
| `Query` | yes | free text |
| `Handle` | yes | |
| `Status` | yes | `active` is the create default; the field is free-form |
| `AgentID` | no | the `agent` ref covers it |
| `ProjectID` | no | agents made through MCP are not project-scoped; see traps |
| `Deleted` | yes, as `includeDeleted` | |
| `Labels` | no | no label tools exist yet |

## 5. Undelete

`UndeleteByID` exists (hand-written). Ships as `system_agent_undelete`, risk
`write`, taking `agentID` — say in the description that a handle will not work.

## 6. Domain ops

None. Status is a plain field on update, not an op: nothing in the service
treats it specially beyond defaulting it to `active` at create.

## 7. Traps

**Update replaces the whole record.** `onUpdate` calls `store.UpdateAgent` with
the payload it is given — there is no field-level merge. A partial update tool
MUST load the agent first and apply only what the caller sent, or every section
the caller omitted is wiped.

**Five nested config sections**, each a struct: `meta`, `behavior`,
`execution`, `access`, `invocation`. Per §8.2 each is one JSON-object param
that replaces that whole section when present and leaves it alone when absent.
Flattening them would be ~25 params and would still not reach
`access.tools[].allow[].moduleIDs`.

**`Revision` is server-owned** — bumped on every update. Never a param.

**`ProjectID` is preserved on update** and never in the payload; an update that
sent a zero would unlink the agent from its project's resource graph, which the
service guards against explicitly.

**The name is the only validated field.** `onCreate`/`onUpdate` refuse an
empty or blank `meta.short` with `missing name`, and check nothing else:
provider, model and temperature are stored as given and first meet the provider
on `system_agent_exec`. The webapp's agent editor requires provider and model
before it saves; the API does not.

**`access.tools` is an allow-list of MCP tool names** — the same names this
registry publishes. An agent editing agents is expressible, so the description
must say plainly that granting `system_agent_update` lets that agent re-grant
itself anything. RBAC still applies to the agent's own identity; the allow-list
narrows, it does not widen.
