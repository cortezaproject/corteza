import { expect, test, type Page } from '@playwright/test'

// Contract under test:
// - components/Translator/Translator.intent.md: one dialog serves the whole
//   section; a key is named through the config's `keyPrettifier`; `highlight`
//   marks and scrolls to a row, matched on resource AND key.
// - components/Admin/Admin.intent.md: a wrapper's `update:*` emit carries the
//   TYPED resource — handing back a plain object crashed the editor on save.
// - views/Admin/Modules/Edit.intent.md: one translate button per field row,
//   covering every key that field has.
// - views/Admin/Pages/Edit.intent.md: a layout row opens the page's translator
//   aimed at that layout.
//
// Every assertion here stands for a defect that shipped: field dialogs that
// opened empty (the server loaded the module without its fields), dialogs that
// listed raw dotted paths, and a save that replaced the editor's resource with
// a plain object and threw.
//
// Setup/teardown are hooks, not tests: inside a serial block a failing test
// skips the rest, and cleanup has to survive that.

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:5173'
const STORAGE_STATE = 'e2e/.auth/state.json'

const SLUG = `e2e_translations_${Date.now()}`
const NAME = `e2e-translations-${Date.now()}`

/** Rows the translation dialog is showing, by their Key column. */
async function dialogKeys(page: Page): Promise<string[]> {
  const dialog = page.getByRole('dialog')
  await expect(dialog.locator('table')).toBeVisible()
  return dialog.locator('tbody tr:not(.bg-emphasis) td:first-child').allInnerTexts()
}

/** Page errors are what the class-loss defect surfaced as; a crashed editor
 *  still renders its last frame, so the console is the only witness. */
function watchForCrashes(page: Page): string[] {
  const errors: string[] = []
  page.on('pageerror', e => errors.push(String(e)))
  return errors
}

