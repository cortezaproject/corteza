# Project Graph Endpoint

## Purpose

Single endpoint that returns all resources scoped to a project as a graph payload (nodes + edges). Computed at query time from existing Corteza entity tables — no separate resource manifest or planning table.

---

## Endpoint

```
GET /api/system/projects/{projectID}/graph
```

Auth: standard JWT + tenant + project scope middleware (see middleware spec).
Permission: `canRead` capability on the project.

---

## Response shape

```json
{
  "nodes": [
    { "id": "uint64", "kind": "module|page|chart|connection|automation|agent|chatbot|role", "name": "string", "sensitivity": "string|null" }
  ],
  "edges": [
    { "sourceID": "uint64", "targetID": "uint64", "reason": "string" }
  ]
}
```

`reason` is a short machine label for why the edge exists. Used for tooltip/debug, not rendered as a label in the graph. Examples: `module-field-ref`, `page-module`, `trigger-module`, `workflow-agent`.

---

## Node sources

Query each entity type scoped to `(tenantID, projectID)`. All require `project_id` column to be present (depends on resource scoping work).

| Kind | Source table | Sensitivity source |
|---|---|---|
| `module` | `compose_modules` | `DalSensitivityLevel` ref on module |
| `page` | `compose_pages` | — |
| `chart` | `compose_charts` | — |
| `connection` | `configured_connections` | sensitivity field on connection |
| `automation` | `automation_workflows` | — |
| `agent` | `system_agents` | — |
| `chatbot` | `system_chatbots` | — |
| `role` | `system_roles` | — |

All queries filter: `project_id = :projectID AND tenant_id = :tenantID AND deleted_at IS NULL`.

---

## Edge derivation

Edges are derived from existing metadata relationships. No separate edge table.

| Edge reason | Source | Derived from |
|---|---|---|
| `module-field-ref` | module → module | `compose_module_fields` where `field_type = 'Record'` and `options.moduleID` points to another module in the same project |
| `page-module` | page → module | `compose_page_blocks` where block references a moduleID |
| `trigger-module` | automation → module | `automation_triggers` where resource ref is a compose module |
| `workflow-agent` | automation → agent | `automation_workflows` steps/config referencing an agent ID |
| `chatbot-agent` | chatbot → agent | `system_chatbots.config` agent reference |
| `chatbot-connection` | chatbot → connection | chatbot config connection ref |
| `agent-connection` | agent → connection | agent config connection ref |

All edge queries are scoped to the same `projectID` — edges to nodes outside the project are dropped.

---

## Service logic

```
ProjectGraphService.Graph(ctx, projectID) (*ProjectGraph, error)
```

1. Load project → verify access via scope middleware (already done before handler)
2. Fan out parallel queries for each node kind (8 queries)
3. Build node map: `id → node`
4. Fan out parallel queries for each edge source (7 queries)
5. For each raw edge row: check both source and target IDs exist in node map → keep or drop
6. Return `ProjectGraph{Nodes, Edges}`

No caching at this stage. Graph is small enough to compute per-request.

---

## Types (`server/system/types/project_graph.go`)

```go
type ProjectGraphNode struct {
    ID          uint64 `json:"id,string"`
    Kind        string `json:"kind"`
    Name        string `json:"name"`
    Sensitivity string `json:"sensitivity,omitempty"`
}

type ProjectGraphEdge struct {
    SourceID uint64 `json:"sourceID,string"`
    TargetID uint64 `json:"targetID,string"`
    Reason   string `json:"reason"`
}

type ProjectGraph struct {
    Nodes []*ProjectGraphNode `json:"nodes"`
    Edges []*ProjectGraphEdge `json:"edges"`
}
```

---

## Files to create / modify

| File | Action |
|---|---|
| `server/system/types/project_graph.go` | New — `ProjectGraph`, `ProjectGraphNode`, `ProjectGraphEdge` |
| `server/system/service/project_graph.go` | New — `ProjectGraphService` with `Graph(ctx, projectID)` |
| `server/system/rest/project.go` | Modify — add `Graph` handler |
| `server/system/rest.yaml` | Modify — add `GET /projects/{projectID}/graph` |
| `server/system/rest/request/project.go` | Modify (auto-gen) — add `ProjectGraph` request struct |

---

## Dependency

Requires `project_id` column on all source tables (resource scoping work). Endpoint cannot return meaningful data until scoping is in place. Can be stubbed to return empty nodes/edges before then.
