import { expect, test, type Page } from '@playwright/test'

// Contract under test:
// - views/system/User/Editor.intent.md: `system:user` manual scripts render in
//   two slots — `infoFooter` under the basic information panel and
//   `passwordFooter` under the password fields — and only while editing.
// - views/system/Role/Editor.intent.md: `system:role` scripts render in the
//   editor's toolbar, beside Archive and Clone permissions.
// - views/system/UserGroup/Editor.intent.md: `system:user-group` scripts render
//   under the group's basic information.
// - views/Dashboard.intent.md: `system` scripts render above the stat cards.
// - lib/vue CManualScriptButtons is the one component all four use; a click
//   dispatches the script on `$ScriptBus` with the edited resource as subject.
//
// The scripts are the fixture extension's `client-scripts/admin/*` bundle, each
// declaring the app/page/slot it belongs to.
//
// What this would have caught: a slot lookup that ignored the `slot` uiProp, so
// both user-editor scripts rendered in both slots; a `uiPage` typo, which makes
// a button vanish with no error anywhere; and a dispatch built from the wrong
// resource, which throws inside the bundle and surfaces only as a page error and
// an error toast.
//
// The gate is the script inventory, not any one slot: it says whether Corredor
// is serving the admin bundle at all. Gating on a button instead would let a
// slot that stopped rendering skip the file it is supposed to fail.

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:5173'
const STORAGE_STATE = 'e2e/.auth/state.json'

const DASHBOARD_SCRIPT = 'Dashboard hello (client)'

/** A crashed bundle still leaves the last frame on screen; the console is the
 *  only witness. */
function watchForCrashes(page: Page): string[] {
  const errors: string[] = []
  page.on('pageerror', e => errors.push(String(e)))
  return errors
}

/** Navigates to an admin list and says whether it has anything to edit. */
async function listHasRows(page: Page, listPath: string): Promise<boolean> {
  await page.goto(listPath)
  return page
    .locator('tbody tr')
    .first()
    .waitFor({ state: 'visible', timeout: 30000 })
    .then(() => true)
    .catch(() => false)
}

/** Opens the first row of an already-loaded admin list; the id is never assumed. */
async function openFirstRow(page: Page, listPath: string) {
  await page.locator('tbody tr').first().click()
  await page.waitForURL(new RegExp(`${listPath}/\\d+`), { timeout: 30000 })
}

/** Clicks a slot button and asserts the dispatch left no wreckage: a script that
 *  throws is reported as an error toast, and a broken bundle as a page error. */
async function clickWithoutFallout(page: Page, label: string, errors: string[]) {
  const button = page.getByRole('button', { name: label })
  await button.click()
  await expect(button).toBeVisible()
  // These scripts only log; any alert here is the failure handler's toast.
  await expect(page.getByRole('alert')).toHaveCount(0)
  expect(errors).toEqual([])
}

