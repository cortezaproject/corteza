// Translation between this section's resource KINDS (config/kinds.js — the
// vocabulary the graph, steps and dialogs speak) and the RBAC/envoy RESOURCE
// REFERENCE strings the backend persists.
//
// This exists because AI system membership stores refs, not kinds+ids: a
// ProjectAiSystemEntry is `(projectAiSystemID, resourceRef)` where resourceRef
// is the SAME string RBAC rules and envoy already use, e.g.
// `corteza::compose:module/2145…`. Reusing that vocabulary rather than
// inventing a `kind:id` one was a deliberate ruling (2026-07-30) — it keeps
// refs parseable by machinery that already exists, and interoperable with
// anything RBAC-shaped later.
//
// THESE STRINGS ARE A PERSISTED INTERFACE, exactly like config/friaTaxonomies
// keys: once an AI system stores a ref, that string is in the database. The
// resource-type halves must track the backend's own `<X>ResourceType`
// constants (server/*/types/resources.gen.go) — they are not ours to rename.

// Kind → backend resource type. Only kinds that can belong to an AI system
// appear; a kind absent here simply cannot be added as a member.
//
// Roles are included because an AI system also records the roles responsible
// for OVERSEEING it (Art. 14 human oversight, Art. 26(2) — the deployer
// assigns oversight to competent staff). They ride in the same entry table,
// distinguished by the ref's own kind rather than a separate column.
export const RESOURCE_TYPE_BY_KIND = {
  module: 'corteza::compose:module',
  page: 'corteza::compose:page',
  chart: 'corteza::compose:chart',
  automation: 'corteza::automation:ng-automation',
  workflow: 'corteza::automation:workflow',
  agent: 'corteza::system:agent',
  chatbot: 'corteza::system:chatbot',
  connection: 'corteza::system:configured-connection',
  role: 'corteza::system:role',
}

export const KIND_BY_RESOURCE_TYPE = Object.fromEntries(
  Object.entries(RESOURCE_TYPE_BY_KIND).map(([kind, type]) => [type, kind]),
)

// The project-graph endpoint names a reference's TARGET by resource type, and
// its vocabulary is a little wider than the AI-system one above: a connection
// reference can name either the project's own configured connection or the
// tenant-level DAL connection behind it, and both draw as a "connection" node.
// The reverse map above can't carry that (two types, one kind), so the graph
// direction gets its own table.
//
// Only the types the graph can emit appear — server/system/service/
// project_graph.go drops refs to any other type before the payload is built.
export const GRAPH_KIND_BY_RESOURCE_TYPE = {
  ...KIND_BY_RESOURCE_TYPE,
  'corteza::system:dal-connection': 'connection',
}

// Kinds whose membership means "this resource is PART OF the AI system", as
// opposed to the oversight relationship roles carry. Kept explicit so the
// editor can present the two as separate concerns without re-deriving the
// distinction from a hardcoded list of exceptions.
export const MEMBER_KINDS = [
  'module',
  'page',
  'chart',
  'automation',
  'agent',
  'chatbot',
  'connection',
]
export const OVERSIGHT_KINDS = ['role']

export function buildResourceRef(kind, id) {
  const type = RESOURCE_TYPE_BY_KIND[kind]
  if (!type || !id) return null
  return `${type}/${id}`
}

// Splits a stored ref back into { kind, id }. Returns nulls for a ref whose
// resource type is unknown to this frontend rather than throwing: refs
// deliberately outlive the resources they point at (a deleted member is kept and
// shown as a tombstone), and a ref may name a type this build has never heard
// of. Callers render what they can and say so otherwise.
export function parseResourceRef(ref) {
  if (typeof ref !== 'string') return { kind: null, id: null, type: null }

  const slash = ref.lastIndexOf('/')
  if (slash < 0) return { kind: null, id: null, type: ref }

  const type = ref.slice(0, slash)
  return { kind: KIND_BY_RESOURCE_TYPE[type] || null, id: ref.slice(slash + 1), type }
}
