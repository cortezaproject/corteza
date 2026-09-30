import { expect, test, type Page } from '@playwright/test'

// Contract under test: sections/admin/views/system/DataSource/Editor.intent.md
// (create, then DAL wiring on the saved connection). The data source is made
// through the editor and removed again in `afterAll`.

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:5173'
const STORAGE_STATE = 'e2e/.auth/state.json'
const STAMP = Date.now()
const NAME = `e2e-data-source-${STAMP}`
const HANDLE = `e2e_data_source_${STAMP}`
const DSN = 'sqlite3://file::memory:?cache=shared&mode=memory'

async function api(page: Page, fn: string, args: Record<string, unknown>) {
  await page.waitForFunction(() => (document.querySelector('#app') as any)?.__vue_app__)
  return page.evaluate(
    ({ fn, args }) =>
      (document.querySelector('#app') as any).__vue_app__.config.globalProperties.$SystemAPI[fn](
        args,
      ),
    { fn, args },
  )
}

async function save(page: Page) {
  await page.getByRole('button', { name: 'Save', exact: true }).click()
}

test.describe.serial('data source editor', () => {
  let connectionID = ''

  test.afterAll(async ({ browser }) => {
    if (!connectionID) return
    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()
      await page.goto('/')
      await api(page, 'dalConnectionDelete', { connectionID })
    } catch {
      // Teardown is best-effort: an e2e- data source left behind is noise, not a failure.
    } finally {
      await context.close()
    }
  })

  test('create sends no DAL wiring and lands on the editor with the DSN type offered', async ({
    page,
  }) => {
    await page.goto('/admin/system/data-sources/new')
    await page.locator('input#name').fill(NAME)
    await page.locator('input#handle').fill(HANDLE)

    const created = page.waitForRequest(
      r => r.method() === 'POST' && /\/dal\/connections\/?$/.test(r.url()),
    )
    await save(page)
    expect((await created).postDataJSON().config?.dal).toBeUndefined()

    await page.waitForURL(/\/admin\/system\/data-sources\/\d+$/, { timeout: 20000 })
    connectionID = page.url().split('/').pop() as string

    await expect(page.locator('input#dalType')).toHaveValue('corteza::dal:connection:dsn')
  })

  test('DSN, ownership and name survive a save and a reload', async ({ page }) => {
    await page.goto(`/admin/system/data-sources/${connectionID}`)
    await expect(page.locator('input#dalType')).toBeVisible()

    await page.locator('input#name').fill(`${NAME} edited`)
    await page.locator('input#ownership').fill('E2E Ops')
    await page.locator('textarea#dalParams').fill(JSON.stringify({ dsn: DSN }))
    await page.locator('textarea#dalParams').blur()

    const updated = page.waitForResponse(
      r => r.request().method() === 'PUT' && r.url().includes(`/dal/connections/${connectionID}`),
    )
    await save(page)
    await updated

    await page.reload()
    await expect(page.locator('input#name')).toHaveValue(`${NAME} edited`)
    await expect(page.locator('input#ownership')).toHaveValue('E2E Ops')
    await expect(page.locator('input#dalType')).toHaveValue('corteza::dal:connection:dsn')
    await expect(page.locator('textarea#dalParams')).toHaveValue(
      new RegExp(DSN.replace(/[?]/g, '\\?')),
    )
  })
})