test.describe.serial('compose resource translations', () => {
  let namespacePath = ''
  let namespaceSlug = ''
  let modulePath = ''
  let pagePath = ''

  test.beforeAll(async ({ browser }) => {
    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()

      // Addressed by the control's own id, not by role-and-order: the shell's
      // resource search sits before the form in the DOM, so `.first()` types
      // the name into the sidebar and leaves the form's required field empty.
      await page.goto('/compose/namespaces/create')
      await page.locator('input#name').fill(NAME)
      await page.locator('input#slug').fill(SLUG)
      await page.getByRole('button', { name: 'Save', exact: true }).first().click()
      // Create lands on the namespace's own pages route; the slug it actually
      // stored is in that URL, so it is read rather than assumed.
      await page.waitForURL(/\/compose\/namespace\/[^/]+/, { timeout: 20000 })
      namespaceSlug = new URL(page.url()).pathname.split('/')[3]
      namespacePath = `/compose/namespace/${namespaceSlug}`

      await page.goto(`${namespacePath}/admin/modules/create`)
      await page.locator('input#name').fill(NAME)
      await page.getByRole('button', { name: 'Add new field' }).click()
      // The new row's name and label inputs are the first two textboxes in it.
      const row = page.locator('[data-drag-item]').first()
      await row.getByRole('textbox').nth(0).fill('probe')
      await row.getByRole('textbox').nth(1).fill('Probe')
      await page.getByRole('button', { name: 'Save', exact: true }).first().click()
      await page.waitForURL(/\/admin\/modules\/\d+\/edit/, { timeout: 20000 })
      modulePath = new URL(page.url()).pathname

      // A page, so the layout-row test has a layout to aim at. The server
      // seeds every new page a `primary` layout.
      await page.goto(`${namespacePath}/admin/pages/create`)
      await page.locator('input#title').fill(NAME)
      await page.getByRole('button', { name: 'Save', exact: true }).first().click()
      await page.waitForURL(/\/admin\/pages\/\d+\/(edit|builder)/, { timeout: 20000 })
      pagePath = new URL(page.url()).pathname.replace(/\/builder$/, '/edit')
    } finally {
      await context.close()
    }
  })

  test.afterAll(async ({ browser }) => {
    if (!namespacePath) return
    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()
      await page.goto(`/compose/namespaces/edit/${namespaceSlug}`)
      const remove = page.getByRole('button', { name: 'Delete', exact: true }).first()
      await remove.click()
      // CInputDelete asks before it deletes.
      const confirm = page.getByRole('button', { name: /^(delete|yes|confirm)$/i }).last()
      await confirm.click()
      await page.waitForURL(/\/compose\/namespaces$/, { timeout: 20000 })
    } catch {
      // Teardown is best-effort: a namespace left behind is noise, not a failure.
    } finally {
      await context.close()
    }
  })

  test('a field row translates every key that field has', async ({ page }) => {
    await page.goto(modulePath)

    // The row has to be on screen before anything is counted: `count()` is a
    // snapshot, and the app is still booting when `goto` resolves — a spec that
    // counts too early skips itself instead of failing.
    const row = page.locator('[data-drag-item]').first()
    await expect(row).toBeVisible({ timeout: 20000 })

    const translate = row.locator('button:has(.pi-language)')
    // Resource translations are off unless the instance has >1 language.
    test.skip(
      (await translate.count()) === 0,
      'resource translations are off on this stack (one language)',
    )

    await translate.first().click()

    const keys = await dialogKeys(page)
    // The field's own keys, and named — not `Meta › Description › View`.
    expect(keys.length).toBeGreaterThan(0)
    expect(keys.some(k => k.includes('›'))).toBe(false)
    expect(keys).toContain('Field label')
  })

  test('saving a translation leaves the editor alive', async ({ page }) => {
    const errors = watchForCrashes(page)
    await page.goto(modulePath)

    await expect(page.locator('[data-drag-item]').first()).toBeVisible({ timeout: 20000 })
    const translate = page.locator('button:has(.pi-language)').first()
    test.skip((await translate.count()) === 0, 'resource translations are off on this stack')
    await translate.click()

    const dialog = page.getByRole('dialog')
    const cell = dialog.locator('tbody tr').first().locator('textarea').first()
    const current = await cell.inputValue()
    await cell.fill(current === 'e2e A' ? 'e2e B' : 'e2e A')
    await dialog.getByRole('button', { name: 'Save', exact: true }).click()
    await expect(dialog).toBeHidden({ timeout: 15000 })

    // The editor holds a typed resource and calls its methods; handed a plain
    // object it throws on the next render.
    expect(errors).toEqual([])
    await expect(page.locator('[data-drag-item]').first()).toBeVisible()
  })

  test("a layout row opens the page's translator on its own row", async ({ page }) => {
    // Straight to the page the hooks made: the page list is a tree, not a table.
    await page.goto(pagePath)

    await expect(page.locator('[data-drag-item]').first()).toBeVisible({ timeout: 20000 })
    const rowTranslate = page.locator('[data-drag-item] button:has(.pi-language)')
    test.skip((await rowTranslate.count()) === 0, 'resource translations are off on this stack')
    await rowTranslate.first().click()

    const dialog = page.getByRole('dialog')
    // One dialog for the whole page: its own keys plus a section per layout.
    await expect(dialog.locator('tbody tr.bg-emphasis').first()).toBeVisible()

    const keys = await dialogKeys(page)
    expect(keys.some(k => k.includes('›'))).toBe(false)
    expect(keys).toContain('Page title')

    // Exactly one row marked, and it is the layout's own title row — this
    // catches a caller aiming at the wrong resource or key. The narrower case
    // the highlight was resource-blind for (two resources sharing a key name)
    // cannot be built from real page data, since a page's title key is `title`
    // and a layout's is `meta.title`; CTranslatorForm.test.ts covers it.
    const marked = dialog.locator('tbody tr.bg-highlight')
    await expect(marked).toHaveCount(1)
    await expect(marked.locator('td').first()).toHaveText('Layout title')
  })
})
