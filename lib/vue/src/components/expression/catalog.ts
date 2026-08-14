// What an author may reference inside an expression input.
//
// `${...}` templates are evaluated client-side, by building a function over
// live JS objects (lib/js `interpolateTemplate`). The authority for what a
// name resolves to is therefore the shape of those objects — a `compose.Record`
// instance and the authenticated user — and not the server's expression types,
// which describe the separate, server-evaluated expression language.

export interface ScopeEntry {
  // Identifier as written in the expression.
  name: string
  // Human-facing kind, shown beside the completion.
  type: string
  label?: string
  fields?: ScopeEntry[]
  // Reachable at runtime but kept out of the suggestion list. Lint still
  // accepts these, so a hand-written reference to one is not reported wrong.
  suggest?: boolean
}

// Minimal structural view of a compose module — enough to read its fields
// without depending on the class.
export interface ScopeModule {
  moduleID?: string
  name?: string
  fields?: Array<{ name: string; label?: string; kind?: string }>
}

export interface ScopeSource {
  // Module of the record the `${record...}` variables resolve against: the
  // module of the page the expression renders on. On a record-list block this
  // is the *page's* module, not the module the block queries.
  recordModule?: ScopeModule | null
  // Module whose fields a bare identifier refers to — the module being
  // queried. Used for QL field completion only.
  queryModule?: ScopeModule | null
  // Whether a record is in scope at all. False on a non-record page, where
  // `${record...}`, `${recordID}` and `${ownerID}` resolve to nothing.
  hasRecord?: boolean
}

// Properties of a compose.Record reachable from a template. `values` is
// handled separately — its members come from the module, not the class.
const RECORD_PROPERTIES: ScopeEntry[] = [
  { name: 'recordID', type: 'ID' },
  { name: 'moduleID', type: 'ID' },
  { name: 'namespaceID', type: 'ID' },
  { name: 'ownedBy', type: 'ID' },
  { name: 'createdAt', type: 'DateTime' },
  { name: 'createdBy', type: 'ID' },
  { name: 'updatedAt', type: 'DateTime' },
  { name: 'updatedBy', type: 'ID' },
  { name: 'deletedAt', type: 'DateTime' },
  { name: 'deletedBy', type: 'ID' },
  { name: 'revision', type: 'Number', suggest: false },
  { name: 'meta', type: 'Object', suggest: false },
  { name: 'valueErrors', type: 'Object', suggest: false },
]

const USER_PROPERTIES: ScopeEntry[] = [
  { name: 'userID', type: 'ID' },
  { name: 'email', type: 'String' },
  { name: 'name', type: 'String' },
  { name: 'username', type: 'String' },
  { name: 'handle', type: 'String' },
  { name: 'emailConfirmed', type: 'Bool' },
  { name: 'userGroupID', type: 'ID', suggest: false },
  { name: 'labels', type: 'Object', suggest: false },
  { name: 'meta', type: 'Object', suggest: false },
]

function moduleValueEntries(module?: ScopeModule | null): ScopeEntry[] {
  return (module?.fields || []).map(f => ({
    name: f.name,
    label: f.label,
    type: f.kind || 'String',
  }))
}

// The roots an expression may open with, given what the surrounding page and
// block provide.
export function buildScope({ recordModule, hasRecord = true }: ScopeSource = {}): ScopeEntry[] {
  const entries: ScopeEntry[] = []

  if (hasRecord) {
    entries.push(
      { name: 'recordID', type: 'ID' },
      { name: 'ownerID', type: 'ID' },
      {
        name: 'record',
        type: 'Record',
        fields: [
          ...RECORD_PROPERTIES,
          // Left without `fields` when the module is unknown, which lint reads
          // as "cannot be checked" rather than "has no members".
          recordModule
            ? { name: 'values', type: 'Object', fields: moduleValueEntries(recordModule) }
            : { name: 'values', type: 'Object' },
        ],
      },
    )
  }

  entries.push(
    { name: 'userID', type: 'ID' },
    { name: 'user', type: 'User', fields: USER_PROPERTIES },
  )

  return entries
}

// Walks a dotted path against the scope.
//
// The third outcome matters as much as the other two: a path that runs into a
// member with no declared `fields` is *unverifiable*, not wrong, and lint must
// leave it alone.
export type PathResult =
  | { status: 'found'; entry: ScopeEntry }
  | { status: 'unknown'; parent: ScopeEntry | null; segment: string }
  | { status: 'unverifiable' }

export function resolvePath(scope: ScopeEntry[], path: string[]): PathResult {
  let candidates = scope
  let parent: ScopeEntry | null = null
  let entry: ScopeEntry | null = null

  for (const segment of path) {
    if (!candidates.length) return { status: 'unverifiable' }

    const hit = candidates.find(e => e.name === segment)
    if (!hit) return { status: 'unknown', parent, segment }

    parent = hit
    entry = hit
    candidates = hit.fields || []
  }

  return entry ? { status: 'found', entry } : { status: 'unverifiable' }
}

// Members offered after `prefix.` — the suggestion list, so entries marked
// `suggest: false` are withheld.
export function membersAt(scope: ScopeEntry[], path: string[]): ScopeEntry[] {
  if (!path.length) return scope.filter(e => e.suggest !== false)

  const res = resolvePath(scope, path)
  if (res.status !== 'found') return []

  return (res.entry.fields || []).filter(e => e.suggest !== false)
}
