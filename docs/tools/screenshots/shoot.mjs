#!/usr/bin/env node
// Capture the docs screenshots from the local dev server, light and dark.
//
//   node docs/tools/screenshots/shoot.mjs [--only scene,scene] [--png]
//
// Logs in as the screenshot user that setup.py creates (docs@example.com),
// switches its theme through the API between passes, and writes
// docs/public/screenshots/<scene>-<theme>.webp. The theme is left on light.
// --png also keeps the lossless capture next to each WebP.
//
// Scenes are the SCENES list below: a path, a selector that means the page's
// main content has rendered, and optionally a prepare() step run after it.
import { execFileSync } from 'node:child_process'
import { mkdirSync, readFileSync, readdirSync, statSync, writeFileSync } from 'node:fs'
import { createRequire } from 'node:module'
import { dirname, join, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const HERE = dirname(fileURLToPath(import.meta.url))
const REPO = resolve(HERE, '..', '..', '..')
const AGENT = join(REPO, 'dev', 'agent')
const OUT = join(REPO, 'docs', 'public', 'screenshots')

const { stack } = await import(join(AGENT, 'stack.mjs'))
const { HUMAN_WEBAPP: WEBAPP, HUMAN_API: API } = stack()
for (const u of [WEBAPP, API]) {
  if (!['localhost', '127.0.0.1', '::1'].includes(new URL(u).hostname)) {
    console.error(`shoot.mjs is local-only; refusing ${u}`)
    process.exit(1)
  }
}

// playwright-core from the workspace's pnpm store, whatever its version
const pnpmDir = join(REPO, 'node_modules', '.pnpm')
const pwDir = readdirSync(pnpmDir).find(d => d.startsWith('playwright-core@'))
if (!pwDir) {
  console.error('playwright-core not found in node_modules/.pnpm — run pnpm install')
  process.exit(1)
}
const pwPath = join(pnpmDir, pwDir, 'node_modules', 'playwright-core')
const { chromium } = createRequire(join(pwPath, 'package.json'))(pwPath)

const args = process.argv.slice(2)
const keepPng = args.includes('--png')
const onlyArg = args[args.indexOf('--only') + 1]
const only = args.includes('--only') ? new Set(onlyArg.split(',')) : null

const EMAIL = 'docs@example.com'
const PASSWORD = readFileSync(join(AGENT, '.state', 'docs-password'), 'utf8').trim()
const VIEWPORT = { width: 1440, height: 900 }
const SCALE = Number(process.env.SHOTS_SCALE || 2)
const QUALITY = Number(process.env.SHOTS_QUALITY || 0.82)

// ── API ─────────────────────────────────────────────────────────────────────

const token = execFileSync(join(AGENT, 'token.sh'), { encoding: 'utf8' }).trim()

async function api(method, path, body) {
  const r = await fetch(API + path, {
    method,
    headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' },
    body: body && JSON.stringify(body),
  })
  const j = await r.json()
  if (j.error) throw new Error(`${method} ${path}: ${j.error.message}`)
  return j.response
}

const byHandle = (set, handle) => {
  const hit = set.find(x => x.handle === handle)
  if (!hit) throw new Error(`nothing with handle ${handle} — run setup.py first`)
  return hit
}
const value = (rec, name) => rec.values.find(v => v.name === name)?.value

const [ns] = (await api('GET', '/compose/namespace/?slug=sales&limit=1')).set
if (!ns) throw new Error('no sales namespace — run docs/tools/screenshots/setup.py first')
const nsid = ns.namespaceID
const pages = (await api('GET', `/compose/namespace/${nsid}/page/?limit=100`)).set
const deal = byHandle(
  (await api('GET', `/compose/namespace/${nsid}/module/?limit=100`)).set,
  'deal',
)
const deals = (
  await api('GET', `/compose/namespace/${nsid}/module/${deal.moduleID}/record/?limit=100`)
).set
const shownDeal = deals.find(r => value(r, 'name') === 'Analytics platform renewal')
const agent = byHandle((await api('GET', '/system/agents/?limit=0')).set, 'sales_assistant')
const chatbot = (await api('GET', '/system/chatbots/?handle=website_chat')).set[0]
const taq = byHandle(
  (await api('GET', '/automation/ng-automation/?query=notify_large_deals')).set ?? [],
  'notify_large_deals',
)
const provider = byHandle((await api('GET', '/system/llm-providers/')).set, 'docs_anthropic')
const [user] = (await api('GET', `/system/users/?email=${encodeURIComponent(EMAIL)}`)).set
if (!user) throw new Error(`no user ${EMAIL} — run setup.py first`)

// Roles the permission grid shows as columns; it restores them from localStorage.
const roles = (await api('GET', '/system/roles/?limit=500')).set
const gridRoles = ['admin', 'sales_team', 'authenticated'].map(h => {
  const r = byHandle(roles, h)
  return { mode: 'edit', ID: `edit-${r.roleID}`, roleID: r.roleID, name: [r.name] }
})

async function setTheme(theme) {
  const u = await api('GET', `/system/users/${user.userID}`)
  u.meta = { ...(u.meta || {}), theme }
  await api('PUT', `/system/users/${user.userID}`, u)
}

// ── Scenes ──────────────────────────────────────────────────────────────────

const ns_ = p => `/compose/namespace/sales${p}`
const SCENES = [
  {
    name: 'home',
    path: '/',
    ready: '.column-panel textarea, .column-panel [contenteditable]',
  },
  {
    name: 'record-list',
    path: ns_(`/pages/${byHandle(pages, 'deals').pageID}`),
    ready: 'table tbody tr',
  },
  {
    name: 'record-page',
    path: ns_(`/pages/${byHandle(pages, 'deal_record').pageID}/records/${shownDeal.recordID}`),
    ready: '.grid-stack-item table tbody tr',
  },
  {
    name: 'module-editor',
    path: ns_(`/admin/modules/${deal.moduleID}/edit`),
    ready: 'input[value="close_date"]',
  },
  {
    name: 'page-builder',
    path: ns_(`/admin/pages/${byHandle(pages, 'dashboard').pageID}/builder`),
    ready: '.grid-stack-item canvas',
  },
  {
    name: 'taq-builder',
    path: `/taq/builder/${taq.automationID}`,
    ready: '.vue-flow__node',
    // open the branch so its condition shows in the side panel
    prepare: async page => {
      await page.locator('.vue-flow__node', { hasText: 'First Match' }).click()
      await page.locator('button:has(.pi-minus)').click()
    },
  },
  {
    name: 'agent-editor',
    path: `/agentic/${agent.agentID}/edit`,
    ready: 'textarea',
  },
  {
    name: 'chatbot-editor',
    path: `/chatbot/${chatbot.chatbotID}/edit`,
    ready: 'input',
    // the journey, not the embed snippet, is what the page is about
    prepare: async page => {
      await page.getByText('Your details', { exact: true }).first().click()
      await page
        .getByText('Journey', { exact: true })
        .first()
        .evaluate(e => {
          e.scrollIntoView({ block: 'start' })
          let el = e.parentElement
          while (el && el.scrollHeight <= el.clientHeight) el = el.parentElement
          if (el) el.scrollTop -= 24
        })
    },
  },
  {
    name: 'permissions',
    path: '/admin/system/permissions',
    ready: ':text("Sales team")',
  },
].filter(s => !only || only.has(s.name))

// Toasts come and go with timing, so they never belong in a screenshot.
const HIDE_CSS = `
  .p-toast, .p-toast-message { display: none !important; }
  #__vue-devtools-container__, [id^="__vue-devtools"], #vue-inspector-container { display: none !important; }
  *, *::before, *::after { caret-color: transparent !important; }
`

// Spinners and skeletons that mean "not rendered yet".
const BUSY = '.p-progressspinner, .p-skeleton, .pi-spinner, .pi-spin, [aria-busy="true"]'

async function settle(page, ready) {
  await page.waitForLoadState('networkidle', { timeout: 30000 }).catch(() => {})
  await page.locator(ready).first().waitFor({ state: 'visible', timeout: 30000 })
  await page
    .waitForFunction(sel => ![...document.querySelectorAll(sel)].some(e => e.offsetParent), BUSY, {
      timeout: 15000,
    })
    .catch(() => console.warn('    still busy after 15s'))
  await page.waitForLoadState('networkidle', { timeout: 15000 }).catch(() => {})
  // chart and flow-canvas entry animations
  await page.waitForTimeout(1500)
}

// ── Capture ─────────────────────────────────────────────────────────────────

async function login(context) {
  const page = await context.newPage()
  await page.goto(WEBAPP, { waitUntil: 'domcontentloaded' })
  await page.waitForURL(/\/auth\//, { timeout: 20000 }).catch(() => {})
  if (page.url().includes('/auth/')) {
    await page.fill('input[name="email"]', EMAIL)
    if (!(await page.locator('input[name="password"]').count()))
      await page.click('button[type="submit"]')
    await page.fill('input[name="password"]', PASSWORD)
    await page.click('button[type="submit"]')
    await page.waitForURL(u => !u.href.includes('/auth/'), { timeout: 20000 })
  }
  return page
}

// PNG → WebP with the browser's own encoder, so nothing extra is installed.
async function toWebp(page, png) {
  const dataUrl = await page.evaluate(
    async ({ src, q }) => {
      const img = new Image()
      img.src = src
      await img.decode()
      const c = document.createElement('canvas')
      c.width = img.naturalWidth
      c.height = img.naturalHeight
      c.getContext('2d').drawImage(img, 0, 0)
      return c.toDataURL('image/webp', q)
    },
    { src: `data:image/png;base64,${png.toString('base64')}`, q: QUALITY },
  )
  return Buffer.from(dataUrl.split(',')[1], 'base64')
}

mkdirSync(OUT, { recursive: true })
const browser = await chromium.launch({ headless: true })
const encoder = await browser.newPage()
const problems = []

try {
  for (const theme of ['light', 'dark']) {
    await setTheme(theme)
    const context = await browser.newContext({
      viewport: VIEWPORT,
      deviceScaleFactor: SCALE,
      reducedMotion: 'reduce',
      colorScheme: theme,
      locale: 'en-US',
      timezoneId: 'Europe/Ljubljana',
    })
    await context.addInitScript(css => {
      addEventListener('DOMContentLoaded', () => {
        const s = document.createElement('style')
        s.textContent = css
        document.head.appendChild(s)
      })
    }, HIDE_CSS)
    await context.addInitScript(rr => {
      try {
        localStorage.setItem('permissionList.roles', JSON.stringify(rr))
        // the browser's own notification prompt, which a headless run cannot answer
        localStorage.setItem('notificationsDesktopPromptDismissed', 'true')
      } catch {}
    }, gridRoles)
    // The placeholder key cannot list the provider's models, and the model
    // picker shows nothing it has no option for.
    await context.route(`**/system/llm-providers/${provider.llmProviderID}/models`, route =>
      route.fulfill({
        contentType: 'application/json',
        body: JSON.stringify({
          response: ['claude-sonnet-4-5', 'claude-opus-4-1', 'claude-haiku-4-5'],
        }),
      }),
    )
    const page = await login(context)
    page.on('pageerror', e => problems.push(`${theme} ${page.url()}: ${e.message}`))

    for (const scene of SCENES) {
      process.stdout.write(`  ${scene.name}-${theme} `)
      await page.goto(WEBAPP + scene.path, { waitUntil: 'domcontentloaded' })
      await settle(page, scene.ready)
      if (scene.prepare) {
        await scene.prepare(page)
        await settle(page, scene.ready)
      }
      await page.mouse.move(0, VIEWPORT.height - 1)
      const png = await page.screenshot({ animations: 'disabled' })
      const webp = await toWebp(encoder, png)
      const file = join(OUT, `${scene.name}-${theme}.webp`)
      writeFileSync(file, webp)
      if (keepPng) writeFileSync(file.replace(/\.webp$/, '.png'), png)
      console.log(`${Math.round(webp.length / 1024)} KB`)
    }
    await context.close()
  }
} finally {
  await setTheme('light').catch(e => console.error(`could not restore light theme: ${e.message}`))
  await browser.close()
}

const total = readdirSync(OUT)
  .filter(f => f.endsWith('.webp'))
  .reduce((n, f) => n + statSync(join(OUT, f)).size, 0)
console.log(
  `✓ ${SCENES.length * 2} screenshots in ${OUT} — ${(total / 1024 / 1024).toFixed(2)} MB total`,
)
if (problems.length) console.warn(`page errors:\n  ${problems.join('\n  ')}`)
