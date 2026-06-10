// Generates inline SVG data URIs per resource kind so echarts can render them as
// node symbols. Uses the *same* PrimeIcons artwork as the metrics strip (the
// system icon for each resource kind), drawn white on the kind-colored square.

import { kindConfig } from '@/sections/project-poc/config/kinds'
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

const buildSvg = kind => {
  const cfg = kindConfig(kind)
  const name = (cfg.icon || '').split(' ').pop() // 'pi pi-database' → 'pi-database'
  const raw = RAW[name] || RAW['pi-circle']
  // Embed the PrimeIcon as a nested <svg>, recolored white (paths inherit fill).
  const glyph = raw.replace(
    /<svg\b[^>]*>/,
    '<svg x="7" y="7" width="18" height="18" viewBox="0 0 24 24" fill="#ffffff">',
  )
  return `
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="32" height="32">
      <rect x="0" y="0" width="32" height="32" rx="8" fill="${cfg.stroke}"/>
      ${glyph}
    </svg>
  `
    .replace(/\s+/g, ' ')
    .trim()
}

const cache = new Map()

export function kindIconDataUri(kind) {
  if (cache.has(kind)) return cache.get(kind)
  const uri = `image://data:image/svg+xml;utf8,${encodeURIComponent(buildSvg(kind))}`
  cache.set(kind, uri)
  return uri
}
