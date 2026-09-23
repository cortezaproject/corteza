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
  out.modules = await settle(human.call('modules'))
  out.create = await settle(human.call('records.create', { module: 'e2e_probes', values: {} }))
  // Reached indirectly: storing a page that names them is refused, and what
  // these tests are about is what the sandbox does when one gets through.
  out.fetch = await settle(window['fet' + 'ch']('/api/system/users/').then(() => 'reached'))
  try { window['local' + 'Storage'].setItem('probe', '1'); out.storage = 'works' } catch (e) { out.storage = e.name }
  document.getElementById('out').textContent = JSON.stringify(out)
})()
</script>`

// Every way a page can try to load another document.
// Every way a page can try to load another document. The links that leave are
// built at run time: storing them is refused, which is the point of the guard,
// while what they do once rendered is the point of these tests.
// Every way a page can try to load another document. The addresses are put
// together at run time: storing them is refused, which is the guard's job,
// while what they do once rendered is what these tests are about.
const LINKS = `
<script>__BRIDGE__</script>
<p>
  <a id="anchor" href="#below">anchor</a>
  <span id="away"></span>
</p>
<form><button id="submit">submit</button></form>
<button id="scripted">scripted</button>
<button id="download" onclick="human.download('probe.csv', 'a,b')">download</button>
<p id="log"></p>
<script>
var elsewhere = 'htt' + 'ps://example.com/'
document.getElementById('away').innerHTML =
  '<a id="site">site</a> <a id="mail">mail</a> <a id="handled">handled</a>'
document.getElementById('site').setAttribute('href', elsewhere)
document.getElementById('mail').setAttribute('href', 'mail' + 'to:ana@example.com')
document.getElementById('handled').setAttribute('href', elsewhere)
document.getElementById('handled').onclick = () => { document.getElementById('log').textContent = 'handled' }
document.getElementById('scripted').onclick = () => { location.href = elsewhere }
</script>
<div style="height: 3000px"></div>
<p id="below">below</p>`

// A page that throws where a real one did: an object literal, then a line
// beginning with a bracket, which JavaScript reads as calling the object.
const THROWS = `
<h1>Contacts by type</h1>
<script>__BRIDGE__</script>
<script>
window.SAMPLE = { 'records.list': function () { return Promise.resolve({ records: [] }) } }

