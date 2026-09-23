#!/usr/bin/env node
// Render-verify custom applications in a real (headless) browser.
//
// Usage: node dev/agent/verify-app.mjs [--out DIR] [--expect TEXT] [--click N] ID [ID...]
//   node dev/agent/verify-app.mjs --expect 'Ada Lovelace' 514958620159967233
//
// `--click N` presses up to N of the app's own buttons afterwards and agrees to
// whatever Human asks, which is the only way to reach a page that saves: the
// consent dialog belongs to the shell, so a page that changes records does
// nothing at all until somebody answers it.
//
// Logs in as agent@local.dev (password from .state/ui-password, created by
// bootstrap.sh), opens /app/<id> for each application and reports what the app
// actually drew. Exits non-zero if any app failed to render. Local-only.
//
// Why this exists: a custom app can pass every static check — the guard
// accepts it, the snippet is intact, it deploys — and still show an empty
// list, a column of record IDs, or "undefined" in every row. None of that is
// visible to the API, to the deploy guard, or to the source itself. It is only
// visible here.
//
// A custom app runs two frames deep: the shell holds an outer `srcdoc` frame
// carrying the CSP, and that holds the sandboxed frame the page runs in. The
// sandbox has an opaque origin, but the browser is being driven rather than
// scripted, so the page's DOM is readable from here.

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
    console.error(`verify-app is local-only; refusing ${u}`)
    process.exit(1)
  }
}

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
let outDir = join(AGENT_DIR, '.state', 'sessions', SESSION, basename(REPO_DIR), 'app')
const expects = []
const ids = []
let clicks = 0
for (let i = 0; i < args.length; i++) {
  if (args[i] === '--out') outDir = args[++i]
  else if (args[i] === '--expect') expects.push(args[++i])
  else if (args[i] === '--click') clicks = Number(args[++i]) || 6
  else ids.push(args[i])
}
if (!ids.length) {
  console.error('usage: verify-app.mjs [--out DIR] [--expect TEXT] ID [ID...]')
  process.exit(1)
}
mkdirSync(outDir, { recursive: true })

