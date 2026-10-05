import { expect, test, type Page } from '@playwright/test'

// A Custom page block in the custom app sandbox: what it is told about where
// it is shown, what it may read, where it may take the viewer, and what the
// server refuses when the page is saved. Everything is created here and removed
// again in `afterAll`.

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:5173'
const STORAGE_STATE = 'e2e/.auth/state.json'
const STAMP = Date.now()
const SLUG = `e2e_custom_block_${STAMP}`

// Prints what the page found out as JSON into #out.
const PROBE = `
<pre id="out"></pre>
<button id="open">open</button>
<script>
(async () => {
  const settle = p => p.then(ok => ({ ok }), e => ({ error: String(e && e.message || e) }))
  const out = {}
  out.ctx = await settle(human.context())
  out.list = await settle(human.records.list({ module: 'e2e_task' }))
  out.undeclared = await settle(human.records.list({ module: 'e2e_other' }))
  out.modules = await settle(human.modules())
  if (out.ctx.ok && out.ctx.ok.recordID) {
    out.read = await settle(human.records.read({ module: 'e2e_task', recordID: out.ctx.ok.recordID }))
  }
  document.getElementById('open').onclick = () => {
    const first = out.list.ok.records[0]
    human.navigate({ module: 'e2e_task', recordID: first.recordID })
  }
  document.getElementById('out').textContent = JSON.stringify(out)
})()
</script>`

type Probe = {
  ctx: { ok?: Record<string, any>; error?: string }
  list: { ok?: { records: Array<{ recordID: string; values: Record<string, any> }> } }
  undeclared: { error?: string }
  read?: { ok?: { record: { values: Record<string, any> } } }
  modules: { ok?: Array<{ handle: string; fields: Array<{ name: string; multi: boolean }> }> }
}

// Creates a record and deletes it, and counts the refresh events it hears.
const ACTIONS = `
<button id="delete">delete</button>
<pre id="log"></pre>
<div id="heard">0</div>
<script>
let heard = 0
human.on('refresh', () => { document.getElementById('heard').textContent = String(++heard) })
document.getElementById('delete').onclick = async () => {
  try {
    const { record } = await human.records.create({ module: 'e2e_task', values: { title: 'doomed' } })
    await human.records.delete({ module: 'e2e_task', recordID: record.recordID })
    document.getElementById('log').textContent = 'deleted ' + record.recordID
  } catch (e) { document.getElementById('log').textContent = 'error ' + e.message }
}
</script>`

// Tells the rest of the page that data changed.
const PINGER = `<button id="ping">ping</button>
<script>document.getElementById('ping').onclick = () => human.refresh()</script>`

const inline = (extra: Record<string, unknown> = {}) => ({
  kind: 'Custom',
  title: 'Probe',
  xywh: [0, 0, 48, 20],
  options: { source: PROBE, modules: ['e2e_task'], ...extra },
})

async function withApp<T>(page: Page, fn: (apis: any, arg: any) => Promise<T>, arg?: unknown) {
  await page.waitForFunction(() => (document.querySelector('#app') as any)?.__vue_app__)
  return page.evaluate(
    async ({ src, arg }) => {
      const apis = (document.querySelector('#app') as any).__vue_app__.config.globalProperties
      // eslint-disable-next-line no-new-func
      return new Function('apis', 'arg', `return (${src})(apis, arg)`)(apis, arg)
    },
    { src: fn.toString(), arg },
  )
}

function frame(page: Page, index = 0) {
  return page
    .frameLocator('iframe[srcdoc]')
    .nth(index)
    .frameLocator('iframe[sandbox^="allow-scripts"]')
    .first()
}

async function probe(page: Page, path: string, index = 0): Promise<Probe> {
  await page.goto(path)
  const out = frame(page, index).locator('#out')
  await expect(out).not.toBeEmpty({ timeout: 20000 })
  return JSON.parse(await out.innerText())
}

