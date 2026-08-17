#!/usr/bin/env node
// Drive the webapp in a real browser: multi-step checks, not just "does it render".
//
// verify-ui.mjs answers "did this path load clean". Anything that needs a
// second step — click Back and assert where you land, search and count rows,
// leave an editor and see whether it warns — meant hand-writing a script with
// ~80 lines of launch/login/listener preamble, every time. This is that
// preamble, once.
//
//   import { drive, check, expectPath } from './drive.mjs'
//
//   drive('back leaves the editor', async page => {
//     await page.open('/compose/namespaces/edit/catalogue')
//     await page.back()
//     expectPath(page, '/compose/namespaces')
//   })
//
// Run it: node dev/agent/drive.mjs my-checks.mjs [--only NAME] [--headed]
// The runner calls run() once the suite is loaded; a suite never calls it.
//
// What you get for free, because each cost a wrong answer once:
//  - logged in already, and the session is cached so later runs skip the form
//  - dialogs are RECORDED, not auto-dismissed (playwright's default dismiss
//    turns an unsaved-changes confirm into a phantom "the button is broken")
//  - console errors, page errors and HTTP >= 500 collected per check and
//    failing the check unless it opts out
//  - screenshots only on failure, because reading one costs more than the
//    assertion it replaces
import { readFileSync, readdirSync, mkdirSync, existsSync, writeFileSync } from 'node:fs'
import { dirname, join, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { createRequire } from 'node:module'

const AGENT_DIR = dirname(fileURLToPath(import.meta.url))
const REPO_DIR = resolve(AGENT_DIR, '..', '..')
const STATE_DIR = join(AGENT_DIR, '.state')
const SHOT_DIR = join(STATE_DIR, 'drive')
const STORAGE = join(STATE_DIR, 'ui-storage.json')

const WEBAPP = process.env.HUMAN_WEBAPP || 'http://localhost:5173'
const API_BASE = (process.env.HUMAN_API || 'http://localhost:1043/api').replace(/\/api$/, '')

// Per-action ceiling. Long enough for a slow render, short enough that a
// selector matching nothing reports as a failure rather than a hang.
const ACTION_TIMEOUT = Number(process.env.DRIVE_TIMEOUT || 8000)

// What an element is given to appear before it counts as absent. Higher than
// the action ceiling on purpose: waiting for something to render and asserting
// it exists are different questions, and a cold vite compiling a route for the
// first time answers the first one slowly without changing the second.
const RENDER_TIMEOUT = Number(process.env.DRIVE_RENDER_TIMEOUT || 20000)

for (const u of [WEBAPP]) {
  if (!['localhost', '127.0.0.1', '::1'].includes(new URL(u).hostname)) {
    console.error(`drive is local-only; refusing ${u}`)
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

// --- ids ---------------------------------------------------------------------

// The handle → ID map ids.sh writes. Checks address fixtures by handle so an
// ID never has to be pasted into a script (and never goes stale in one).
export function ids(slug) {
  const path = join(STATE_DIR, 'ids.json')
  if (!existsSync(path)) {
    throw new Error(`no ${path} — run: dev/agent/ids.sh ${slug || ''}`.trim())
  }
  const all = JSON.parse(readFileSync(path, 'utf8'))
  if (!slug) return all
  if (!all[slug]) {
    throw new Error(`no ids for '${slug}' — run: dev/agent/ids.sh ${slug}`)
  }
  return all[slug]
}

// --- checks ------------------------------------------------------------------

const cases = []
const results = []

/** Register a check. Body receives a driven page. */
export function drive(name, body, opts = {}) {
  cases.push({ name, body, opts })
}

let current = null

/** Assert, with the failure detail that makes it diagnosable. */
export function check(name, ok, detail = '') {
  current.checks.push({ name, ok: !!ok, detail })
  return !!ok
}

/** Assert the app settled on a path (the query string is ignored). */
export function expectPath(page, expected, name = `settles on ${expected}`) {
  const at = page.path()
  return check(name, at === expected, at === expected ? '' : `at ${at}`)
}

// --- the driven page ---------------------------------------------------------

function drivePage(page, state) {
  const api = {
    raw: page,

    path: () => new URL(page.url()).pathname,
    url: () => page.url(),

    /** Navigate, then wait for the shell to actually be up.
     *
     *  networkidle is deliberately not used: editors holding a live connection
     *  (the chatbot inbox) never reach it. Nor is a fixed sleep enough on its
     *  own — vite compiles a route the first time it is asked for, so the first
     *  check of a run waits seconds longer than the rest and a timeout tuned to
     *  the warm case reports a cold one as a missing element. The topbar
     *  landmark is the honest signal that the app has mounted. */
    async open(path, { settle = 800 } = {}) {
      await page.goto(WEBAPP + path, { waitUntil: 'domcontentloaded', timeout: 30000 })
      await page
        .locator('[data-testid="app-topbar"]')
        .waitFor({ state: 'visible', timeout: 30000 })
        .catch(() => {})
      await page.waitForTimeout(settle)
      return api
    },

    async click(selector, { hasText, nth = 0, settle = 2000 } = {}) {
      let loc = page.locator(selector)
      if (hasText) loc = loc.filter({ hasText })
      await loc.nth(nth).waitFor({ state: 'visible', timeout: RENDER_TIMEOUT })
      await loc.nth(nth).click()
      await page.waitForTimeout(settle)
      return api
    },

    async fill(selector, value, { settle = 800 } = {}) {
      const loc = page.locator(selector).first()
      await loc.waitFor({ state: 'visible', timeout: RENDER_TIMEOUT })
      await loc.fill(value)
      await page.waitForTimeout(settle)
      return api
    },

    /** Click an editor's Back button. The commonest interaction there is, and
     *  the one whose selector was guessed at most often before the shell grew
     *  a test id for it. */
    async back({ settle = 2500 } = {}) {
      const loc = page.locator('[data-testid="editor-back"]').first()
      await loc.waitFor({ state: 'visible', timeout: RENDER_TIMEOUT })
      await loc.click()
      await page.waitForTimeout(settle)
      return api
    },

    async text(selector) {
      const loc = page.locator(selector).first()
      if (!(await loc.count())) return ''
      return (await loc.innerText()).trim()
    },

    /** Visible rows of the shell's left sidebar — a PrimeVue Drawer, so the
     *  landmark carries data-testid="app-sidebar".
     *
     *  The drawer slides open after the route renders, and on a cold vite it is
     *  slower than any fixed settle: reading straight away returns [] and the
     *  check reports an empty nav that is merely late. So it waits for a row,
     *  and only an empty sidebar actually returns nothing. */
    async sidebarRows({ expectRows = true } = {}) {
      if (expectRows) {
        await page
          .locator('[data-testid="app-sidebar"] button')
          .first()
          .waitFor({ state: 'visible' })
          .catch(() => {})
      }
      const rows = await page.locator('[data-testid="app-sidebar"] button').allInnerTexts()
      return rows.map(s => s.trim()).filter(Boolean)
    },

    /** The shell remembers sidebar expansion per section; a section whose
     *  default is collapsed renders zero rows and reads as a broken nav.
     *
     *  Usually called before the first open(), where the page may still be on
     *  about:blank — and localStorage there is a SecurityError, not an empty
     *  store. So the app origin is established first. */
    async expandSidebar(section) {
      // Always navigated, never conditionally: the page may be on about:blank,
      // on the auth server's origin, or midway through a redirect that swaps
      // the document out from under evaluate() — and all three raise a
      // SecurityError on localStorage rather than returning an empty store.
      // Landing on the app root first makes the origin a fact, not a guess.
      await api.open('/')
      await page.evaluate(
        s => localStorage.setItem('ui.sidebar.expanded', JSON.stringify({ [s]: true })),
        section,
      )
      return api
    },

    /** The last native dialog seen, or null. Dialogs are accepted so the flow
     *  continues, and recorded so a check can assert on them. */
    dialog: () => state.dialog,
    clearDialog: () => {
      state.dialog = null
    },

    problems: () => state.problems.slice(),
    clearProblems: () => {
      state.problems.length = 0
    },

    evaluate: (...a) => page.evaluate(...a),
    locator: (...a) => page.locator(...a),
  }

  return api
}

// index.html requests both of these unconditionally — they are the deployment's
// optional CSS and JS customization hooks — and the dev server answers 500 when
// none is configured. Every page load in dev carries them, so treating them as
// failures would fail every check forever and teach the reader to ignore the
// one line that matters.
const DEV_NOISE = ['/custom.css', '/code-snippets.js']

function isKnownDevNoise(url) {
  const { pathname } = new URL(url)
  return DEV_NOISE.includes(pathname)
}

// A repo edited underneath a run breaks the app in ways that are not results:
// the API restarting mid-request (a concurrent session rebuilding Go), or vite
// serving a module halfway through someone's save. Both leave every check
// failing on a page that never mounted, both read as product bugs, and both
// cost a full diagnosis to disprove. So they are named instead.
const CHURN = [
  'App setup failed',
  'Network Error',
  'ERR_CONNECTION',
  'does not provide an export named',
  'Failed to fetch dynamically imported module',
]

function looksLikeEnvironmentChurn(problems) {
  return problems.some(p => CHURN.some(sig => p.includes(sig)))
}

// --- session -----------------------------------------------------------------

async function login(context, page) {
  const email = 'agent-ui@local.dev'
  const password = readFileSync(join(STATE_DIR, 'ui-password'), 'utf8').trim()

  // The caller has already loaded the app and seen it bounce to auth. Loading
  // it again from here restarts that redirect chain underneath the page we are
  // about to type into, and the form goes out from under the locator.
  if (!page.url().includes('/auth/')) {
    await page.goto(WEBAPP, { waitUntil: 'domcontentloaded' })
    await page.waitForURL(/\/auth\//, { timeout: 20000 }).catch(() => {})
  }
  if (!page.url().includes('/auth/')) return false

  // The form is server-rendered by the auth app after its own redirect, so it
  // is waited for rather than assumed present the moment the URL matches.
  await page
    .locator('input[name="email"]')
    .waitFor({ state: 'visible', timeout: 30000 })
    .catch(() => {
      throw new Error(`the auth form never appeared; the page settled on ${page.url()}`)
    })
  await page.fill('input[name="email"]', email)
  if (!(await page.locator('input[name="password"]').count())) {
    await page.click('button[type="submit"]')
  }
  await page.fill('input[name="password"]', password)
  await page.click('button[type="submit"]')
  await page.waitForURL(u => !u.href.includes('/auth/'), { timeout: 20000 })

  // Cached so later runs skip the form. It also removes a race that reads as a
  // product bug: a fresh context deep-linked before its session settled lands
  // on /auth/login, which looks exactly like a route that refuses to open.
  writeFileSync(STORAGE, JSON.stringify(await context.storageState()))
  return true
}

// --- runner ------------------------------------------------------------------

function report(r) {
  const bad = r.checks.filter(c => !c.ok)
  console.log(`${bad.length ? 'FAIL' : 'PASS'}  ${r.name}`)
  for (const c of r.checks) {
    if (!c.ok || process.env.VERBOSE) {
      console.log(`   ${c.ok ? 'ok  ' : 'BAD '} ${c.name}${c.detail ? ' — ' + c.detail : ''}`)
    }
  }
  if (r.shot) console.log(`   screenshot: ${r.shot}`)
}

// Whether the checks have already been run in this process. The CLI reads it
// to keep a suite that runs itself from being run a second time.
let hasRun = false

export async function run({ only = null, headed = false } = {}) {
  hasRun = true
  const picked = only ? cases.filter(c => c.name.toLowerCase().includes(only.toLowerCase())) : cases

  if (!picked.length) {
    console.error(only ? `no check matches --only ${only}` : 'no checks registered')
    process.exit(1)
  }

  // A dev server that is mid-restart — another session rebuilding Go, most
  // often — makes every check fail on a page whose API calls never landed.
  // Those failures look exactly like product bugs and cost a full diagnosis
  // before the cause turns out to be the clock. Say it up front instead.
  for (const [what, url] of [
    ['the webapp (vite)', WEBAPP],
    ['the API', API_BASE],
  ]) {
    const reachable = await fetch(url, { method: 'HEAD' }).then(
      () => true,
      () => false,
    )
    if (!reachable) {
      console.error(`${what} is not answering at ${url} — is it restarting? Nothing was checked.`)
      process.exitCode = 1
      return false
    }
  }

  mkdirSync(SHOT_DIR, { recursive: true })
  const browser = await chromium.launch({ headless: !headed })

  let failed = 0
  let churn = false

  // Anything thrown outside a check body — a login that never finds its form,
  // a browser that dies — must still close the browser, or node keeps its
  // handle open and never exits. A finished run that cannot exit is
  // indistinguishable from a hung one, and costs the same to wait on.
  // One attempt at one check, in its own browser context.
  async function runCase(c) {
    // A context per check: no state bleeds between them, and a check that
    // needs a virgin history (a deep link with nothing to go Back to) gets one
    // without saying so.
    const context = await browser.newContext({
      viewport: { width: 1440, height: 900 },
      ...(existsSync(STORAGE) && !c.opts.freshLogin ? { storageState: STORAGE } : {}),
    })

    const page = await context.newPage()

    const state = { dialog: null, problems: [] }
    page.on('dialog', async d => {
      state.dialog = d.message()
      await d.accept()
    })
    page.on('pageerror', e => state.problems.push(`pageerror: ${e.message}`))
    page.on('console', m => {
      if (m.type() === 'error' && !m.text().startsWith('Failed to load resource')) {
        state.problems.push(`console.error: ${m.text().slice(0, 300)}`)
      }
    })
    page.on('response', r => {
      if (r.status() >= 500 && !isKnownDevNoise(r.url())) {
        state.problems.push(`HTTP ${r.status()}: ${r.url()}`)
      }
    })

    const driven = drivePage(page, state)
    current = { name: c.name, checks: [] }

    // storageState may be absent or stale; log in and retry once. The login
    // form gets playwright's generous default — the auth server is a separate
    // app and a cold one is slow. Only the checks are held to the fast
    // timeout, where a miss means a wrong selector rather than a slow page.
    // Where an unauthenticated browser ends up is decided by the app's auth
    // handshake, so the URL right after domcontentloaded is mid-flight and
    // says nothing yet. Settling first is what separates "needs a login" from
    // "was already on its way in".
    await page.goto(WEBAPP, { waitUntil: 'domcontentloaded' }).catch(() => {})
    await page.waitForTimeout(2500)
    if (page.url().includes('/auth/')) await login(context, page)

    context.setDefaultTimeout(ACTION_TIMEOUT)
    state.problems.length = 0
    state.dialog = null

    try {
      await c.body(driven)
    } catch (e) {
      current.checks.push({ name: 'threw', ok: false, detail: e.message })
    }

    const result = current
    result.churn = looksLikeEnvironmentChurn(state.problems)

    if (!c.opts.allowProblems && state.problems.length) {
      result.checks.push({
        name: 'no console or network errors',
        ok: false,
        detail: state.problems.slice(0, 5).join(' | '),
      })
    }

    if (result.checks.some(x => !x.ok)) {
      const shot = join(SHOT_DIR, c.name.replace(/[^a-z0-9]+/gi, '-') + '.png')
      await page.screenshot({ path: shot, fullPage: true }).catch(() => {})
      result.shot = shot
    }

    await context.close()
    return result
  }

  try {
    for (const c of picked) {
      let result = await runCase(c)

      // vite re-optimizes its dependency set the moment a new dependency shows
      // up, invalidating every module URL already served, and swaps modules
      // under a running page on each save. Either leaves the app unable to boot
      // for exactly one load. A reload is the recovery — verify-ui.mjs has
      // retried once for the same reason for as long as it has existed — so a
      // churned check gets a second attempt before it counts as a result.
      if (result.churn) {
        console.log(`RETRY ${c.name} — the app did not load; dev server churn, not a result`)
        result = await runCase(c)
        if (result.churn) churn = true
      }

      if (result.checks.some(x => !x.ok)) failed++

      // Reported as it finishes, not collected for the end: a check that hangs
      // should leave a trail of what already passed, otherwise a stuck run and a
      // slow one look identical from outside.
      report(result)
      results.push(result)
    }
  } finally {
    await browser.close().catch(() => {})
  }

  const total = results.reduce((n, r) => n + r.checks.length, 0)
  console.log(
    failed ? `\n${failed}/${results.length} checks FAILED` : `\nall ${total} assertions passed`,
  )

  if (churn) {
    console.log(
      '\nThe app failed to load for at least one check: the API was unreachable, or vite served a\n' +
        'module mid-save. Another session was editing or restarting the server while this ran.\n' +
        'These failures are NOT results — check dev_server_status and run it again.',
    )
  }

  process.exitCode = failed ? 1 : 0
  return failed === 0
}

// --- cli ---------------------------------------------------------------------

// `node drive.mjs suite.mjs --only foo` loads the suite then runs it, so a
// check file needs no boilerplate of its own.
//
// The dispatch is deferred rather than awaited here: the suite imports this
// module back, and awaiting it mid-evaluation deadlocks the two against each
// other. Letting this module finish first makes the suite's import resolve.
if (process.argv[1] && resolve(process.argv[1]) === resolve(fileURLToPath(import.meta.url))) {
  const args = process.argv.slice(2)
  const file = args.find(a => !a.startsWith('--'))
  const only = args.includes('--only') ? args[args.indexOf('--only') + 1] : null

  if (!file) {
    console.error('usage: drive.mjs SUITE.mjs [--only NAME] [--headed]')
    process.exit(1)
  }

  setImmediate(async () => {
    try {
      await import(resolve(file))
      // A suite that calls run() itself has already run by the time the import
      // resolves. Running it again here executes every check a second time and
      // adds those assertions to the totals, so the run reads as passing twice
      // as much as it checked.
      if (hasRun) {
        console.error(
          `${file} calls run() itself, so the runner did not run it again` +
            (only ? ` (--only ${only} was ignored)` : '') +
            '\nDrop the run() call from the suite: the runner is what runs it.',
        )
      } else {
        await run({ only, headed: args.includes('--headed') })
      }
    } catch (e) {
      console.error(e)
      process.exitCode = 1
    }
  })
}
