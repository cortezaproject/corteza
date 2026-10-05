import { expect, test } from '@playwright/test'

// Fixture: namespace `e2e_ns` with page "E2E Dashboard" holding a Progress
// block "Empty progress" that counts records of a module without any records.
//
// Before the fix the block rendered "NaN%": an empty report divided 0 by 0.
test('progress block over an empty module shows the default instead of NaN', async ({ page }) => {
  await page.goto('/ns/e2e_ns/pages')
  await page.getByRole('link', { name: 'E2E Dashboard' }).first().click()

  // the dashboard holds a single progress block
  await expect(page.getByText('Empty progress').first()).toBeVisible()
  const label = page.locator('.progress strong').first()

  // the label sits inside the 0% wide bar, so it has no size of its own
  await expect(label).toBeAttached()
  await expect(label).not.toContainText('NaN')
  await expect(label).toHaveText(/^\s*0(\.0+)?%\s*$/)
})
