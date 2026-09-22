import { expect, test, type Locator, type Page } from '@playwright/test'

// Contract under test:
// - components/PageBlocks/PageBlocks.intent.md: a page button naming a `script`
//   goes through `$ScriptBus.Dispatch`, which runs a client script in the browser
//   and forwards a server script to the API; the record a server script returns
//   is applied to the one on the page, NOT saved.
// - views/Pages/RecordView.intent.md: the `ui:compose:record-page`
//   `beforeFormSubmit` event runs on the very record that is then saved, so a
//   script may correct it.
//
// The scripts are the fixture extension's `agent-sandbox` set against the seeded
// `agent-contact` module.
//
// What each group would have caught:
// - greet: the compose script context is what gives a client script its toast and
//   its router; handed the plain webapp context instead, `ComposeUI` was
//   undefined and the button failed silently. Ending on the record VIEWER also
//   proves `gotoRecordViewer` still reaches the router.
// - activate: the trigger endpoint answers a record payload the caller has to
//   apply — a dispatch that dropped the response left the page showing the old
//   status, and one that saved it would persist a change the author never asked
//   for. The reload asserts both halves.
// - email: `beforeFormSubmit` must run on the record that is then saved, not on a
//   clone taken before it; against a clone the script's correction was thrown
//   away and the typed value was stored as typed. Note the browser's own email
//   input strips the padding, so what the script is proved by here is the
//   lower-casing: the server trims but never lower-cases.
//
// Nothing is restored for the activate case, because the contract is that nothing
// was stored — the reload that proves it is the restore. The email case does
// store, and `afterAll` puts the address back.
//
// A button naming a script that is NOT loaded renders danger-outlined with its
// tooltip; that case is not asserted here. The fixture record page carries
// exactly the two buttons below, and the Automation configurator only offers
// scripts the registry already has, so the case cannot be built through the UI
// without editing the fixture page. AutomationButtons.test.js covers it.

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:5173'
const STORAGE_STATE = 'e2e/.auth/state.json'

const NAMESPACE = '/compose/namespace/agent-sandbox'
const GREET = 'Greet contact (client)'
const ACTIVATE = 'Activate contact (server)'
const PROBE_EMAIL = '  Probe@Example.COM  '
const STORED_EMAIL = 'probe@example.com'

/** One field row of the record form, addressed by the field's own label. Label
 *  and control are not associated, so the row — `field-item`, the field
 *  renderer's own class — is the scope, and `field-value` holds what the field
 *  shows or offers. */
function fieldRow(page: Page, field: string): Locator {
  return page.locator('.field-item').filter({ has: page.getByText(field, { exact: true }) })
}

/** The editable control of a record form field. */
function recordInput(page: Page, field: string): Locator {
  return fieldRow(page, field).getByRole('textbox')
}

/** What a field shows on the record viewer. */
function recordField(page: Page, field: string): Locator {
  return fieldRow(page, field).locator('.field-value')
}

function watchForCrashes(page: Page): string[] {
  const errors: string[] = []
  page.on('pageerror', e => errors.push(String(e)))
  return errors
}

