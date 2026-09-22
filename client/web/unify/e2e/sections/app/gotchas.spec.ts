import { expect, test, type Page } from '@playwright/test'
import { BRIDGE_SCRIPT } from '../../../src/sections/app/bridge'

// Contract under test: sections/app/app.intent.md, Bridge and Sandbox. One
// test per gotcha a custom app meets — each is something the store or the
// sandbox does that an app author would otherwise find out in production.
//
// The probes run against a namespace this spec creates through the webapp's
// own API clients (custom apps have no editor to drive), with an `e2e_` slug,
// removed again in `afterAll`.

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:5173'
const STORAGE_STATE = 'e2e/.auth/state.json'
const STAMP = Date.now()
const SLUG = `e2e_custom_apps_${STAMP}`
const NAME = `e2e-custom-apps-${STAMP}`

// Everything a probe app finds out, printed as JSON into #out.
const PROBE = `
<pre id="out"></pre>
<script>
window.SAMPLE = {
  'records.list': () => Promise.resolve({
    records: [{ recordID: 's1', values: { name: 'sample row', flag: false } }],
    refs: {}, nextPageCursor: null
  })
}
</script>
<script>__BRIDGE__</script>
<script>
(async () => {
  const settle = p => p.then(ok => ({ ok }), e => ({ error: String(e && e.message || e) }))
  const out = { live: await human.ready }
  out.list = await settle(human.records.list({ module: 'e2e_probes', limit: 10 }))
  out.undeclared = await settle(human.records.list({ module: 'e2e_companies' }))
  out.create = await settle(human.call('records.create', { module: 'e2e_probes', values: {} }))
  out.fetch = await settle(fetch('/api/system/users/').then(() => 'reached'))
  try { localStorage.setItem('probe', '1'); out.storage = 'works' } catch (e) { out.storage = e.name }
  document.getElementById('out').textContent = JSON.stringify(out)
})()
</script>`

const probeSource = (version: number) =>
  PROBE.replace('__BRIDGE__', BRIDGE_SCRIPT.replace(/v: \d+ \}/, `v: ${version} }`))

type Probe = {
  live: boolean
  list: { ok?: { records: any[]; refs: Record<string, string> }; error?: string }
  undeclared: { error?: string }
  create: { error?: string }
  fetch: { ok?: string; error?: string }
  storage: string
}

async function runProbe(page: Page, applicationID: string): Promise<Probe> {
  await page.goto(`/app/${applicationID}`)
  const out = page
    .frameLocator('iframe[srcdoc]')
    .first()
    .frameLocator('iframe[sandbox="allow-scripts"]')
    .first()
    .locator('#out')
  await expect(out).not.toBeEmpty({ timeout: 20000 })
  return JSON.parse(await out.innerText())
}

const byName = (probe: Probe, name: string) =>
  probe.list.ok!.records.find(r => r.values.name === name)

