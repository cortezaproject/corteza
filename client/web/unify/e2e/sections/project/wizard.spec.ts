import { expect, test } from '@playwright/test'

// Smoke: the wizard's locked four-tab shape (Build / Govern / Manage &
// Monitor / Publish, per sections/project/project.intent.md) renders with
// icons and switches tabs. Uses the first non-live project in the list (live
// rows open the overview, not the wizard); skips when none exist.
test('wizard shows the four tabs and switches between them', async ({ page }) => {
  await page.goto('/project/projects')
  await page.waitForLoadState('networkidle')

  // Rows navigate on click (no anchors in the list); filter out live projects
  // by their status tag — those route to project.overview instead.
  const editable = page
    .locator('tbody tr')
    .filter({ hasNot: page.getByText('Active', { exact: true }) })
    .filter({ hasNot: page.getByText('Published', { exact: true }) })
  if ((await editable.count()) === 0) test.skip(true, 'no non-live projects in dev data')

  await editable.first().click()
  await page.waitForURL(/\/wizard/)

  const tabs = page.getByRole('tab')
  await expect(tabs).toHaveCount(4)
  for (const icon of ['pi-wrench', 'pi-shield', 'pi-chart-line', 'pi-cloud-upload']) {
    await expect(page.locator(`[role="tab"] .${icon}`)).toBeVisible()
  }

  await tabs.nth(1).click()
  await expect(tabs.nth(1)).toHaveAttribute('aria-selected', 'true')
  await tabs.nth(0).click()
  await expect(tabs.nth(0)).toHaveAttribute('aria-selected', 'true')
})
