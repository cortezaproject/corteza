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

// Permission flags a serialized record carries. `serialize()` spreads the whole
// instance, so these reach the evaluator alongside the stored fields and are
// the usual way a visibility rule asks "may this user edit?".
const RECORD_PERMISSIONS: ScopeEntry[] = [
  { name: 'canUpdateRecord', type: 'Bool' },
  { name: 'canDeleteRecord', type: 'Bool' },
  { name: 'canReadRecord', type: 'Bool' },
  { name: 'canUndeleteRecord', type: 'Bool', suggest: false },
  { name: 'canManageOwnerOnRecord', type: 'Bool', suggest: false },
  { name: 'canSearchRevision', type: 'Bool', suggest: false },
  { name: 'canGrant', type: 'Bool', suggest: false },
]

const SCREEN_PROPERTIES: ScopeEntry[] = [
  { name: 'width', type: 'Number' },
  { name: 'height', type: 'Number' },
  { name: 'breakpoint', type: 'String', label: 'xxs | xs | sm | md | lg' },
  { name: 'userAgent', type: 'String', suggest: false },
]

// Variables for the server-evaluated expression language — block and layout
// visibility, and record field conditions.
//
// A different scope from buildScope(): these are the keys the webapp puts in
// the `variables` payload of POST /system/expressions/evaluate (see
// usePageVisibility and RecordBlock), not the bindings a `${...}` template is
// built over. There is no bare `recordID`/`ownerID`/`userID` here, and `record`
// arrives serialized, which is what brings the permission flags with it.
export function buildExprScope({ recordModule, hasRecord = true }: ScopeSource = {}): ScopeEntry[] {
  const entries: ScopeEntry[] = []

  if (hasRecord) {
    entries.push({
      name: 'record',
      type: 'Record',
      fields: [
        ...RECORD_PROPERTIES,
        recordModule
          ? { name: 'values', type: 'Object', fields: moduleValueEntries(recordModule) }
          : { name: 'values', type: 'Object' },
        ...RECORD_PERMISSIONS,
      ],
    })
  }

  entries.push(
    { name: 'user', type: 'User', fields: USER_PROPERTIES },
    { name: 'screen', type: 'Screen', fields: SCREEN_PROPERTIES },
  )

  if (hasRecord) {
    entries.push(
      { name: 'isView', type: 'Bool' },
      { name: 'isCreate', type: 'Bool' },
      { name: 'isEdit', type: 'Bool' },
    )
  }

  return entries
}

// Record shape as the *server* builds it for a field expression
// (types.Record.Dict) — close to the serialized JS record but not identical: it
// carries `ID` and `createdByAgent`, and none of the `can…` permission flags.
const RECORD_DICT: ScopeEntry[] = [
  ...RECORD_PROPERTIES,
  { name: 'ID', type: 'ID', suggest: false },
  { name: 'createdByAgent', type: 'ID', suggest: false },
]

function recordDict(module?: ScopeModule | null): ScopeEntry {
  return {
    name: 'record',
    type: 'Record',
    fields: [
      ...RECORD_DICT,
      module
        ? { name: 'values', type: 'Object', fields: moduleValueEntries(module) }
        : { name: 'values', type: 'Object' },
    ],
  }
}

// Which field expression is being written. Each gets a different scope from the
// server, so they cannot share one:
//
//   sanitizer  value                                    (values/sanitizer.go)
//   validator  value, oldValue, values.<field>          (values/validator.go)
//   value      <field> at the top level, new, old       (values/expr.go)
export type FieldExprKind = 'sanitizer' | 'validator' | 'value'

export function buildFieldExprScope(
  kind: FieldExprKind,
  module?: ScopeModule | null,
): ScopeEntry[] {
  const values: ScopeEntry = module
    ? { name: 'values', type: 'Object', fields: moduleValueEntries(module) }
    : { name: 'values', type: 'Object' }

  if (kind === 'sanitizer') {
    return [{ name: 'value', type: 'Any', label: 'the value being sanitized' }]
  }

  if (kind === 'validator') {
    return [
      { name: 'value', type: 'Any', label: 'the value being validated' },
      { name: 'oldValue', type: 'Any', label: 'its value before this change' },
      values,
    ]
  }

  // A value expression reads the record's own fields as bare names.
  return [
    ...moduleValueEntries(module),
    { ...recordDict(module), name: 'new', label: 'the record being saved' },
    { ...recordDict(module), name: 'old', label: 'the record before this save' },
  ]
}

// A workflow expression sees the variables its trigger puts in scope. The
// server describes those per event type (GET /automations/event-types →
// `properties`), and the dry-run form already reads the same list, so this
// takes it as given rather than restating it.
//
// A workflow may carry several triggers; pass the union. Types come through as
// the expr type names the automation registry uses (`ComposeRecord`, `User`…),
// which is what the reference panel shows too.
export interface TriggerProperty {
  name: string
  types?: string[]
}

export function buildWorkflowScope(properties: TriggerProperty[] = []): ScopeEntry[] {
  const seen = new Map<string, ScopeEntry>()

  for (const p of properties) {
    if (!p?.name || seen.has(p.name)) continue
    seen.set(p.name, { name: p.name, type: p.types?.[0] || 'Any' })
  }

  return [...seen.values()].sort((a, b) => a.name.localeCompare(b.name))
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
