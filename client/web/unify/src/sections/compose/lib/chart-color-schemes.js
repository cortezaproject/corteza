// A chart's palette is either one of the tables compiled into lib/js or one an
// admin defined for this instance. The instance's own live in a single global
// setting, written whole; this module is the one place that knows its name, its
// shape, and how the two sources are put in front of a chart author.

export const COLOR_SCHEMES_SETTING = 'ui.charts.colorSchemes'

/** The instance's own schemes, or an empty list where settings never loaded. */
export function readColorSchemes($Settings) {
  const value = $Settings?.get(COLOR_SCHEMES_SETTING, [])
  return Array.isArray(value) ? value : []
}

// getColorschemeColors picks a custom scheme out by this substring, so an id
// without it is looked up in the built-in tables and resolves to nothing.
export const isCustomScheme = id => !!id && String(id).includes('custom')

export const newCustomScheme = (now = Date.now()) => ({
  id: `custom-${now}`,
  name: '',
  colors: ['#6C757D', '#000000'],
})

const capitalize = w => (w ? `${w[0].toUpperCase()}${w.slice(1)}` : w)

// A built-in key carries its swatch count as a suffix ("ClassicOrangeBlue13");
// the picker shows the two apart.
function splitCount(key) {
  const m = /(\D+)(\d+)$/.exec(key)
  return { label: m?.[1] || key, count: m?.[2] || '' }
}

export function builtinColorSchemes(colorschemes = {}, countLabel = c => `${c} colors`) {
  const out = []

  for (const family in colorschemes) {
    for (const key in colorschemes[family]) {
      const { label, count } = splitCount(key)
      out.push({
        id: `${family}.${key}`,
        name: `${capitalize(family)}: ${capitalize(label)} (${countLabel(count)})`,
        colors: [...colorschemes[family][key]],
      })
    }
  }

  return out
}

// The instance's own schemes come first: there is a handful of them among some
// hundreds of built-ins, and they are the ones it actually uses.
export function colorSchemeOptions(custom, colorschemes, countLabel) {
  const own = (custom || [])
    .filter(s => s && s.id)
    .map(({ id, name, colors }) => ({ id, name: name || id, colors: [...(colors || [])] }))

  return [...own, ...builtinColorSchemes(colorschemes, countLabel)]
}

// Both writes rebuild the array rather than mutating the one the app renders
// from: the setting is stored whole, and a half-applied edit that then fails to
// save leaves the picker showing something the server never took.
export function upsertColorScheme(list, scheme) {
  const next = (list || []).map(s => ({ ...s }))
  const entry = {
    id: scheme.id,
    name: (scheme.name || '').trim(),
    colors: [...(scheme.colors || [])],
  }

  const at = next.findIndex(({ id }) => id === entry.id)
  if (at < 0) {
    next.push(entry)
  } else {
    next.splice(at, 1, entry)
  }

  return next
}

export function removeColorScheme(list, id) {
  return (list || []).filter(s => s.id !== id)
}