test.describe.serial('custom page block', () => {
  let namespaceID = ''
  let dashboardID = ''
  let actionsID = ''
  let recordPageID = ''
  let recordID = ''

  test.beforeAll(async ({ browser }) => {
    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()
      await page.goto('/')
      const made = await withApp(
        page,
        async ({ $ComposeAPI: compose }, { slug, dashboard, recordBlock, actionBlocks }) => {
          const ns = await compose.namespaceCreate({ name: slug, slug, enabled: true, meta: {} })
          const namespaceID = ns.namespaceID
          const fields = [
            { name: 'title', kind: 'String', label: 'Title' },
            { name: 'tags', kind: 'String', label: 'Tags', isMulti: true },
          ]
          const task = await compose.moduleCreate({
            namespaceID,
            name: 'Task',
            handle: 'e2e_task',
            meta: {},
            fields,
          })
          await compose.moduleCreate({
            namespaceID,
            name: 'Other',
            handle: 'e2e_other',
            meta: {},
            fields,
          })
          const first = await compose.recordCreate({
            namespaceID,
            moduleID: task.moduleID,
            values: [{ name: 'title', value: 'first' }],
          })
          await compose.recordCreate({
            namespaceID,
            moduleID: task.moduleID,
            values: [{ name: 'title', value: 'second' }],
          })
          const dash = await compose.pageCreate({
            namespaceID,
            title: 'Dashboard',
            handle: 'dashboard',
            visible: true,
            blocks: [dashboard],
          })
          const rec = await compose.pageCreate({
            namespaceID,
            title: 'Task',
            handle: 'task_record',
            moduleID: task.moduleID,
            visible: true,
            blocks: [recordBlock],
          })
          const actions = await compose.pageCreate({
            namespaceID,
            title: 'Actions',
            handle: 'actions',
            visible: true,
            blocks: actionBlocks,
          })
          return {
            namespaceID,
            actionsID: actions.pageID,
            dashboardID: dash.pageID,
            recordPageID: rec.pageID,
            recordID: first.recordID,
          }
        },
        {
          slug: SLUG,
          dashboard: inline({ params: { colour: 'teal' } }),
          recordBlock: inline(),
          actionBlocks: [
            inline({ source: ACTIONS, writes: ['e2e_task'], deletes: ['e2e_task'] }),
            { ...inline({ source: PINGER }), xywh: [0, 20, 48, 10] },
          ],
        },
      )
      ;({ namespaceID, actionsID, dashboardID, recordPageID, recordID } = made)
    } finally {
      await context.close()
    }
  })

  test.afterAll(async ({ browser }) => {
    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()
      await page.goto('/')
      await withApp(
        page,
        async ({ $ComposeAPI: compose }, id) => {
          if (id) await compose.namespaceDelete({ namespaceID: id }).catch(() => {})
        },
        namespaceID,
      )
    } catch {
      // Teardown is best-effort: an e2e_ namespace left behind is noise, not a failure.
    } finally {
      await context.close()
    }
  })

  test('an inline block reads its page namespace and its own parameters', async ({ page }) => {
    const out = await probe(page, `/compose/namespace/${SLUG}/pages/${dashboardID}`)
    expect(out.ctx.ok).toMatchObject({ namespace: SLUG, namespaceID, pageID: dashboardID })
    expect(out.ctx.ok?.params).toEqual({ colour: 'teal' })
    expect(out.list.ok?.records.map(r => r.values.title).sort()).toEqual(['first', 'second'])
  })

  test('a module of the namespace the block did not declare is refused', async ({ page }) => {
    const out = await probe(page, `/compose/namespace/${SLUG}/pages/${dashboardID}`)
    expect(out.undeclared.error).toContain('module "e2e_other" is not declared')
  })

  test('on a record page the block is told the record it sits on', async ({ page }) => {
    const out = await probe(
      page,
      `/compose/namespace/${SLUG}/pages/${recordPageID}/records/${recordID}`,
    )
    expect(out.ctx.ok?.recordID).toBe(recordID)
    expect(out.read?.ok?.record.values.title).toBe('first')
  })

  test('the block opens a record page in the shell', async ({ page }) => {
    await probe(page, `/compose/namespace/${SLUG}/pages/${dashboardID}`)
    await frame(page).locator('#open').click()
    await expect(page).toHaveURL(new RegExp(`/pages/${recordPageID}/records/\\d+$`))
  })

  test('saving a page refuses HTML the sandbox cannot run, and writes it does not read', async ({
    page,
  }) => {
    await page.goto('/')
    const refused = await withApp(
      page,
      async ({ $ComposeAPI: compose }, { namespaceID, blocks }) => {
        const out: string[] = []
        for (const block of blocks) {
          await compose
            .pageCreate({ namespaceID, title: 'Refused', handle: '', blocks: [block] })
            .then(() => out.push('saved'))
            .catch((e: any) => out.push(String(e?.message || e)))
        }
        return out
      },
      {
        namespaceID,
        blocks: [
          inline({ source: '<script>fetch("/api")</script>' }),
          inline({ writes: ['e2e_other'] }),
          inline({ modules: ['missing'] }),
        ],
      },
    )
    expect(refused[0]).toContain("connect-src 'none'")
    expect(refused[1]).toContain('may change module "e2e_other" but does not read it')
    expect(refused[2]).toContain('there is no module "missing"')
  })

  test('modules tells a multi-value field apart', async ({ page }) => {
    const out = await probe(page, `/compose/namespace/${SLUG}/pages/${dashboardID}`)
    const fields = out.modules.ok?.find(m => m.handle === 'e2e_task')?.fields || []
    expect(fields.find(f => f.name === 'tags')?.multi).toBe(true)
    expect(fields.find(f => f.name === 'title')?.multi).toBe(false)
  })

  test('a block deletes a record of a module declared for it, once the viewer allows', async ({
    page,
  }) => {
    await page.goto(`/compose/namespace/${SLUG}/pages/${actionsID}`)
    const block = frame(page, 0)
    await block.locator('#delete').click({ timeout: 20000 })
    await page.getByRole('button', { name: 'Allow this app to make changes' }).click()
    await page.getByRole('button', { name: 'Allow this app to delete records' }).click()
    await expect(block.locator('#log')).toContainText('deleted ', { timeout: 15000 })
  })

  test('a block hears another block change data', async ({ page }) => {
    await page.goto(`/compose/namespace/${SLUG}/pages/${actionsID}`)
    const block = frame(page, 0)
    await expect(block.locator('#heard')).toHaveText('0', { timeout: 20000 })
    await frame(page, 1).locator('#ping').click()
    await expect(block.locator('#heard')).toHaveText('1')
  })
})
