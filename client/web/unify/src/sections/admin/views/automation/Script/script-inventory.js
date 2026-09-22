import { constraintChips } from '@planetcrust/human-vue'

/**
 * Inventory facts about a Corredor script, read off the shape
 * `automationList` returns.
 *
 * A script's name is its path inside the Corredor extension —
 * `/server-scripts/<extension>/<file>.js:<export>` for one that runs inside
 * Corredor, `/client-scripts/<bundle>/<extension>/<file>.js:<export>` for one
 * bundled into the browser. The name is the only place the kind, the bundle and
 * the extension are always filled in: the wire `type` and `bundle` fields are
 * omitted when the bundler leaves them empty.
 */

const clientScriptPrefix = '/client-scripts/'
const chipSeparator = ' · '

function segments(name) {
  return String(name || '')
    .split(':')[0]
    .split('/')
    .filter(Boolean)
}

/**
 * Where a script runs: `client` in the browser, `server` inside Corredor.
 */
export function scriptKind(script = {}) {
  return String(script?.name || '').startsWith(clientScriptPrefix) ? 'client' : 'server'
}

/**
 * The client bundle a script is built into, empty for a server script.
 */
export function scriptBundle(script = {}) {
  if (script?.bundle) return script.bundle
  if (scriptKind(script) !== 'client') return ''
  return segments(script?.name)[1] || ''
}

/**
 * The extension folder a script belongs to: the directory under the kind and,
 * for a client script, under its bundle. A script sitting directly there has
 * none.
 */
export function scriptExtension(script = {}) {
  const rest = segments(script?.name).slice(scriptKind(script) === 'client' ? 2 : 1)
  return rest.length > 1 ? rest[0] : ''
}

/**
 * One row per trigger: what fires it, and the constraints that narrow it.
 */
export function triggerRows(script = {}) {
  return (script?.triggers || []).map(trigger => ({
    label: [(trigger?.eventTypes || []).join(', '), (trigger?.resourceTypes || []).join(', ')]
      .filter(Boolean)
      .join(chipSeparator),
    constraints: constraintChips(trigger?.constraints),
  }))
}

/**
 * Scripts by the extension they came from, extensions alphabetically and the
 * loose scripts last.
 */
export function groupScripts(scripts) {
  const groups = new Map()

  ;(scripts || []).forEach(script => {
    const extension = scriptExtension(script)
    if (!groups.has(extension)) groups.set(extension, [])
    groups.get(extension).push(script)
  })

  return [...groups.entries()]
    .map(([extension, items]) => ({
      extension,
      items: [...items].sort((a, b) => String(a?.name).localeCompare(String(b?.name))),
    }))
    .sort((a, b) => {
      if (!a.extension) return 1
      if (!b.extension) return -1
      return a.extension.localeCompare(b.extension)
    })
}

/**
 * Whether a script of this kind passes the server/client toggle pair. The pair
 * narrows the list only while exactly one of the two is on: neither and both
 * mean "every kind".
 */
export function matchesKindFilter(kind, { server, client } = {}) {
  if (!!server === !!client) return true
  return server ? kind === 'server' : kind === 'client'
}

/**
 * What the list says about the Corredor connection, or null when the server
 * reported nothing about it.
 */
export function corredorBanner({ enabled, connected } = {}) {
  if (typeof enabled !== 'boolean') return null
  if (!enabled) return { state: 'disabled', severity: 'info' }
  if (connected === true) return { state: 'connected', severity: 'success' }
  if (connected === false) return { state: 'unreachable', severity: 'warn' }
  return null
}

/**
 * How long ago a timestamp was, in words. Coarse by design: the exact instant
 * is shown beside it.
 */
export function relativeTime(value, locale) {
  const then = value ? new Date(value) : null
  if (!then || Number.isNaN(then.getTime())) return ''

  const seconds = Math.round((then.getTime() - Date.now()) / 1000)
  const abs = Math.abs(seconds)
  const rtf = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' })

  if (abs < 60) return rtf.format(seconds, 'second')
  if (abs < 3600) return rtf.format(Math.round(seconds / 60), 'minute')
  if (abs < 86400) return rtf.format(Math.round(seconds / 3600), 'hour')
  return rtf.format(Math.round(seconds / 86400), 'day')
}
