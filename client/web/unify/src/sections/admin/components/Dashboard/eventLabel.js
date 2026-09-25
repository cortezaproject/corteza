// Readable names for action-log entries. A stored resource looks like
// "corteza::compose:namespace/123" or "system:user"; an action is a camelCase
// verb such as "execAndWait". Known ones are translated, the rest are
// spelled out from their identifiers so nothing shows raw.

// "corteza::compose:namespace/123" → { type: 'compose:namespace', id: '123' }
export function parseResource(resource) {
  const stripped = String(resource || '').replace(/^corteza::/, '')
  const slash = stripped.indexOf('/')
  const type = slash === -1 ? stripped : stripped.slice(0, slash)
  const id = slash === -1 ? '' : stripped.slice(slash + 1)
  return { type, id: id === '*' ? '' : id }
}

// "system:user-group" → "user group"; "compose:page-layout" → "page layout"
export function humanizeType(type) {
  const name = type.includes(':') ? type.slice(type.indexOf(':') + 1) : type
  return name.replace(/[-_]+/g, ' ').trim()
}

// "execAndWait" → "exec and wait"; "markAllAsRead" → "mark all as read"
export function humanizeAction(action) {
  return String(action || '')
    .replace(/([a-z0-9])([A-Z])/g, '$1 $2')
    .replace(/[-_]+/g, ' ')
    .toLowerCase()
    .trim()
}

function keyOf(s) {
  return s.replace(/[^a-zA-Z0-9]+/g, '_')
}

function capitalize(s) {
  return s ? s[0].toUpperCase() + s.slice(1) : s
}

// Builds the label with the app's translator: t/te from useI18n.
export function describeEvent({ t, te }, resource, action) {
  const { type } = parseResource(resource)

  const typeKey = `dashboard.events.resources.${keyOf(type)}`
  const noun = te(typeKey) ? t(typeKey) : humanizeType(type)

  const actionKey = `dashboard.events.actions.${keyOf(action)}`
  const verb = te(actionKey) ? t(actionKey) : humanizeAction(action)

  return capitalize(`${noun} ${verb}`.trim())
}