// Requests the shell makes on every route, custom app or not. They are the
// webapp's business and would otherwise be reported against every app.
const SHELL_NOISE = [/\/federation\/permissions\/effective/, /\/system\/attachment\/avatar\//]

// What gets typed into a page's own fields, so a record it creates can be told
// apart from the fixture's afterwards.
const PROBE_MARK = 'probe-' + Date.now().toString(36)

const email = 'agent@local.dev'
const password = readFileSync(join(AGENT_DIR, '.state', 'ui-password'), 'utf8').trim()

const browser = await chromium.launch({ headless: true })
const page = await browser.newPage({ viewport: { width: 1440, height: 900 } })

let problems = []
page.on('pageerror', e => problems.push(`pageerror: ${e.message.slice(0, 300)}`))
page.on('console', m => {
  // resource-load failures are reported (and filtered) through the response
  // listener below — the console duplicate has no URL and only adds noise
  if (m.type() === 'error' && !m.text().startsWith('Failed to load resource'))
    problems.push(`console.error: ${m.text().slice(0, 300)}`)
})
page.on('requestfailed', r => {
  // Leaving a route cancels whatever it still had in flight; only a request the
  // browser refused says anything about the app.
  const why = r.failure()?.errorText || ''
  if (why === 'net::ERR_ABORTED') return
  if (!SHELL_NOISE.some(re => re.test(r.url()))) problems.push(`${why}: ${r.url().slice(0, 160)}`)
})
page.on('response', r => {
  if (r.status() >= 400 && !SHELL_NOISE.some(re => re.test(r.url())))
    problems.push(`HTTP ${r.status()}: ${r.url().slice(0, 160)}`)
})

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

// What the app drew, read out of the frame it drew it in.
//
// Every fault below survives the deploy guard, so this is the only place they
// can be seen: a Human ID printed where a name belongs (the page read a field
// the module does not have), `undefined` in a cell (a value the record left
// empty), a frame holding nothing at all (the app threw before drawing).
const APP_REPORT = () => {
  const text = document.body ? document.body.innerText : ''
  // Repeated things, whatever the page built them out of: a table, a list, or
  // the grid of cards a page is just as likely to draw. Counted for the report
  // rather than judged — the honest test of live data is the text itself.
  const rows = document.querySelectorAll(
    'tbody tr, li, [role="row"], [role="listitem"], [class*="card"], [class*="row"], [class*="item"]',
  ).length
  // A Human ID is 18-19 digits. Shown where a name belongs it is a mistake —
  // the page read the wrong field, or never looked the reference up. Shown
  // under a label that says it is an ID, it is a record's own detail and is
  // exactly what the reader asked to see.
  const excused = new Set()
  for (const el of document.querySelectorAll('*')) {
    if (el.children.length) continue
    const own = (el.textContent || '').trim()
    if (!/^\d{18,19}$/.test(own)) continue
    const beside = el.previousElementSibling?.textContent || ''
    const around = el.parentElement ? (el.parentElement.textContent || '').replace(own, '') : ''
    if (/\bid\b/i.test(beside + ' ' + around)) excused.add(own)
  }
  const ids = [...new Set(text.match(/\b\d{18,19}\b/g) || [])].filter(id => !excused.has(id))
  const holes = ['undefined', 'NaN', '[object Object]', 'null'].filter(h =>
    new RegExp(`(^|[\\s>])${h.replace(/[[\]]/g, '\\$&')}([\\s<]|$)`).test(text),
  )
  // Text the reader cannot get to. Two things are fine and only one is a
  // fault: a container with overflow-x auto/scroll can be scrolled, and a cell
  // that truncates with an ellipsis says so on screen. Content cut off by
  // `overflow: hidden` with neither is simply gone, and looks complete in a
  // preview that happened to be wider.
  const root = document.documentElement
  const clipped = Math.max(0, root.scrollWidth - root.clientWidth)
  const lost = []
  for (const el of document.querySelectorAll('*')) {
    if (el.clientWidth < 80 || el.children.length) continue
    const over = el.scrollWidth - el.clientWidth
    if (over <= 4) continue
    const st = getComputedStyle(el)
    if (st.overflowX !== 'hidden' || st.textOverflow === 'ellipsis') continue
    lost.push(
      `${el.tagName.toLowerCase()} loses ${over}px of "${(el.textContent || '').trim().slice(0, 30)}"`,
    )
  }
  return {
    text,
    rows,
    ids,
    holes,
    clipped,
    lost,
    canvases: document.querySelectorAll('canvas').length,
  }
}

const MODE = text => {
  if (/\blive data\b/i.test(text)) return 'live'
  if (/\bsample data\b/i.test(text)) return 'sample'
  return 'unstated'
}

let failed = false
for (const id of ids) {
  problems = []
  const faults = []
  await page.goto(`${WEBAPP}/app/${id}`, { waitUntil: 'networkidle' }).catch(e => {
    faults.push(`navigation: ${e.message.split('\n')[0]}`)
  })
  // The bridge answers over a port, so the first paint is empty by design.
  await page.waitForTimeout(2500)

  // shell, outer host, sandboxed app — in that order, the app being the last
  // srcdoc frame the shell holds.
  const srcdoc = page.frames().filter(f => f.url() === 'about:srcdoc')
  const app = srcdoc[srcdoc.length - 1]

  let report = null
  if (!srcdoc.length) faults.push('no sandbox frame — the app never rendered')
  else if (srcdoc.length < 2) faults.push('the app frame is gone — it navigated itself away')
  else {
    try {
      report = await app.evaluate(APP_REPORT)
    } catch (e) {
      faults.push(`unreadable app frame: ${e.message.split('\n')[0]}`)
    }
  }

  // Fill what the page asks for, press what it offers, and agree to whatever
  // Human asks about it. A page that saves is inert until the consent dialog is
  // answered, and a page that saves something new is inert until its form has
  // been filled — press-only leaves both reading as working pages.
  //
  // A search or filter box is left alone: typing in one hides the very rows the
  // checks above were made against.
  let pressed = 0
  let typed = 0
  if (clicks && report) {
    // Never pressed: these undo the work rather than doing it, and pressing a
    // dialog's close button straight after its open button is how a form that
    // was about to be filled disappears again.
    const DISMISSES = /close|cancel|dismiss|back|reset|clear|delete|remove|×|✕/i
    const COMMITS = /save|add|create|submit|apply|confirm|update|mark|set|send/i

    const controls = app.locator(
      'button:visible, [role="button"]:visible, [role="switch"]:visible, ' +
        'label:has(input[type="checkbox"]):visible, label:has(input[type="radio"]):visible, ' +
        'input[type="checkbox"]:visible, input[type="radio"]:visible',
    )

    // A form's own submit button is the commit, whatever it is called. Pressing
    // it before anything else that merely reads like one keeps a button called
    // "+ Add contact" from reopening — and emptying — the form just filled in.
    const submits = app.locator('button[type="submit"]:visible, input[type="submit"]:visible')

    const press = async only => {
      if (only === 'commits') {
        const count = Math.min(await submits.count().catch(() => 0), clicks)
        for (let i = 0; i < count && pressed < clicks; i++) {
          await submits
            .nth(i)
            .click({ timeout: 3000, noWaitAfter: true })
            .catch(() => {})
          pressed++
          const agree = page.getByRole('alertdialog').getByRole('button', { name: /^Allow/ })
          if (await agree.isVisible().catch(() => false)) await agree.click().catch(() => {})
          await page.waitForTimeout(400)
        }
      }
      const total = Math.min(await controls.count().catch(() => 0), 20)
      for (let i = 0; i < total && pressed < clicks; i++) {
        const one = controls.nth(i)
        const words =
          ((await one.innerText().catch(() => '')) || '') +
          ' ' +
          ((await one.getAttribute('aria-label').catch(() => '')) || '')
        if (DISMISSES.test(words)) continue
        if (only === 'commits' && !COMMITS.test(words)) continue
        if (only === 'opens' && COMMITS.test(words)) continue
        await one.click({ timeout: 3000, noWaitAfter: true }).catch(() => {})
        pressed++
        const allow = page.getByRole('alertdialog').getByRole('button', { name: /^Allow/ })
        if (await allow.isVisible().catch(() => false)) await allow.click().catch(() => {})
        await page.waitForTimeout(400)
      }
    }

    const fill = async () => {
      const boxes = app.locator(
        'input:visible:not([type="search"]):not([type="checkbox"]):not([type="radio"]), textarea:visible',
      )
      const fields = Math.min(await boxes.count().catch(() => 0), 12)
      for (let i = 0; i < fields; i++) {
        const box = boxes.nth(i)
        const what = (
          ((await box.getAttribute('placeholder').catch(() => '')) || '') +
          ' ' +
          ((await box.getAttribute('name').catch(() => '')) || '') +
          ' ' +
          ((await box.getAttribute('id').catch(() => '')) || '')
        ).toLowerCase()
        if (/search|filter|query/.test(what)) continue
        if (await box.inputValue().catch(() => 'x')) continue // leave what is already there
        const type = (await box.getAttribute('type').catch(() => '')) || 'text'
        const value =
          type === 'number'
            ? '7'
            : type === 'date'
              ? '2026-09-23'
              : type === 'email'
                ? `${PROBE_MARK}@example.tld`
                : PROBE_MARK
        await box.fill(value, { timeout: 2000 }).catch(() => {})
        typed++
      }
    }

    // What a person does: open the thing, fill it in, commit it. Twice over,
    // because the button that reveals a form is as likely to be called "Add
    // contact" as anything else, and a form only exists to be filled once
    // whatever hides it has been pressed.
    for (const round of [1, 2]) {
      await press('opens')
      await fill()
      await press('commits')
      if (round === 1) await page.waitForTimeout(400)
    }

    // What the page said for itself before anything was pressed still decides
    // whether it drew live data; pressing can only add faults, never excuse one.
    try {
      const after = await app.evaluate(APP_REPORT)
      report.holes = [...new Set([...report.holes, ...after.holes])]
      report.ids = [...new Set([...report.ids, ...after.ids])]
      report.lost = [...report.lost, ...after.lost]
    } catch {
      faults.push('the app frame stopped answering after its own controls were used')
    }
  }

  const shot = join(outDir, `app-${id}.png`)
  await page.screenshot({ path: shot, fullPage: true }).catch(() => {})

  console.log(`\n▸ /app/${id}`)
  if (pressed || typed)
    console.log(`  used        ${typed} field(s), ${pressed} control(s) · typed ${PROBE_MARK}`)
  if (report) {
    const mode = MODE(report.text)
    const visible = report.text.trim()
    console.log(`  mode        ${mode}`)
    console.log(`  drew        ${report.rows} row(s), ${visible.length} chars of text`)
    if (mode === 'unstated') faults.push('the page never says whether it is on live or sample data')
    if (mode === 'sample') faults.push('the page fell back to sample data inside Human')
    if (visible.length < 40) faults.push(`the app drew almost nothing: ${JSON.stringify(visible)}`)
    if (report.ids.length)
      faults.push(`Human IDs shown to the reader: ${report.ids.slice(0, 3).join(', ')}`)
    if (report.holes.length) faults.push(`empty values rendered as ${report.holes.join(', ')}`)
    if (report.clipped > 4)
      faults.push(`${report.clipped}px of the page is cut off the right edge and cannot be reached`)
    for (const l of report.lost.slice(0, 3))
      faults.push(`text cut off with no way to read it: ${l}`)
    for (const want of expects) {
      if (!report.text.includes(want))
        faults.push(`expected text not on the page: ${JSON.stringify(want)}`)
    }
    console.log(`  text        ${JSON.stringify(visible.slice(0, 200))}`)
  }

  for (const p of problems) faults.push(p)
  if (faults.length) {
    failed = true
    for (const f of faults) console.log(`  ✗ ${f}`)
  } else {
    console.log('  ✓ renders live data with no faults')
  }
  console.log(`  screenshot  ${shot}`)
}

await browser.close()
process.exit(failed ? 1 : 0)
