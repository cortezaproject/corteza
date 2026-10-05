import { expect, test, type Page } from '@playwright/test'

// Fixture: page "E2E Dashboard" in namespace `e2e_ns` with a selectable
// record list over `e2e_items` (3 records) with bulk editing enabled.
//
// Before the fix the "Update selected records" modal kept the fields and
// values chosen in a previous round after being closed, so the next bulk edit
// started with stale data.

const ADD_FIELD = 'Add a field to update'
const ADD_ANOTHER = 'Add another field to update'

async function openDashboard (page: Page) {
  await page.goto('/ns/e2e_ns/pages')
  await page.getByRole('link', { name: 'E2E Dashboard' }).first().click()
  await expect(page.locator('[data-test-id="table-record-list"] tbody tr').first()).toBeVisible()
}

async function selectFirstRow (page: Page) {
  // the custom checkbox draws its box over the hidden input, which would
  // otherwise make Playwright wait for the input to receive the click
  const checkbox = page.locator('[data-test-id="table-record-list"] tbody tr').first().locator('input[type="checkbox"]').first()
  await checkbox.check({ force: true })
  await expect(checkbox).toBeChecked()
}

async function openBulkEdit (page: Page) {
  // the pen button appears in the record list toolbar once rows are selected
  await page.locator('button.inline-button', { has: page.locator('svg[data-icon="pen"]') }).first().click()
  const modal = page.locator('.modal-dialog', { hasText: 'Update selected records' })
  await expect(modal).toBeVisible()
  return modal
}

test('bulk edit modal forgets the fields and values after it is closed', async ({ page }) => {
  await openDashboard(page)
  await selectFirstRow(page)

  // the field picker's placeholder tells whether any field is chosen yet
  const picker = () => modal.locator('.vs__search').first()

  // first round: pick a field and type a value, then cancel
  let modal = await openBulkEdit(page)
  await expect(picker()).toHaveAttribute('placeholder', ADD_FIELD)
  await picker().click()
  await page.keyboard.type('Title')
  await page.keyboard.press('Enter')
  await expect(picker()).toHaveAttribute('placeholder', ADD_ANOTHER)
  const editor = modal.locator('.card input').first()
  await expect(editor).toBeVisible()
  await editor.fill('stale value')
  await modal.getByRole('button', { name: 'Cancel' }).click()
  await expect(modal).toBeHidden()

  // second round: the modal must start clean
  modal = await openBulkEdit(page)
  await expect(picker()).toHaveAttribute('placeholder', ADD_FIELD)
  await expect(modal.locator('.card input')).toHaveCount(0)
})
