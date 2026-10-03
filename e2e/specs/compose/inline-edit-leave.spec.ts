import { expect, test, type Page } from '@playwright/test'

// Fixture: page "E2E Inline" in namespace `e2e_ns` with a record list over
// `e2e_items` in inline-edit mode (editable cells).
//
// Before the fix nothing warned about unsaved inline edits: navigating away
// or reloading silently dropped them, while a record page asked first.

async function openInlinePage (page: Page) {
  await page.goto('/ns/e2e_ns/pages')
  // the pages index lands on the first page; switch to the inline one
  await page.getByRole('link', { name: 'E2E Inline', exact: true }).first().click()
  await expect(page.getByRole('heading', { name: 'E2E Inline' })).toBeVisible()
  await expect(page.locator('[data-test-id="table-record-list"] tbody tr').first()).toBeVisible()
}

function firstTitleInput (page: Page) {
  return page.locator('[data-test-id="table-record-list"] tbody tr').first().locator('input[type="text"]').first()
}

test('leaving a record list with unsaved inline edits asks first', async ({ page }) => {
  await openInlinePage(page)
  const inlineURL = page.url()

  const input = firstTitleInput(page)
  await expect(input).toBeVisible()
  await input.fill('edited inline')
  // commit the change the way a user does, by leaving the cell
  await input.press('Tab')

  // in-app navigation: decline and stay
  const seen: string[] = []
  page.once('dialog', async d => { seen.push(d.message()); await d.dismiss() })
  await page.goBack()
  await expect.poll(() => seen.length).toBe(1)
  expect(seen[0]).toContain('Unsaved changes')
  await expect(page).toHaveURL(inlineURL)
  await expect(firstTitleInput(page)).toHaveValue('edited inline')

  // in-app navigation: accept and leave
  page.once('dialog', async d => { await d.accept() })
  await page.goBack()
  await expect(page).not.toHaveURL(inlineURL)
})

test('reloading a record list with unsaved inline edits warns', async ({ page }) => {
  await openInlinePage(page)

  const input = firstTitleInput(page)
  await input.fill('edited inline')
  await input.press('Tab')

  const types: string[] = []
  page.on('dialog', async d => { types.push(d.type()); await d.accept() })
  await page.reload()
  await expect.poll(() => types).toContain('beforeunload')
})

test('a record list without edits does not ask', async ({ page }) => {
  await openInlinePage(page)
  const inlineURL = page.url()

  let dialogs = 0
  page.on('dialog', async d => { dialogs++; await d.dismiss() })
  await page.goBack()
  await expect(page).not.toHaveURL(inlineURL)
  expect(dialogs).toBe(0)
})