test.describe.serial('custom app gotchas', () => {
  let namespaceID = ''
  const apps: Record<string, string> = {}

  test.beforeAll(async ({ browser }) => {
    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()
      await page.goto('/')
      await page.waitForFunction(() => (document.querySelector('#app') as any)?.__vue_app__)

      const made = await page.evaluate(
        async ({ slug, name, sources }) => {
          const app = (document.querySelector('#app') as any).__vue_app__
          const { $ComposeAPI: compose, $SystemAPI: system } = app.config.globalProperties

          const ns = await compose.namespaceCreate({ name, slug, enabled: true, meta: {} })
          const namespaceID = ns.namespaceID

          const companies = await compose.moduleCreate({
            namespaceID,
            name: `${name} companies`,
            handle: 'e2e_companies',
            meta: {},
            fields: [{ name: 'name', kind: 'String', label: 'Name' }],
          })
          const probes = await compose.moduleCreate({
            namespaceID,
            name: `${name} probes`,
            handle: 'e2e_probes',
            meta: {},
            fields: [
              { name: 'name', kind: 'String', label: 'Name' },
              { name: 'flag', kind: 'Bool', label: 'Flag' },
              { name: 'amount', kind: 'Number', label: 'Amount', options: { precision: 2 } },
              { name: 'tags', kind: 'String', label: 'Tags', isMulti: true },
              {
                name: 'company',
                kind: 'Record',
                label: 'Company',
                options: { moduleID: companies.moduleID, labelField: 'name' },
              },
            ],
          })

          const acme = await compose.recordCreate({
            namespaceID,
            moduleID: companies.moduleID,
            values: [{ name: 'name', value: 'Acme' }],
          })
          await compose.recordCreate({
            namespaceID,
            moduleID: probes.moduleID,
            values: [
              { name: 'name', value: 'on' },
              { name: 'flag', value: '1' },
              { name: 'amount', value: '12.50' },
              { name: 'tags', value: 'a' },
              { name: 'tags', value: 'b' },
              { name: 'company', value: acme.recordID },
            ],
          })
          await compose.recordCreate({
            namespaceID,
            moduleID: probes.moduleID,
            values: [{ name: 'name', value: 'off' }],
          })

          const apps: Record<string, string> = {}
          for (const [version, source] of Object.entries(sources)) {
            const created = await system.applicationCreate({
              name: `${name}-${version}`,
              enabled: true,
              unify: { name: `${name}-${version}`, listed: false, url: '', kind: 'custom' },
            })
            await system.applicationSourceSet({
              applicationID: created.applicationID,
              source,
              namespace: slug,
              modules: ['e2e_probes'],
            })
            apps[version] = created.applicationID
          }

          return { namespaceID, apps }
        },
        { slug: SLUG, name: NAME, sources: { v1: probeSource(1), v2: probeSource(2) } },
      )

      namespaceID = made.namespaceID
      Object.assign(apps, made.apps)
    } finally {
      await context.close()
    }
  })

  test.afterAll(async ({ browser }) => {
    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()
      await page.goto('/')
      await page.waitForFunction(() => (document.querySelector('#app') as any)?.__vue_app__)
      await page.evaluate(
        async ({ namespaceID, apps }) => {
          const { $ComposeAPI: compose, $SystemAPI: system } = (
            document.querySelector('#app') as any
          ).__vue_app__.config.globalProperties
          for (const applicationID of Object.values(apps)) {
            await system.applicationDelete({ applicationID }).catch(() => {})
          }
          if (namespaceID) await compose.namespaceDelete({ namespaceID }).catch(() => {})
        },
        { namespaceID, apps },
      )
    } catch {
      // Teardown is best-effort: an e2e_ namespace left behind is noise, not a failure.
    } finally {
      await context.close()
    }
  })

  test('contract 2 types a Bool and a Number, and a false Bool is present', async ({ page }) => {
    const probe = await runProbe(page, apps.v2)
    expect(probe.live).toBe(true)
    expect(byName(probe, 'on').values).toMatchObject({ flag: true, amount: 12.5, tags: ['a', 'b'] })
    // The store keeps a false Bool as nothing at all.
    expect(byName(probe, 'off').values).toEqual({ name: 'off', flag: false })
  })

  test('contract 1 keeps the strings it was written against', async ({ page }) => {
    const probe = await runProbe(page, apps.v1)
    // The store writes a number back in its own form: 12.50 comes out as "12.5".
    expect(byName(probe, 'on').values).toMatchObject({
      flag: '1',
      amount: '12.5',
      tags: ['a', 'b'],
    })
    expect(byName(probe, 'off').values).toEqual({ name: 'off' })
  })

  test('a reference arrives labelled, in every contract', async ({ page }) => {
    for (const version of ['v1', 'v2']) {
      const probe = await runProbe(page, apps[version])
      const on = byName(probe, 'on')
      expect(probe.list.ok!.refs[on.values.company]).toBe('Acme')
      expect(probe.list.ok!.refs[on.ownedBy]).toBeTruthy()
    }
  })

  test('a module the app did not declare is refused by name', async ({ page }) => {
    const probe = await runProbe(page, apps.v2)
    expect(probe.undeclared.error).toBe('module "e2e_companies" is not declared for this app')
  })

  test('writes are not part of the contract', async ({ page }) => {
    const probe = await runProbe(page, apps.v2)
    expect(probe.create.error).toBe('operation "records.create" is not available to an app')
  })

  test('the page reaches neither the network nor storage', async ({ page }) => {
    const probe = await runProbe(page, apps.v2)
    expect(probe.fetch.ok).toBeUndefined()
    expect(probe.fetch.error).toBeTruthy()
    expect(probe.storage).toBe('SecurityError')
  })

  test('outside Human the page falls back to its own sample data', async ({ page }) => {
    await page.setContent(probeSource(2))
    const out = page.locator('#out')
    await expect(out).not.toBeEmpty({ timeout: 5000 })
    const probe: Probe = JSON.parse(await out.innerText())
    expect(probe.live).toBe(false)
    expect(probe.list.ok!.records[0].values.name).toBe('sample row')
  })
})