test.describe.serial('corredor scripts on a compose record page', () => {
  let unavailable = ''
  let recordPath = ''
  let contactName = ''
  let originalEmail = ''

  test.beforeAll(async ({ browser }) => {
    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()

      // The namespace's landing page is the fixture dashboard, whose record list
      // is the only thing this spec needs to find a contact — no id is assumed.
      await page.goto(NAMESPACE)
      const rows = page.locator('tbody tr')
      const listed = await rows
        .first()
        .waitFor({ state: 'visible', timeout: 30000 })
        .then(() => true)
        .catch(() => false)

      if (!listed) {
        unavailable =
          'the agent-sandbox namespace and its contacts are not on this stack — seed it with dev/agent/seed.sh agent-sandbox'
        return
      }

      // The server script activates a lead, so the spec needs one.
      const lead = rows.filter({ hasText: 'Lead' }).first()
      if ((await lead.count()) === 0) {
        unavailable =
          'no agent-contact record has status lead on this stack — reseed agent-sandbox so the activate script has something to change'
        return
      }

      await lead.click()
      await page.waitForURL(/\/pages\/\d+\/records\/\d+/, { timeout: 30000 })
      recordPath = new URL(page.url()).pathname

      // The buttons appear once the client bundle is registered, after boot.
      const wired = await page
        .getByRole('button', { name: GREET })
        .waitFor({ state: 'visible', timeout: 20000 })
        .then(() => true)
        .catch(() => false)

      if (!wired) {
        unavailable =
          "the contact page's Automation block has no fixture script buttons — Corredor is off, or dev/fixtures/agent-sandbox predates the block"
        return
      }

      // The editor is where the field values can be read back verbatim.
      await page.goto(`${recordPath}?edit=1`)
      await expect(recordInput(page, 'email')).toBeVisible({ timeout: 30000 })
      contactName = await recordInput(page, 'name').inputValue()
      originalEmail = await recordInput(page, 'email').inputValue()

      if (!contactName || !originalEmail) {
        unavailable = 'the contact this spec picked has no name or email to assert on'
      }
    } finally {
      await context.close()
    }
  })

  test.afterAll(async ({ browser }) => {
    if (!recordPath || !originalEmail) return

    const context = await browser.newContext({ baseURL: BASE_URL, storageState: STORAGE_STATE })
    try {
      const page = await context.newPage()
      await page.goto(`${recordPath}?edit=1`)

      const email = recordInput(page, 'email')
      await expect(email).toBeVisible({ timeout: 30000 })
      if ((await email.inputValue()) === originalEmail) return

      await email.fill(originalEmail)
      await page.getByRole('button', { name: 'Save', exact: true }).click()
      await page.waitForURL(url => !url.searchParams.has('edit'), { timeout: 30000 })
      await expect(recordField(page, 'email')).toContainText(originalEmail)
    } finally {
      await context.close()
    }
  })

  test('the client script greets the contact and leaves the record viewer open', async ({
    page,
  }) => {
    test.skip(!!unavailable, unavailable)
    const errors = watchForCrashes(page)

    // From the editor, so the script's `gotoRecordViewer` has somewhere to go.
    await page.goto(`${recordPath}?edit=1`)
    const greet = page.getByRole('button', { name: GREET })
    await expect(greet).toBeVisible({ timeout: 30000 })
    await greet.click()

    const toast = page.getByRole('alert')
    await expect(toast).toContainText(contactName)
    await expect(toast).toContainText('(client script)')

    // Same record, now on its viewer.
    await page.waitForURL(url => !url.searchParams.has('edit'), { timeout: 30000 })
    expect(new URL(page.url()).pathname).toBe(recordPath)
    expect(errors).toEqual([])
  })

  test('the server script activates the lead on the page without storing it', async ({ page }) => {
    test.skip(!!unavailable, unavailable)
    const errors = watchForCrashes(page)

    await page.goto(recordPath)
    const status = recordField(page, 'status')
    await expect(status).toContainText('Lead', { timeout: 30000 })

    await page.getByRole('button', { name: ACTIVATE }).click()

    // The endpoint answers the record the script returned; the page adopts it.
    await expect(status).toContainText('Active')
    await expect(status).not.toContainText('Lead')
    expect(errors).toEqual([])

    // And only the page: a manual trigger never writes the record.
    await page.reload()
    await expect(recordField(page, 'status')).toContainText('Lead', { timeout: 30000 })
  })

  test('the client script normalises the email the form submits', async ({ page }) => {
    test.skip(!!unavailable, unavailable)
    const errors = watchForCrashes(page)

    await page.goto(`${recordPath}?edit=1`)
    const email = recordInput(page, 'email')
    await expect(email).toBeVisible({ timeout: 30000 })

    await email.fill(PROBE_EMAIL)
    await page.getByRole('button', { name: 'Save', exact: true }).click()
    await page.waitForURL(url => !url.searchParams.has('edit'), { timeout: 30000 })

    await expect(recordField(page, 'email')).toContainText(STORED_EMAIL)
    expect(errors).toEqual([])

    // Stored, not merely displayed: the script corrected the record that was sent.
    await page.reload()
    await expect(recordField(page, 'email')).toContainText(STORED_EMAIL, { timeout: 30000 })
  })
})
