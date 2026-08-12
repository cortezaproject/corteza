#!/usr/bin/env node
// Render-verify webapp paths in a real (headless) browser.
//
// Usage: node dev/agent/verify-ui.mjs [--out DIR] PATH [PATH...]
//   node dev/agent/verify-ui.mjs '/compose/namespace/agent-sandbox/pages/123'
//
// Logs in as agent-ui@local.dev (password from .state/ui-password, created by
// bootstrap.sh), navigates each path on the vite dev server, captures console
// and page errors, and writes a full-page screenshot per path. Exits non-zero
// if any path produced page errors. Local-only.
//
// Why this exists: API checks cannot see rendering bugs (mis-sized grids,
// blocks that clip to nothing, unresolved refs showing raw IDs). Both
// toolkit field tests caught real bugs only through this kind of check.

import { readFileSync, readdirSync, mkdirSync } from 'node:fs'
import { dirname, join, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { createRequire } from 'node:module'

const AGENT_DIR = dirname(fileURLToPath(import.meta.url))
const REPO_DIR = resolve(AGENT_DIR, '..', '..')
const WEBAPP = process.env.HUMAN_WEBAPP || 'http://localhost:5173'
const API_BASE = (process.env.HUMAN_API || 'http://localhost:1043/api').replace(/\/api$/, '')

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
let outDir = join(AGENT_DIR, '.state', 'ui')
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

const email = 'agent-ui@local.dev'
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

let failed = false
for (const p of paths) {
  // stale vite optimized-dep chunks make first loads fail with
  // "App setup failed … app.use" — a reload fixes it, so retry once
  for (let attempt = 0; attempt < 2; attempt++) {
    problems.length = 0
    await page
      .goto(WEBAPP + p, { waitUntil: 'networkidle', timeout: 30000 })
      .catch(e => problems.push(`nav: ${e.message}`))
    await page.waitForTimeout(1500) // charts/metrics fetch after load
    if (attempt === 0 && problems.some(x => x.includes('App setup failed'))) continue
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

  const shot = join(outDir, p.replace(/[^a-z0-9-]+/gi, '_').replace(/^_+|_+$/g, '') + '.png')
  await page.screenshot({ path: shot, fullPage: true })
  const status = problems.length ? 'FAIL' : 'OK'
  if (problems.length) failed = true
  console.log(`${status}  ${p}`)
  problems.forEach(x => console.log(`      ${x}`))
  console.log(`      screenshot: ${shot}`)
}

await browser.close()
process.exit(failed ? 1 : 0)
