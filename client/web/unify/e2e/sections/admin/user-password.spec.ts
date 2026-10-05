import { expect, test, type Page } from '@playwright/test'

// Contract under test: sections/admin/views/system/User/Editor.intent.md
// (a password is set from its own dialog in the editor's action row, write-only,
// and only once both entries match).
//
// The target user is made through the API and deleted again in `afterAll`.

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:5173'
const STORAGE_STATE = 'e2e/.auth/state.json'
const STAMP = Date.now()
const PASSWORD = 'E2e-Strong-Passw0rd!'

// Opens the user list and waits for its rows, so the API calls that follow go
// through an app that has finished signing in.
async function settled(page: Page) {
  await page.goto('/admin/system/users')
  await expect(page.locator('tbody tr').first()).toBeVisible({ timeout: 30000 })
}

async function api(page: Page, fn: string, args: Record<string, unknown>) {
  return page.evaluate(
    ({ fn, args }) =>
      (document.querySelector('#app') as any).__vue_app__.config.globalProperties.$SystemAPI[fn](
        args,
      ),
    { fn, args },
  )
}

function watchPasswordRequests(page: Page) {
  const bodies: string[] = []
  page.on('request', r => {
    if (r.method() === 'POST' && /\/system\/users\/\d+\/password$/.test(r.url())) {
      bodies.push(r.postData() || '')
    }
  })
  return bodies
}

test.describe.serial('setting a user password', () => {
  let userID = ''

  test.beforeAll(async ({ browser }) => {
    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()
      await settled(page)
      const created = await api(page, 'userCreate', {
        email: `e2e-password-${STAMP}@local.dev`,
        name: `e2e-password-${STAMP}`,
        userGroupID: '0',
      })
      userID = created.userID
    } finally {
      await context.close()
    }
  })

  test.afterAll(async ({ browser }) => {
    if (!userID) return
    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()
      await settled(page)
      await api(page, 'userDelete', { userID })
    } catch {
      // Teardown is best-effort: an e2e- user left behind is noise, not a failure.
    } finally {
      await context.close()
    }
  })

  async function openDialog(page: Page) {
    await page.goto(`/admin/system/users/${userID}`)
    await expect(page.locator('input[name="email"]')).toBeVisible({ timeout: 30000 })
    await page.getByRole('button', { name: 'Set password' }).click()
    const dialog = page.getByRole('dialog', { name: 'Set a new password' })
    await expect(dialog).toBeVisible()
    return dialog
  }

  test('the editor offers a button, not password fields in a panel', async ({ page }) => {
    await page.goto(`/admin/system/users/${userID}`)
    await expect(page.locator('input[name="email"]')).toBeVisible({ timeout: 30000 })

    await expect(page.getByRole('button', { name: 'Set password' })).toBeVisible()
    await expect(page.locator('input[type="password"]')).toHaveCount(0)
    await expect(page.getByLabel('New password')).toHaveCount(0)
  })

  test('mismatched entries are refused without reaching the server', async ({ page }) => {
    const sent = watchPasswordRequests(page)
    const dialog = await openDialog(page)

    await dialog.locator('#user-password-new').fill(PASSWORD)
    await dialog.locator('#user-password-confirm').fill('something-else')
    await dialog.getByRole('button', { name: 'Set password' }).click()

    await expect(dialog.getByText('The two passwords do not match')).toBeVisible()
    await expect(dialog).toBeVisible()
    expect(sent).toHaveLength(0)
  })

  test('cancelling discards what was typed', async ({ page }) => {
    let dialog = await openDialog(page)
    await dialog.locator('#user-password-new').fill(PASSWORD)
    await dialog.getByRole('button', { name: 'Cancel' }).click()
    await expect(dialog).toBeHidden()

    await page.getByRole('button', { name: 'Set password' }).click()
    dialog = page.getByRole('dialog', { name: 'Set a new password' })
    await expect(dialog.locator('#user-password-new')).toHaveValue('')
  })

  test('matching entries set the password at once and leave the editor clean', async ({ page }) => {
    const sent = watchPasswordRequests(page)
    const dialog = await openDialog(page)

    await dialog.locator('#user-password-new').fill(PASSWORD)
    await dialog.locator('#user-password-confirm').fill(PASSWORD)
    await dialog.getByRole('button', { name: 'Set password' }).click()

    await expect(dialog).toBeHidden()
    await expect(page.getByText('Password changed')).toBeVisible()
    expect(sent).toHaveLength(1)
    expect(JSON.parse(sent[0])).toEqual({ password: PASSWORD })

    // The password is listed among the user's sign-in methods straight away.
    await expect(
      page.getByRole('region', { name: 'Sign-in methods' }).getByText('Password'),
    ).toBeVisible()

    // Nothing waits on Save: leaving through the editor raises no unsaved-changes prompt.
    const prompts: string[] = []
    page.on('dialog', d => {
      prompts.push(d.message())
      d.dismiss()
    })
    await page.getByRole('button', { name: 'Back' }).click()
    await expect(page).toHaveURL(/\/admin\/system\/users(\?|$)/)
    expect(prompts).toEqual([])
  })

  test('a suspended user carries a tag and is offered Unsuspend', async ({ page }) => {
    await settled(page)
    await api(page, 'userSuspend', { userID })

    await page.goto(`/admin/system/users/${userID}`)
    await expect(page.locator('input[name="email"]')).toBeVisible({ timeout: 30000 })
    await expect(page.getByTestId('user-state')).toHaveText('Suspended')
    await expect(page.getByRole('button', { name: 'Unsuspend' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Suspend', exact: true })).toHaveCount(0)
  })
})

test('revoking your own sessions is disabled and says why', async ({ page }) => {
  await settled(page)
  const selfID = await page.evaluate(
    () =>
      (document.querySelector('#app') as any).__vue_app__.config.globalProperties.$Auth.user.userID,
  )

  await page.goto(`/admin/system/users/${selfID}`)
  await expect(page.locator('input[name="email"]')).toBeVisible({ timeout: 30000 })
  await expect(page.getByTestId('user-state')).toHaveCount(0)

  const revoke = page.getByRole('button', { name: 'Revoke all active sessions' })
  await expect(revoke).toBeDisabled()
  await revoke.locator('..').hover()
  await expect(page.getByText('This would end your own session as well')).toBeVisible()
})

test('a system user is not offered a password', async ({ page }) => {
  await settled(page)
  const systemID = await page.evaluate(async () => {
    const api = (document.querySelector('#app') as any).__vue_app__.config.globalProperties
      .$SystemAPI
    const { set } = await api.userList({ kind: 'sys', limit: 1 })
    return set?.[0]?.userID || ''
  })
  test.skip(!systemID, 'no system user on this stack')

  await page.goto(`/admin/system/users/${systemID}`)
  await expect(page.locator('input[name="email"]')).toBeVisible({ timeout: 30000 })
  await page.waitForLoadState('networkidle')
  await expect(page.getByRole('button', { name: 'Revoke all active sessions' })).toBeVisible()
  await expect(page.getByRole('button', { name: 'Set password' })).toHaveCount(0)
})