(async function () {
  document.body.appendChild(document.createTextNode('never reached'))
})()
</script>`

// What a page that changes records finds out.
const WRITES = `
<pre id="out"></pre>
<script>__BRIDGE__</script>
<script>
(async () => {
  const settle = p => p.then(ok => ({ ok }), e => ({ error: String(e && e.message || e) }))
  const out = {}
  out.create = await settle(human.records.create({ module: 'e2e_probes', values: { name: 'written by the app', flag: true } }))
  const id = out.create.ok && out.create.ok.record.recordID
  out.update = id ? await settle(human.records.update({ module: 'e2e_probes', recordID: id, values: { flag: false } })) : { error: 'no record' }
  out.readOnlyModule = await settle(human.records.create({ module: 'e2e_companies', values: { name: 'x' } }))
  document.getElementById('out').textContent = JSON.stringify(out)
})()
</script>`

const probeSource = (version: number) =>
  PROBE.replace('__BRIDGE__', BRIDGE_SCRIPT.replace(/v: \d+ \}/, `v: ${version} }`))

type Probe = {
  live: boolean
  list: { ok?: { records: any[]; refs: Record<string, string> }; error?: string }
  undeclared: { error?: string }
  modules: { ok?: { handle: string; fields: any[] }[]; error?: string }
  create: { error?: string }
  fetch: { ok?: string; error?: string }
  storage: string
}

function appFrame(page: Page) {
  return page
    .frameLocator('iframe[srcdoc]')
    .first()
    .frameLocator('iframe[sandbox="allow-scripts"]')
    .first()
}

async function runProbe(page: Page, applicationID: string): Promise<Probe> {
  await page.goto(`/app/${applicationID}`)
  const out = appFrame(page).locator('#out')
  await expect(out).not.toBeEmpty({ timeout: 20000 })
  return JSON.parse(await out.innerText())
}

const byName = (probe: Probe, name: string) =>
  probe.list.ok!.records.find(r => r.values.name === name)

let linksApp = ''

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
        async ({ slug, name, sources, writable }) => {
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
                name: 'stage',
                kind: 'Select',
                label: 'Stage',
                options: { options: [{ value: 'won', text: 'Closed won' }] },
              },
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
              // Only the writing app may read the second module; the others
              // need one they cannot touch.
              modules: version === 'writes' ? ['e2e_probes', 'e2e_companies'] : ['e2e_probes'],
              writes: version === 'writes' ? writable : [],
            })
            apps[version] = created.applicationID
          }

          return { namespaceID, apps }
        },
        {
          slug: SLUG,
          name: NAME,
          sources: {
            v1: probeSource(1),
            v2: probeSource(2),
            links: LINKS.replace('__BRIDGE__', BRIDGE_SCRIPT),
            writes: WRITES.replace('__BRIDGE__', BRIDGE_SCRIPT),
            throws: THROWS.replace('__BRIDGE__', BRIDGE_SCRIPT),
          },
          writable: ['e2e_probes'],
        },
      )

      namespaceID = made.namespaceID
      Object.assign(apps, made.apps)
      linksApp = made.apps.links
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

  // A page labels a Select from the options the module holds. Read through the
  // API those options say `text`, so an option that arrives under `label`
  // alone renders as nothing at all.
  test('a Select option carries its wording under the name the module uses', async ({ page }) => {
    const probe = await runProbe(page, apps.v2)
    const probes = probe.modules.ok!.find((m: { handle: string }) => m.handle === 'e2e_probes')
    const stage = probes.fields.find((f: { name: string }) => f.name === 'stage')
    expect(stage.options).toEqual([{ value: 'won', text: 'Closed won', label: 'Closed won' }])
  })

  // Silence is the worst answer here: the page draws its headings, stops, and
  // reads as a working app that happens to hold no data.
  test('an app that throws says so to the viewer', async ({ page }) => {
    test.skip(!apps.throws, 'the throwing app was not created')
    await page.goto(`/app/${apps.throws}`)
    const said = page.locator('[data-test-id="custom-app-failure"]')
    await expect(said).toBeVisible({ timeout: 20000 })
    await expect(said).toContainText('is not a function')
  })

  test('a module the app did not declare is refused by name', async ({ page }) => {
    const probe = await runProbe(page, apps.v2)
    expect(probe.undeclared.error).toBe('module "e2e_companies" is not declared for this app')
  })

  // An app changes records only where it was deployed to; a page that declared
  // none is refused before the viewer is asked anything.
  test('an app that declared no writable module changes nothing', async ({ page }) => {
    const probe = await runProbe(page, apps.v2)
    expect(probe.create.error).toBe(
      'module "e2e_probes" is not declared as one this app may change',
    )
    expect(page.getByRole('alertdialog')).toHaveCount(0)
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

  test('an in-page anchor scrolls, and the app keeps running', async ({ page }) => {
    test.skip(!linksApp, 'the links app was not created')
    await page.goto(`/app/${linksApp}`)
    const app = appFrame(page)
    await app.locator('#anchor').click()
    await expect
      .poll(() => app.locator('body').evaluate(() => Math.round(scrollY)))
      .toBeGreaterThan(1000)
    await expect(page.locator('iframe[srcdoc]')).toHaveCount(1)
  })

  test('a link or form that leaves the page does nothing, and the app keeps running', async ({
    page,
  }) => {
    test.skip(!linksApp, 'the links app was not created')
    for (const id of ['site', 'mail', 'handled', 'submit']) {
      await page.goto(`/app/${linksApp}`)
      const app = appFrame(page)
      await app.locator(`#${id}`).click()
      await page.waitForTimeout(1000)
      await expect(page.locator('iframe[srcdoc]'), id).toHaveCount(1)
      if (id === 'handled') await expect(app.locator('#log')).toHaveText('handled')
    }
  })

  test('a page that navigates itself by script is stopped', async ({ page }) => {
    test.skip(!linksApp, 'the links app was not created')
    await page.goto(`/app/${linksApp}`)
    await appFrame(page).locator('#scripted').click()
    await expect(page.locator('iframe[srcdoc]')).toHaveCount(0)
    await expect(
      page.getByText('This app tried to leave the sandbox and was stopped.'),
    ).toBeVisible()
  })

  // A change asks the person first: the dialog belongs to Human, outside the
  // sandbox, so the app can neither draw it nor answer it.
  async function runWrites(page: Page, answer: RegExp) {
    await page.goto(`/app/${apps.writes}`)
    const dialog = page.getByRole('alertdialog')
    await dialog.waitFor({ timeout: 20000 })
    await dialog.getByRole('button', { name: answer }).click()
    const out = appFrame(page).locator('#out')
    await expect(out).not.toBeEmpty({ timeout: 20000 })
    return JSON.parse(await out.innerText())
  }

  test('with the viewer agreeing, a record is created and changed in place', async ({ page }) => {
    test.skip(!apps.writes, 'the writes app was not created')
    const out = await runWrites(page, /Allow/)

    expect(out.create.ok.record.values).toMatchObject({ name: 'written by the app', flag: true })
    // The update names one field; the rest of the record stays as it was.
    expect(out.update.ok.record.values).toMatchObject({ name: 'written by the app', flag: false })
  })

  test('a refusal stops the change, and a module it may only read is refused', async ({ page }) => {
    test.skip(!apps.writes, 'the writes app was not created')
    const out = await runWrites(page, /Do not allow/)

    expect(out.create.error).toContain('did not agree')
    expect(out.readOnlyModule.error).toBe(
      'module "e2e_companies" is not declared as one this app may change',
    )
  })

  // The sandbox saves no file of its own; Human does it for the app.
  test('an app hands the viewer a file', async ({ page }) => {
    test.skip(!linksApp, 'the links app was not created')
    await page.goto(`/app/${linksApp}`)
    const app = appFrame(page)
    await app.locator('#download').waitFor({ timeout: 20000 })

    const [download] = await Promise.all([
      page.waitForEvent('download', { timeout: 20000 }),
      app.locator('#download').click(),
    ])
    expect(download.suggestedFilename()).toBe('probe.csv')
  })
})