test.describe('corredor manual script slots in admin', () => {
  let bundleMissing = ''

  test.beforeAll(async ({ browser }) => {
    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()
      await page.goto('/admin/automation/scripts')
      await expect(page.getByText('Search query')).toBeVisible({ timeout: 30000 })

      if (await page.getByText(/Corredor is turned off on this server/).count()) {
        bundleMissing =
          'Corredor is off on this stack (CORREDOR_ENABLED) — bring the stack up with Corredor: dev/agent/worktree.sh up'
        return
      }

      if (await page.getByText(/cannot reach it/).count()) {
        bundleMissing =
          'the server cannot reach Corredor on this stack — check the Corredor gRPC address (CORREDOR_ADDR)'
        return
      }

      // The inventory lists what the server loaded; a script listed here and
      // missing from its slot is the defect these tests are for.
      const listed = await page
        .getByText(DASHBOARD_SCRIPT)
        .first()
        .waitFor({ state: 'visible', timeout: 20000 })
        .then(() => true)
        .catch(() => false)

      if (!listed) {
        bundleMissing =
          "the fixture extension's admin client scripts are not deployed on this stack — Corredor serves dev/fixtures/corredor when the stack comes up with dev/agent/worktree.sh up"
      }
    } finally {
      await context.close()
    }
  })

  test('the admin dashboard renders its toolbar script', async ({ page }) => {
    test.skip(!!bundleMissing, bundleMissing)
    const errors = watchForCrashes(page)

    await page.goto('/admin')
    await expect(page.getByRole('button', { name: DASHBOARD_SCRIPT })).toBeVisible({
      timeout: 30000,
    })

    await clickWithoutFallout(page, DASHBOARD_SCRIPT, errors)
  })

  test('the user editor puts each script in the slot it asked for', async ({ page }) => {
    test.skip(!!bundleMissing, bundleMissing)
    const errors = watchForCrashes(page)

    expect(await listHasRows(page, '/admin/system/users')).toBe(true)
    await openFirstRow(page, '/admin/system/users')
    await expect(page.locator('input[name="email"]')).toBeVisible({ timeout: 30000 })

    const info = page.getByRole('region', { name: 'Basic information' })
    const security = page.getByRole('region', { name: 'Security' })

    // Two scripts, same resource and page, different slots — each belongs to one
    // panel and not the other.
    await expect(info.getByRole('button', { name: 'Greet user (client)' })).toBeVisible()
    await expect(info.getByRole('button', { name: 'Password hint (client)' })).toHaveCount(0)
    await expect(security.getByRole('button', { name: 'Password hint (client)' })).toBeVisible()
    await expect(security.getByRole('button', { name: 'Greet user (client)' })).toHaveCount(0)

    await clickWithoutFallout(page, 'Greet user (client)', errors)
    await clickWithoutFallout(page, 'Password hint (client)', errors)
  })

  test('the user create form carries no slot buttons', async ({ page }) => {
    test.skip(!!bundleMissing, bundleMissing)

    await page.goto('/admin/system/users/new')
    await expect(page.locator('input[name="email"]')).toBeVisible({ timeout: 30000 })

    // Both slots are edit-only: there is no user yet to hand a script.
    await expect(page.getByRole('button', { name: 'Greet user (client)' })).toHaveCount(0)
    await expect(page.getByRole('button', { name: 'Password hint (client)' })).toHaveCount(0)
  })

  test("the role editor's toolbar renders its script", async ({ page }) => {
    test.skip(!!bundleMissing, bundleMissing)
    const errors = watchForCrashes(page)

    expect(await listHasRows(page, '/admin/system/roles')).toBe(true)
    await openFirstRow(page, '/admin/system/roles')
    const greet = page.getByRole('button', { name: 'Greet role (client)' })
    await expect(greet).toBeVisible({ timeout: 30000 })

    // The toolbar sits above the panels, beside the editor's own actions.
    await expect(page.getByRole('button', { name: 'Archive' })).toBeVisible()
    await expect(
      page.getByRole('region', { name: 'Basic information' }).getByRole('button', {
        name: 'Greet role (client)',
      }),
    ).toHaveCount(0)

    await clickWithoutFallout(page, 'Greet role (client)', errors)
  })

  test('the user group editor renders its script under the basic information', async ({ page }) => {
    test.skip(!!bundleMissing, bundleMissing)
    const errors = watchForCrashes(page)

    const hasGroups = await listHasRows(page, '/admin/system/user-groups')
    test.skip(
      !hasGroups,
      'this stack has no user group to edit — create one, or seed a fixture that does',
    )

    await openFirstRow(page, '/admin/system/user-groups')

    const info = page.getByRole('region', { name: 'Basic information' })
    await expect(info.getByRole('button', { name: 'Greet user group (client)' })).toBeVisible({
      timeout: 30000,
    })

    await clickWithoutFallout(page, 'Greet user group (client)', errors)
  })
})
