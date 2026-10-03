import { expect, test, type Page } from '@playwright/test'

// Fixture: page "E2E Buttons" with an Automation block holding four buttons:
// "Always" (no condition), "Hidden" (visibility expression false), "Shown"
// (visibility expression true) and "Prompt" (runs a workflow whose only step
// is a blocking alert prompt).

async function openButtonsPage (page: Page) {
  await page.goto('/ns/e2e_ns/pages')
  await page.getByRole('link', { name: 'E2E Buttons', exact: true }).first().click()
  await expect(page.getByRole('heading', { name: 'E2E Buttons' })).toBeVisible()
  await expect(page.getByRole('button', { name: 'Always' })).toBeVisible()
}

test('buttons with a visibility expression are shown or hidden accordingly', async ({ page }) => {
  await openButtonsPage(page)

  await expect(page.getByRole('button', { name: 'Shown' })).toBeVisible()
  await expect(page.getByRole('button', { name: 'Hidden' })).toHaveCount(0)
})

test('a blocking prompt is a modal that only closes by answering', async ({ page }) => {
  await openButtonsPage(page)
  await page.getByRole('button', { name: 'Prompt' }).click()

  const modal = page.locator('.modal-dialog', { hasText: 'Blocking message' })
  await expect(modal).toBeVisible()
  await expect(modal.locator('.modal-header')).toContainText('Prompt title')
  await expect(modal.locator('.modal-header button.close')).toHaveCount(0)

  // backdrop and escape do not close it
  await page.keyboard.press('Escape')
  await expect(modal).toBeVisible()

  // answering does
  await modal.getByRole('button').first().click()
  await expect(modal).toBeHidden()
})
