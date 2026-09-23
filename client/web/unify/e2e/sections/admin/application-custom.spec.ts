import { expect, test, type Page } from '@playwright/test'

// Contract under test: sections/admin/views/system/Application/Editor.intent.md
// (custom kind, locked url, the Custom app panel) and sections/app/app.intent.md
// (a switched-off custom app previews for whoever may change its page).
//
// The application is made through the editor itself; its page is stored through
// the API, since the editor only shows one. Removed again in `afterAll`.

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:5173'
const STORAGE_STATE = 'e2e/.auth/state.json'
const NAME = `e2e-custom-app-${Date.now()}`
const PAGE = '<title>Editor probe</title><p id="hello">hello from the stored page</p>'

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

// The shell renders the route under its loading overlay and only then decides
// whether the application is switched off, so a view seen before the page has
// settled is one the disabled screen may still replace.
async function shellSettled(page: Page) {
  await page.waitForLoadState('networkidle')
  await expect(page.getByText('Application Disabled')).toHaveCount(0)
}

test.describe.serial('custom application in the editor', () => {
  let applicationID = ''

  test.afterAll(async ({ browser }) => {
    if (!applicationID) return
    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()
      await page.goto('/')
      await api(page, 'applicationDelete', { applicationID })
    } catch {
      // Teardown is best-effort: an e2e- application left behind is noise, not a failure.
    } finally {
      await context.close()
    }
  })

  test('creating a custom app gives it its own address and an empty page to write', async ({
    page,
  }) => {
    await page.goto('/admin/system/applications/new')
    await page.locator('input#name').fill(NAME)
    await page.locator('#unifyKind').getByText('Custom app').click()
    await expect(page.locator('#unifyUrl')).toBeDisabled()
    await page.getByRole('button', { name: 'Save', exact: true }).click()

    await page.waitForURL(/\/admin\/system\/applications\/\d+$/, { timeout: 20000 })
    applicationID = page.url().split('/').pop()!

    await expect(page.locator('#unifyUrl')).toHaveValue(`app/${applicationID}`)
    // With nothing stored yet, whoever may change the page is offered the
    // editor rather than told there is nothing to see.
    await expect(page.locator('#customSource .cm-content')).toBeVisible()

    // What the launcher opens is the stored url, not what the editor displays.
    const stored = await api(page, 'applicationRead', { applicationID })
    expect(stored.unify.url).toBe(`app/${applicationID}`)
  })

  test('the editor shows the stored page and keeps the kind when saved', async ({ page }) => {
    test.skip(!applicationID, 'the application was not created')
    await page.goto('/')
    // No namespace: this spec is about the editor, and a declaration naming
    // something that is not here is refused while the page is stored.
    await api(page, 'applicationSourceSet', { applicationID, source: PAGE })

    await page.goto(`/admin/system/applications/${applicationID}`)
    await expect(page.locator('#customSource .cm-content')).toContainText(
      'hello from the stored page',
    )

    await page.locator('#description').fill('saved by e2e')
    await page.getByRole('button', { name: 'Save', exact: true }).click()
    await expect(page.locator('#customSource .cm-content')).toContainText(
      'hello from the stored page',
    )

    const saved = await api(page, 'applicationRead', { applicationID })
    expect(saved.unify.kind).toBe('custom')
    expect(saved.unify.url).toBe(`app/${applicationID}`)
  })

  // Whoever may change the page does it here; what the sandbox could not run
  // is refused in the editor exactly as it is over the API.
  test('the page is edited and saved from the editor', async ({ page }) => {
    test.skip(!applicationID, 'the application was not created')
    await page.goto(`/admin/system/applications/${applicationID}`)

    const editor = page.locator('#customSource .cm-content')
    await editor.waitFor({ timeout: 20000 })
    await expect(page.locator('[data-test-id="button-save-page"]')).toBeDisabled()

    // A mark of its own, so a page left behind by an earlier run cannot pass
    // this test for it.
    const mark = `added-${Date.now()}`
    await editor.click()
    await page.keyboard.press('Control+End')
    await page.keyboard.type(`<p id="added">${mark}</p>`)
    await page.locator('[data-test-id="button-save-page"]').click()

    await expect
      .poll(async () => (await api(page, 'applicationSourceRead', { applicationID })).source)
      .toContain(mark)

    // the page's own meta keeps up with what was saved
    const after = await api(page, 'applicationSourceRead', { applicationID })
    expect(after.sourceMeta.size).toBe(after.source.length)
  })

  test('the editor refuses a page the sandbox could not run', async ({ page }) => {
    test.skip(!applicationID, 'the application was not created')
    await page.goto(`/admin/system/applications/${applicationID}`)

    const editor = page.locator('#customSource .cm-content')
    await editor.waitFor({ timeout: 20000 })
    await editor.click()
    await page.keyboard.press('Control+End')
    await page.keyboard.type('<a href="https://example.com/">leave</a>')
    await page.locator('[data-test-id="button-save-page"]').click()

    await expect(page.getByText(/could not be saved|leaves the page/)).toBeVisible()
    const after = await api(page, 'applicationSourceRead', { applicationID })
    expect(after.source).not.toContain('example.com')
  })

  test('open the app from the editor', async ({ page }) => {
    test.skip(!applicationID, 'the application was not created')
    await page.goto(`/admin/system/applications/${applicationID}`)
    await page.locator('[data-test-id="button-open-custom-app"]').click()
    await expect(page).toHaveURL(new RegExp(`/app/${applicationID}$`))
    await expect(
      page
        .frameLocator('iframe[srcdoc]')
        .first()
        .frameLocator('iframe[sandbox="allow-scripts"]')
        .first()
        .locator('#hello'),
    ).toHaveText('hello from the stored page')
  })

  test('switched off, it previews for whoever may change its page', async ({ page }) => {
    test.skip(!applicationID, 'the application was not created')
    await page.goto('/')
    const app = await api(page, 'applicationRead', { applicationID })
    await api(page, 'applicationUpdate', {
      applicationID,
      name: app.name,
      enabled: false,
      weight: app.weight,
      meta: app.meta,
      unify: app.unify,
    })

    await page.goto(`/app/${applicationID}`)
    await shellSettled(page)
    await expect(page.locator('[data-test-id="custom-app-preview"]')).toBeVisible()
    await expect(page.locator('iframe[srcdoc]')).toHaveCount(1)
  })
})
