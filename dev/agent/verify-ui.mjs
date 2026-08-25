#!/usr/bin/env node
// Render-verify webapp paths in a real (headless) browser.
//
// Usage: node dev/agent/verify-ui.mjs [--out DIR] PATH [PATH...]
//   node dev/agent/verify-ui.mjs '/compose/namespace/agent-sandbox/pages/123'
//
// Logs in as agent@local.dev (password from .state/ui-password, created by
// bootstrap.sh), navigates each path on the vite dev server, captures console
// and page errors, and writes a full-page screenshot per path. Exits non-zero
// if any path produced page errors. Local-only.
//
// Why this exists: API checks cannot see rendering bugs (mis-sized grids,
// blocks that clip to nothing, unresolved refs showing raw IDs). Both
// toolkit field tests caught real bugs only through this kind of check.
//
// Each path prints a text report of the compose blocks it found — geometry,
// clipped content, empty blocks, raw IDs, uninterpolated ${...} templates —
// before naming the screenshot. Read the report first: it catches the common
// layout faults on its own, and a screenshot is only worth opening when the
// report flags something or the check is about visual design.

import { readFileSync, readdirSync, mkdirSync } from 'node:fs'
import { basename, dirname, join, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { createRequire } from 'node:module'
import { stack } from './stack.mjs'

const AGENT_DIR = dirname(fileURLToPath(import.meta.url))
const REPO_DIR = resolve(AGENT_DIR, '..', '..')
const { HUMAN_WEBAPP: WEBAPP, HUMAN_BASE: API_BASE } = stack()

for (const u of [WEBAPP, API_BASE]) {
  const host = new URL(u).hostname
  if (!['localhost', '127.0.0.1', '::1'].includes(host)) {
    console.error(`verify-ui is local-only; refusing ${u}`)
    process.exit(1)
  }
}

// resolve playwright-core out of the pnpm store without depending on version
const pnpmDir = join(REPO_DIR, 'node_modules', '.pnpm')
const pwDir = readdirSync(pnpmDir).find(d => d.startsWith('playwright-core@'))
if (!pwDir) {
  console.error('playwright-core not found in node_modules/.pnpm — run pnpm install')
  process.exit(1)
}
const require = createRequire(
  join(pnpmDir, pwDir, 'node_modules', 'playwright-core', 'package.json'),
)
const { chromium } = require(join(pnpmDir, pwDir, 'node_modules', 'playwright-core'))

const args = process.argv.slice(2)
const SESSION = process.env.CLAUDE_CODE_SESSION_ID || 'unknown'
let outDir = join(AGENT_DIR, '.state', 'sessions', SESSION, basename(REPO_DIR), 'ui')
const paths = []
for (let i = 0; i < args.length; i++) {
  if (args[i] === '--out') outDir = args[++i]
  else paths.push(args[i])
}
if (!paths.length) {
  console.error('usage: verify-ui.mjs [--out DIR] PATH [PATH...]')
  process.exit(1)
}
mkdirSync(outDir, { recursive: true })

const email = 'agent@local.dev'
const password = readFileSync(join(AGENT_DIR, '.state', 'ui-password'), 'utf8').trim()

const browser = await chromium.launch({ headless: true })
const page = await browser.newPage({ viewport: { width: 1440, height: 900 } })

const problems = []
page.on('pageerror', e => problems.push(`pageerror: ${e.message}`))
page.on('console', m => {
  // resource-load failures are reported (and filtered) via the response
  // listener below — the console duplicate has no URL and only adds noise
  if (m.type() === 'error' && !m.text().startsWith('Failed to load resource'))
    problems.push(`console.error: ${m.text().slice(0, 300)}`)
})
const httpFloor = process.env.VERBOSE_HTTP ? 400 : 500
page.on('response', r => {
  if (r.status() >= httpFloor) problems.push(`HTTP ${r.status()}: ${r.url()}`)
})

// login through the real auth flow
async function login() {
  await page.goto(WEBAPP, { waitUntil: 'domcontentloaded' })
  await page.waitForURL(/\/auth\//, { timeout: 20000 }).catch(() => {})
  if (page.url().includes('/auth/')) {
    await page.fill('input[name="email"]', email)
    const pw = page.locator('input[name="password"]')
    if (!(await pw.count())) await page.click('button[type="submit"]') // two-step login
    await page.fill('input[name="password"]', password)
    await page.click('button[type="submit"]')
    await page.waitForURL(u => !u.href.includes('/auth/'), { timeout: 20000 })
  }
  problems.length = 0 // ignore login-phase noise
}
await login()

let failed = false

// What the page actually rendered, read out of the live DOM.
//
// Compose lays blocks out with gridstack: .grid-stack-item carries the
// configured cell geometry in gs-*, and its .grid-stack-item-content is
// overflow:auto — so content taller than the block does not spill, it
// silently scrolls. That is exactly what a "clipped" block looks like, and
// it is invisible to every API check. A single-block page skips the grid.
const PAGE_REPORT = () => {
  // Two different faults, measured apart because only one is a defect.
  //
  // The grid cell (.grid-stack-item-content) is overflow:auto, so a block
  // shorter than its content does not spill — it silently scrolls, and the
  // last field of a Record block simply is not there. That is the fault.
  // A datatable viewport inside a list block scrolls too, but that is how a
  // list pages and it carries a paginator saying so; only its *horizontal*
  // overflow matters, because that is a column cut off the right edge.
  //
  // Only auto/scroll counts. PrimeVue icon buttons are overflow:hidden and
  // report a phantom ~21px, which is what makes a naive scan pure noise.
  const measure = root => {
    // Content the reader cannot reach, vs a table viewport paging normally.
    // The distinction is the element doing the scrolling, not how deep it is:
    // a datatable viewport carries a paginator that says there is more, so its
    // vertical scroll is by design. Anything else scrolling vertically — a
    // Record block taller than its cell — has simply lost its last fields.
    // Horizontal overflow is a cut-off column wherever it happens.
    const lost = { v: 0, h: 0 }
    const listScroll = { v: 0 }
    const isListViewport = el => /p-datatable|p-virtualscroller/.test(String(el.className))
    const consider = el => {
      const st = getComputedStyle(el)
      if (['auto', 'scroll'].includes(st.overflowY)) {
        const dv = el.scrollHeight - el.clientHeight
        if (isListViewport(el)) listScroll.v = Math.max(listScroll.v, dv)
        else lost.v = Math.max(lost.v, dv)
      }
      if (['auto', 'scroll'].includes(st.overflowX))
        lost.h = Math.max(lost.h, el.scrollWidth - el.clientWidth)
    }
    consider(root)
    for (const el of root.querySelectorAll('*')) {
      // PrimeVue icon buttons are overflow:hidden and report a phantom ~21px;
      // only real content areas are worth measuring
      if (el.clientWidth < 150 || el.clientHeight < 50) continue
      consider(el)
    }
    return { lost, listScroll }
  }

  const items = [...document.querySelectorAll('.grid-stack-item')]
  const boxes = items.length
    ? items.map(el => ({
        content: el.querySelector('.grid-stack-item-content'),
        xywh: ['x', 'y', 'w', 'h'].map(a => Number(el.getAttribute('gs-' + a))),
      }))
    : [...document.querySelectorAll('.single-block-wrapper')].map(el => ({
        content: el,
        xywh: null,
      }))

  // Ink on a block's canvases. A Chart block that drew nothing still carries
  // its title as text, so the text-emptiness check can never see it.
  const inkOf = root =>
    [...root.querySelectorAll('canvas')].map(c => {
      if (!c.width || !c.height) return 0
      try {
        const d = c.getContext('2d').getImageData(0, 0, c.width, c.height).data
        let ink = 0
        for (let i = 3; i < d.length; i += 4 * 16) if (d[i] > 8) ink++
        return ink
      } catch {
        return -1 // tainted canvas; unknowable, so never reported as blank
      }
    })

  const blocks = boxes
    .filter(b => b.content)
    .map(b => {
      // textContent, not innerText: innerText is layout-dependent and returns
      // only the visible part of a subtree, so a clipped block can read empty
      const text = (b.content.textContent || '').trim()
      return {
        xywh: b.xywh,
        label: (b.content.innerText || text).trim().split('\n')[0].slice(0, 38),
        empty: !text,
        ink: inkOf(b.content),
        ...measure(b.content),
      }
    })

  // A snowflake ID rendered where a label belongs — the signature of a ref
  // the webapp could not resolve.
  const rawIds = new Set()
  const templates = new Set()
  const walk = document.createTreeWalker(document.body, NodeFilter.SHOW_TEXT)
  for (let n = walk.nextNode(); n; n = walk.nextNode()) {
    const t = (n.nodeValue || '').trim()
    if (/^\d{15,20}$/.test(t)) rawIds.add(t)
    if (t.includes('${')) templates.add(t.slice(0, 60))
  }
  return { blocks, rawIds: [...rawIds].slice(0, 6), templates: [...templates].slice(0, 4) }
}

// A fixed wait cannot cover a cold vite load, a records fetch and a chart.js
// render at once — charts came back blank often enough to make an OK verdict
// worthless. Poll the canvases until their ink stops changing instead, so the
// report describes a settled page rather than whatever was on screen at 1.5s.
const settleCanvases = async () => {
  const signature = () =>
    page
      .evaluate(() =>
        [...document.querySelectorAll('canvas')]
          .map(c => {
            if (!c.width || !c.height) return '0'
            try {
              const d = c.getContext('2d').getImageData(0, 0, c.width, c.height).data
              let ink = 0
              for (let i = 3; i < d.length; i += 4 * 16) if (d[i] > 8) ink++
              return String(ink)
            } catch {
              return 'x'
            }
          })
          .join('|'),
      )
      .catch(() => '')

  let prev = await signature()
  if (prev === '') return // no canvases on this page
  for (let i = 0; i < 12; i++) {
    await page.waitForTimeout(400)
    const now = await signature()
    if (now === prev) return
    prev = now
  }
}

// A TAQ builder (and anything else drawn with vue-flow) has no compose blocks,
// so the block report finds nothing there. Read the graph instead: the failure
// this catches is a severed chain — an automation that runs correctly while its
// canvas draws a step no trigger reaches, which is what an empty `paths` array
// produces.
//
// Reachability is the test, not the shape. A gatewayExclusive draws one End per
// branch and a TAQ may carry several triggers, so counting Ends or triggers
// calls a correct branching graph broken — and tells its author to go fix
// `paths` that were right.
const FLOW_REPORT = () => {
  const nodes = [...document.querySelectorAll('.vue-flow__node')].map(n => ({
    id: n.getAttribute('data-id') || '',
    kind: ([...n.classList].find(c => c.startsWith('vue-flow__node-')) || '').replace(
      'vue-flow__node-',
      '',
    ),
    label: (n.textContent || '').trim().slice(0, 30),
  }))
  if (!nodes.length) return null
  // vue-flow names an edge "<sourceID>_<targetID>"; node ids contain
  // underscores of their own, so test membership rather than splitting
  const edges = [...document.querySelectorAll('.vue-flow__edge')].map(
    e => e.getAttribute('data-id') || '',
  )
  const touches = id => id && edges.some(e => e.startsWith(id + '_') || e.endsWith('_' + id))

  // Which node an edge end names, longest id first: ids are not prefix-free,
  // so "1" would claim an edge belonging to "1001" if the short one won.
  const byLength = [...nodes].sort((a, b) => b.id.length - a.id.length)
  const source = e => byLength.find(n => n.id && e.startsWith(n.id + '_'))
  const target = e => byLength.find(n => n.id && e.endsWith('_' + n.id))

  const next = new Map()
  for (const e of edges) {
    const s = source(e)
    const t = target(e)
    if (!s || !t) continue
    if (!next.has(s.id)) next.set(s.id, [])
    next.get(s.id).push(t.id)
  }

  // A severed chain is a step no trigger can reach. Counting Ends cannot see
  // it: a gateway legitimately ends in one End per branch, and a TAQ may carry
  // several triggers, so both are normal shapes rather than defects.
  const seen = new Set()
  const queue = nodes.filter(n => n.kind === 'trigger').map(n => n.id)
  queue.forEach(id => seen.add(id))
  while (queue.length) {
    for (const to of next.get(queue.shift()) || []) {
      if (seen.has(to)) continue
      seen.add(to)
      queue.push(to)
    }
  }

  return {
    nodeCount: nodes.length,
    edgeCount: edges.length,
    triggers: nodes.filter(n => n.kind === 'trigger').length,
    ends: nodes.filter(n => n.kind === 'end').length,
    isolated: nodes.filter(n => !touches(n.id)).map(n => n.label || n.id),
    unreachable: nodes
      .filter(n => n.kind !== 'trigger' && touches(n.id) && !seen.has(n.id))
      .map(n => n.label || n.id),
  }
}

// fullPage stops at the viewport because compose scrolls an inner container,
// not the document — so grow the viewport to the page's own scroll height
// before capturing, or every block below the fold goes unverified.
const fitViewport = async () => {
  const needed = await page.evaluate(() => {
    let node = document.querySelector('.grid-stack, .single-block-wrapper')
    while (node && node !== document.body) {
      const st = getComputedStyle(node)
      if (['auto', 'scroll'].includes(st.overflowY) && node.scrollHeight > node.clientHeight + 4)
        return Math.ceil(node.getBoundingClientRect().top + node.scrollHeight + 24)
      node = node.parentElement
    }
    return Math.ceil(document.documentElement.scrollHeight)
  })
  const height = Math.min(Math.max(needed, 900), 5000)
  if (height > 900) {
    await page.setViewportSize({ width: 1440, height })
    await page.waitForTimeout(400) // gridstack relayouts on resize
  }
  return height
}

for (const p of paths) {
  // stale vite optimized-dep chunks make first loads fail with
  // "App setup failed … app.use"; a expired session answers 401 to every
  // request and renders an empty shell. Both survive exactly one retry.
  for (let attempt = 0; attempt < 2; attempt++) {
    problems.length = 0
    await page.setViewportSize({ width: 1440, height: 900 })
    await page
      .goto(WEBAPP + p, { waitUntil: 'networkidle', timeout: 30000 })
      .catch(e => problems.push(`nav: ${e.message}`))
    await page.waitForTimeout(1500) // charts/metrics fetch after load
    if (attempt === 1) break
    if (problems.some(x => x.includes('App setup failed'))) continue
    if (problems.some(x => x.includes('401'))) {
      await login() // session expired mid-run
      continue
    }
    break
  }
  // Where the app ended up. An unmatched path is not an error anywhere — the
  // webapp's catch-all route redirects it to the home section and that page
  // renders cleanly — so without this check a typo'd path passes, and the
  // screenshot shows a perfectly healthy page nobody asked about.
  const landed = new URL(page.url()).pathname
  const bare = s => s.split(/[?#]/)[0].replace(/\/$/, '')
  if (bare(landed) !== bare(p)) {
    problems.push(
      `redirected: asked for ${p}, settled on ${landed} — the route may redirect, or nothing matched ` +
        `the path (compose pages are /compose/namespace/<slug>/pages/<pageID>)`,
    )
  }

  await settleCanvases()
  const report = await page.evaluate(PAGE_REPORT).catch(() => null)
  const flow = await page.evaluate(FLOW_REPORT).catch(() => null)
  const height = await fitViewport()
  const shot = join(outDir, p.replace(/[^a-z0-9-]+/gi, '_').replace(/^_+|_+$/g, '') + '.png')
  await page.screenshot({ path: shot, fullPage: true })

  // Findings are defects the DOM can prove; they are reported loudly but do
  // not fail the run, which stays reserved for page/HTTP errors.
  const findings = []
  const notes = []
  for (const b of report?.blocks || []) {
    const where = b.xywh ? `[${b.xywh.join(',')}]` : '[single]'
    const name = b.label || '(no text)'
    if (b.empty) {
      findings.push(`empty block ${where} — rendered nothing`)
      continue
    }
    if (b.lost.v > 8)
      findings.push(
        `clipped ${where} "${name}" — ${b.lost.v}px of content unreachable; ` +
          `raise h by ~${Math.ceil(b.lost.v / 10)}`,
      )
    if (b.lost.h > 8)
      findings.push(
        `columns cut off ${where} "${name}" — ${b.lost.h}px hidden; ` +
          `drop a field or widen the block`,
      )
    if (b.ink.length && b.ink.every(i => i === 0))
      findings.push(`blank canvas ${where} "${name}" — a chart block that drew nothing`)
    if (b.listScroll.v > 8)
      notes.push(
        `${where} "${name}" scrolls internally — ${b.listScroll.v}px of rows hidden; ` +
          // A lone block renders through a flex wrapper that fills the view, so
          // its xywh is never read and "raise h" would be advice that cannot work.
          (b.xywh
            ? `raise h by ~${Math.ceil(b.listScroll.v / 10)} unless the list really holds ` +
              `more rows than its perPage`
            : `this block sets the page height itself, so lower perPage if those rows matter`),
      )
  }
  if (flow) {
    if (flow.unreachable.length)
      findings.push(
        `no trigger reaches ${flow.unreachable.map(n => `"${n}"`).join(', ')} — ` +
          `a severed chain; check "paths" wires each step to a trigger`,
      )
    if (!flow.triggers) findings.push(`graph has no trigger node — nothing can start it`)
    for (const n of flow.isolated) findings.push(`node "${n}" is joined to nothing`)
  }
  for (const id of report?.rawIds || [])
    findings.push(`raw ID rendered: ${id} — a ref the webapp could not resolve to a label`)
  for (const t of report?.templates || []) findings.push(`uninterpolated template: ${t}`)

  const status = problems.length ? 'FAIL' : findings.length ? 'WARN' : 'OK'
  if (problems.length) failed = true
  console.log(`${status}  ${p}`)
  problems.forEach(x => console.log(`      ${x}`))
  findings.forEach(x => console.log(`      ${x}`))
  notes.forEach(x => console.log(`      note: ${x}`))
  if (report) {
    const n = report.blocks.length
    if (!n && flow) {
      console.log(
        `      graph: ${flow.nodeCount} nodes, ${flow.edgeCount} edges, ` +
          `${flow.triggers} trigger, ${flow.ends} End` +
          (findings.length ? '' : ' — every step reachable from a trigger'),
      )
    } else if (!n) {
      // No compose blocks here (a TAQ builder, an admin screen). Saying
      // "none clipped or empty" over one would be a health verdict drawn from
      // an empty inspection — reassuring, and true of a visibly broken page.
      console.log(`      no compose blocks on this path — nothing the report can judge`)
      console.log(`      read the screenshot; this check cannot vouch for it`)
    } else {
      console.log(
        `      ${n} block${n === 1 ? '' : 's'} rendered` +
          (findings.length ? '' : ', none clipped or empty') +
          `; captured ${height}px`,
      )
    }
  }
  console.log(`      screenshot: ${shot}`)
}

await browser.close()
process.exit(failed ? 1 : 0)
