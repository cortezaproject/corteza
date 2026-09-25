import { expect, test, type Page } from '@playwright/test'

// Contract under test: the instance-wide home application set in the
// application editor (sections/admin/views/system/Application/Editor.intent.md)
// and the entry redirect it drives (sections/home/home.intent.md), including
// the warning when that application cannot be opened. The application is made
// through the API and removed again in `afterAll`.

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:5173'
const STORAGE_STATE = 'e2e/.auth/state.json'
const NAME = `e2e-home-app-${Date.now()}`

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

function homeToggle(page: Page) {
  return page
    .getByText('Home page for everyone', { exact: true })
    .locator('xpath=ancestor::div[.//input[@type="checkbox"]][1]')
    .locator('input[type="checkbox"]')
}

test.describe.serial('instance-wide home application', () => {
  test.setTimeout(120000)
  let applicationID = ''

  test.afterAll(async ({ browser }) => {
    if (!applicationID) return
    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()
      await page.goto('/?stay=1')
      const app = await api(page, 'applicationRead', { applicationID })
      await api(page, 'applicationUpdate', { ...app, unify: { ...app.unify, home: false } })
      await api(page, 'applicationDelete', { applicationID })
    } finally {
      await context.close()
    }
  })

  test('the editor makes an application the home page for everyone', async ({ page }) => {
    await page.goto('/?stay=1')
    const created = await api(page, 'applicationCreate', {
      name: NAME,
      enabled: true,
      unify: { name: NAME, listed: true, url: 'taq/' },
    })
    applicationID = created.applicationID

    await page.goto(`/admin/system/applications/${applicationID}`)
    await expect(homeToggle(page)).not.toBeChecked({ timeout: 20000 })
    await homeToggle(page).check()
    await page.getByRole('button', { name: 'Save', exact: true }).click()
    await expect(homeToggle(page)).toBeChecked()
    await expect
      .poll(async () => (await api(page, 'applicationRead', { applicationID })).unify.home)
      .toBe(true)

    await page.goto('/')
    await expect(page).toHaveURL(/\/taq/, { timeout: 20000 })
  })

  test('a home page that cannot be opened leaves the list, with a warning', async ({ page }) => {
    test.skip(!applicationID, 'the application was not created')
    await page.goto('/?stay=1')
    const app = await api(page, 'applicationRead', { applicationID })
    await api(page, 'applicationUpdate', { ...app, enabled: false })

    await page.goto('/')
    await expect(
      page.getByText(`Your home page application ${NAME} could not be opened`),
    ).toBeVisible({
      timeout: 20000,
    })
    await expect(page).toHaveURL(/\/\?homeUnavailable=/)
  })
})
