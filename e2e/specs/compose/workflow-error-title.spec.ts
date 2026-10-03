import { expect, test, type Page } from '@playwright/test'

// Fixture: page "E2E Guarded" with an inline-editable record list over
// `e2e_guarded`; a beforeUpdate workflow on that module always stops with an
// error step whose message is "Order is not ready" and title "Order check".
//
// The message must reach the user exactly as configured (no "workflow N step
// M execution failed" prefix) and the title must become the notification's
// heading.

async function openGuardedPage (page: Page) {
  await page.goto('/ns/e2e_ns/pages')
  await page.getByRole('link', { name: 'E2E Guarded', exact: true }).first().click()
  await expect(page.getByRole('heading', { name: 'E2E Guarded' })).toBeVisible()
  await expect(page.locator('[data-test-id="table-record-list"]')).toBeVisible()
}

function rows (page: Page) {
  return page.locator('[data-test-id="table-record-list"] tbody tr')
}

function saveButton (page: Page, row: ReturnType<typeof rows>) {
  return row.locator('button', { has: page.locator('svg[data-icon="check"]') }).first()
}

test('error step message and title are shown as configured', async ({ page }) => {
  await openGuardedPage(page)

  // creating a record is not guarded; make one through the app's API client
  // so there is a row to update (the inline "add" row needs an existing list)
  await page.evaluate(async () => {
    let el: any = document.querySelector('[data-test-id="table-record-list"]')
    let vm: any = null
    while (el && !vm) { vm = el.__vue__; el = el.parentElement }
    while (vm && typeof vm.dirtyRecordsCount === 'undefined') vm = vm.$parent
    const { namespaceID, moduleID } = vm.recordListModule
    await vm.$ComposeAPI.recordCreate({ namespaceID, moduleID, values: [{ name: 'subject', value: `Guarded ${Date.now()}` }] })
  })
  await page.reload()
  await expect(page.locator('[data-test-id="table-record-list"] tbody tr').first()).toBeVisible()

  // updating it runs the workflow, which refuses the change
  const row = rows(page).first()
  const input = row.locator('input[type="text"]').first()
  await input.fill('Guarded edited')
  await input.press('Tab')
  await saveButton(page, row).click()

  const toast = page.locator('.toast', { hasText: 'Order is not ready' })
  await expect(toast).toBeVisible()
  await expect(toast.locator('.toast-header')).toContainText('Order check')
  await expect(toast).not.toContainText('execution failed')
})
