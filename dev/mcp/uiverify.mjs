// Browser driver for dev_ui_verify.
//
// Lives in JS because playwright does, and is invoked by the Go tool with one
// JSON argument on argv[2], answering with one JSON object on stdout. Anything
// diagnostic goes to stderr — stdout is the protocol.
//
// It logs in as the agent-dev user with the password dev/agent/bootstrap.sh
// caches in .state/ui-password, so the browser acts as the same identity the
// rest of the toolkit uses. The session is reused across calls via a saved
// storage state; a stale one is detected by being bounced back to /auth/.
import { createRequire } from 'node:module'
import { mkdirSync, existsSync, readFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { pathToFileURL } from 'node:url'

const input = JSON.parse(process.argv[2] || '{}')
const root = input.root

// playwright is a dependency of the webapp workspace, not of this directory,
// and ESM resolves imports relative to the importing FILE rather than the
// working directory. Resolving through the workspace's own package.json is what
// lets this script live in dev/ without adding a node package there.
const webapp = join(root, 'client', 'web', 'unify', 'package.json')
const requireFromWebapp = createRequire(webapp)
//
// @playwright/test is CommonJS, so importing it from ESM wraps the exports in
// .default — reading chromium off the namespace directly yields undefined.
const playwright = await import(pathToFileURL(requireFromWebapp.resolve('@playwright/test')).href)
const chromium = playwright.chromium ?? playwright.default?.chromium
if (!chromium) {
  console.log(JSON.stringify({ error: 'could not load playwright from the webapp workspace' }))
  process.exit(0)
}
const baseURL = input.baseURL || 'http://localhost:5173'
const statePath = join(root, 'dev', 'mcp', '.state', 'ui-session.json')

const out = {
  url: null,
  title: null,
  consoleErrors: [],
  failedRequests: [],
  steps: [],
  screenshot: null,
}

function fail (message, hint) {
  console.log(JSON.stringify({ error: message, hint }))
  process.exit(0)
}

function credentials () {
  const passwordFile = join(root, 'dev', 'agent', '.state', 'ui-password')
  if (!existsSync(passwordFile)) {
    fail(
      'no UI password for the agent user',
      'run dev/agent/bootstrap.sh — it provisions one and caches it in dev/agent/.state/ui-password',
    )
  }

  // agent-ui, not agent — bootstrap.sh provisions a SEPARATE user for browser
  // login (dev/agent/bootstrap.sh:84). The agent-dev user authenticates by
  // client_credentials and has no usable password, so logging in as it fails
  // with "invalid username and password combination".
  return {
    email: input.email || 'agent-ui@local.dev',
    password: readFileSync(passwordFile, 'utf8').trim(),
  }
}

// settle waits for the app to actually render.
//
// networkidle is not enough: it fires during the lull while the SPA boots, so
// the page is still showing its splash. Screenshotting there produced a
// convincing picture of an app that "never boots" — it had simply not been
// given the three seconds it takes. The header is what auth.setup.ts waits for
// too, and the fallback keeps a page without one from failing the whole run.
async function settle (page, selector) {
  const target = selector || input.waitFor || 'header'

  try {
    await page.locator(target).first().waitFor({ state: 'visible', timeout: 15000 })
  } catch {
    out.consoleErrors.push(`(driver) never saw "${target}"; the page may not have finished rendering`)
    await page.waitForTimeout(2000)
  }
}

async function login (page) {
  const { email, password } = credentials()

  await page.goto(baseURL + '/')
  await page.waitForURL(/\/auth\//, { timeout: 15000 })
  await page.fill('input[name="email"]', email)
  await page.fill('input[name="password"]', password)
  await page.click('button[type="submit"]')

  // Back in the app once the auth server hands the session over.
  await page.waitForURL(url => !url.pathname.includes('/auth/'), { timeout: 20000 })
}

// Launch inside the guard: a missing browser binary is the most likely failure
// on a fresh checkout, and it deserves an answer the caller can act on rather
// than an uncaught exception.
let browser
try {
  browser = await chromium.launch()
} catch (e) {
  const missing = String(e.message || e).includes("Executable doesn't exist")
  fail(
    missing ? 'the chromium browser is not installed' : String(e.message || e).slice(0, 300),
    missing
      ? 'run: cd client/web/unify && npx playwright install chromium (this is what make setup does)'
      : undefined,
  )
}

try {
  const hasState = existsSync(statePath)
  const context = await browser.newContext(hasState ? { storageState: statePath } : {})
  const page = await context.newPage()

  // Collected for the whole run: a console error explains a click that
  // "did nothing" far better than the final URL does.
  page.on('console', msg => {
    if (msg.type() === 'error') out.consoleErrors.push(msg.text().slice(0, 500))
  })
  page.on('requestfailed', req => {
    out.failedRequests.push(`${req.method()} ${req.url()} — ${req.failure()?.errorText || 'failed'}`)
  })
  page.on('response', res => {
    if (res.status() >= 400) out.failedRequests.push(`${res.status()} ${res.request().method()} ${res.url()}`)
  })

  await page.goto(baseURL + (input.path || '/'), { waitUntil: 'networkidle', timeout: 30000 })

  // A saved session that has expired lands back on the auth server; logging in
  // again is cheaper than making the caller notice and retry.
  if (page.url().includes('/auth/')) {
    await login(page)
    await page.goto(baseURL + (input.path || '/'), { waitUntil: 'networkidle', timeout: 30000 })
  }

  await settle(page)

  for (const step of input.steps || []) {
    const target = page.locator(step.selector).first()

    try {
      if (step.action === 'fill') {
        await target.fill(step.value ?? '', { timeout: 10000 })
      } else {
        await target.click({ timeout: 10000 })
      }
      // urlAfter is read AFTER the page settles, not immediately: a click that
      // navigates has not navigated yet when click() resolves, so recording it
      // straight away reports the URL the caller already knew.
      await page.waitForLoadState('networkidle', { timeout: 15000 }).catch(() => {})
      await settle(page)

      // A short settle after that, because the default wait target (the app
      // header) is present on every page: it is satisfied instantly by the
      // page we came FROM, so a client-side route change is still in flight
      // when the URL is read. Without this, a click that navigated correctly
      // reports the URL it started on — which reads as "the click did nothing".
      await page.waitForTimeout(1000)

      out.steps.push({ ...step, ok: true, urlAfter: page.url() })
    } catch (e) {
      // A failed step is a result, not a crash: what the page looked like when
      // it failed is the thing the caller wants.
      out.steps.push({ ...step, ok: false, error: String(e.message || e).slice(0, 300) })
      break
    }
  }

  await page.waitForTimeout(500)

  out.url = page.url()
  out.title = await page.title()

  if (input.screenshot !== false) {
    const shot = join(root, 'dev', 'mcp', '.state', 'ui-verify.png')
    mkdirSync(dirname(shot), { recursive: true })
    await page.screenshot({ path: shot, fullPage: false })
    out.screenshot = shot
  }

  mkdirSync(dirname(statePath), { recursive: true })
  await context.storageState({ path: statePath })

  console.log(JSON.stringify(out))
} catch (e) {
  console.log(JSON.stringify({ ...out, error: String(e.message || e).slice(0, 500) }))
} finally {
  await browser.close()
}
