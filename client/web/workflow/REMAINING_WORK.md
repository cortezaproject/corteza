# Workflow Migration — Remaining Work

Items needed to reach full feature parity with the original mxGraph editor, cross-referenced against `MIGRATION_PLAN.md`.

---

## ✅ 1. Gateway Edge Renumbering on Delete (Plan line 204)

> "On excl gateway edge delete: renumber remaining edges."

**FIXED.** Added `renumberGatewayEdges(sourceId)` function that re-labels remaining exclusive gateway edges as `#1 - If`, `#2 - Else (if)`, etc. Called from both `sidebarDelete()` and `deleteSelected()`.

---

## ✅ 2. Sidebar `item.node.edges` Source/Target Style (Plan §Sidebar)

**FIXED.** `buildSidebarItem()` now includes `style` on both `source` and `target` in the edge map, matching what `Gateway.vue` (line 49) and other Configurators expect.

---

## ✅ 3. Edge Deletion from Sidebar (Standalone Edges)

**FIXED.** `sidebarDelete()` now detects `sidebar.value.itemType === 'edge'` and deletes the edge from `edges.value` instead of trying to match a node.

---

## ✅ 4. Sidebar `item.node.style`

Already correct — `buildSidebarItem()` sets `style: getStyleFromKind(data)?.style || data.kind || ''`.

---

## ✅ 5. Configurator Mutations Synced to Nodes (CRITICAL)

**FIXED.** Added deep watcher on `sidebar.value.item` that syncs mutations back into the corresponding VueFlow node's data. Synced fields: `arguments`, `results`, `ref`, `kind`, `defaultName`, `triggers`, `label`.

---

## ✅ 6. Trigger `stepID` Persistence

**FIXED** (by #5). Now that sidebar mutations are synced back to node data, `encodeWorkflow()` reads the correct `triggers.stepID` from the node.

---

## ✅ 7. Visual Node z-index (Plan line 175)

**FIXED.** Visual nodes now get `zIndex: -1` in both `decodeWorkflow()` and `useWorkflowDnD()`.

---

## ✅ 8. Visual Node "Not Connectable" (Plan line 175)

**FIXED.** Added `isValidConnection()` that rejects connections to/from visual nodes. Also set `connectable: false` on visual nodes created via DnD.

---

## ⚠️ 9. `parentNode` Nesting for Swimlanes

**Status: Implemented but untested.** `decodeWorkflow` sets `parentNode: vis.parent` when parent !== '1'. Needs testing with actual swimlane workflows to verify child nodes render inside their parent.

---

## ✅ 10. Edge Click → Select + Delete

**FIXED.** `deleteSelected()` now handles both selected nodes AND selected edges via the Delete/Backspace key.

---

## ⚠️ 11. Undo/Redo Loses VueFlow Internal State

**Status: Known limitation.** `useWorkflowHistory` uses JSON snapshots which lose VueFlow's internal computed dimensions. After undo, nodes may need a micro-adjustment. This is a low-severity issue that can be improved later.

---

## ✅ 12. `checkExistingTriggerPaths()` Not Called

**FIXED.** Now called via `nextTick()` at the end of `render()`.

---

## Summary Table

| # | Issue | Status |
|---|---|---|
| 1 | Gateway edge renumbering on delete | ✅ Fixed |
| 2 | `item.node.edges[].source.style` missing | ✅ Fixed |
| 3 | Edge deletion from sidebar broken | ✅ Fixed |
| 4 | `item.node.style` correctness | ✅ Already correct |
| 5 | **Configurator mutations not synced to nodes** | ✅ Fixed |
| 6 | Trigger stepID encode correctness | ✅ Fixed (via #5) |
| 7 | Visual node z-index | ✅ Fixed |
| 8 | Visual node not connectable | ✅ Fixed |
| 9 | Swimlane/parentNode nesting | ⚠️ Untested |
| 10 | Edge delete via keyboard | ✅ Fixed |
| 11 | Undo/redo position drift | ⚠️ Known limitation |
| 12 | `checkExistingTriggerPaths` not called | ✅ Fixed |
