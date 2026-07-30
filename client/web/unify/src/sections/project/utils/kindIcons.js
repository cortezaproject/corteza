// Generates inline SVG data URIs per resource kind so echarts can render them as
// node symbols. Uses the *same* PrimeIcons artwork as the metrics strip (the
// system icon for each resource kind), drawn white on the kind-colored square.

import { kindConfig } from '@/sections/project/config/kinds'
import boltSvg from 'primeicons/raw-svg/bolt.svg?raw'
import sparklesSvg from 'primeicons/raw-svg/sparkles.svg?raw'
import chartBarSvg from 'primeicons/raw-svg/chart-bar.svg?raw'
import circleSvg from 'primeicons/raw-svg/circle.svg?raw'
import commentsSvg from 'primeicons/raw-svg/comments.svg?raw'
import databaseSvg from 'primeicons/raw-svg/database.svg?raw'
import folderSvg from 'primeicons/raw-svg/folder.svg?raw'
import idCardSvg from 'primeicons/raw-svg/id-card.svg?raw'
import linkSvg from 'primeicons/raw-svg/link.svg?raw'
import userSvg from 'primeicons/raw-svg/user.svg?raw'
import windowMaximizeSvg from 'primeicons/raw-svg/window-maximize.svg?raw'

// Keyed by the PrimeIcon class (KIND_CONFIG[kind].icon), so the chart stays in
// sync with whatever icon the metrics use.
const RAW = {
  'pi-database': databaseSvg,
  'pi-link': linkSvg,
  'pi-bolt': boltSvg,
  'pi-sparkles': sparklesSvg,
  'pi-comments': commentsSvg,
  'pi-window-maximize': windowMaximizeSvg,
  'pi-chart-bar': chartBarSvg,
  'pi-id-card': idCardSvg,
  'pi-user': userSvg,
  'pi-folder': folderSvg,
  'pi-circle': circleSvg,
}

// Severity colours for the reference badge. A node whose configuration points
// at something that isn't there is an error (red); one whose target is only
// decided at run time is a caveat (amber). Literal hex, not theme tokens: the
// symbol is rasterised into a canvas image, where a CSS variable means nothing.
const BADGE_COLOR = { missing: '#dc2626', warning: '#f59e0b' }

// A tile with the kind colour drained out of it, for a resource that is
// referenced from this project but defined outside it: dashed outline, washed
// fill, glyph in the kind colour rather than white. It reads as an outline of
// the kind rather than one of the project's own things.
const externalTile = color =>
  `<rect x="1.5" y="1.5" width="29" height="29" rx="7" fill="${color}" fill-opacity="0.16"` +
  ` stroke="${color}" stroke-width="2" stroke-dasharray="4 3"/>`

// Exclamation mark in a filled disc, top-right, with a white gap around it so
// it survives on top of the tile in either theme.
const badgeMark = severity => {
  const color = BADGE_COLOR[severity]
  if (!color) return ''
  return (
    `<circle cx="25.5" cy="6.5" r="6.5" fill="#ffffff"/>` +
    `<circle cx="25.5" cy="6.5" r="5.2" fill="${color}"/>` +
    `<rect x="24.75" y="3.2" width="1.5" height="4" rx="0.75" fill="#ffffff"/>` +
    `<circle cx="25.5" cy="9" r="0.9" fill="#ffffff"/>`
  )
}

const buildSvg = (kind, { external, badge }) => {
  const cfg = kindConfig(kind)
  const name = (cfg.icon || '').split(' ').pop() // 'pi pi-database' → 'pi-database'
  const raw = RAW[name] || RAW['pi-circle']
  // Embed the PrimeIcon as a nested <svg>, recolored (paths inherit fill).
  const glyph = raw.replace(
    /<svg\b[^>]*>/,
    `<svg x="7" y="7" width="18" height="18" viewBox="0 0 24 24" fill="${
      external ? cfg.stroke : '#ffffff'
    }">`,
  )
  const tile = external
    ? externalTile(cfg.stroke)
    : `<rect x="0" y="0" width="32" height="32" rx="8" fill="${cfg.stroke}"/>`
  return `
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="32" height="32">
      ${tile}
      ${glyph}
      ${badgeMark(badge)}
    </svg>
  `
    .replace(/\s+/g, ' ')
    .trim()
}

const cache = new Map()

// `badge` is 'missing' | 'warning' | undefined and `external` a boolean; both
// are part of the cache key, so a kind can have several tiles at once.
export function kindIconDataUri(kind, { external = false, badge = '' } = {}) {
  const key = `${kind}|${external ? 'x' : ''}|${badge}`
  if (cache.has(key)) return cache.get(key)
  const uri = `image://data:image/svg+xml;utf8,${encodeURIComponent(
    buildSvg(kind, { external, badge }),
  )}`
  cache.set(key, uri)
  return uri
}
