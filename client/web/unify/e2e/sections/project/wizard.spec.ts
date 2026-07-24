import { expect, test } from '@playwright/test'

// Smoke: the wizard's locked three-tab shape (Build / Govern / Manage &
// Monitor, per sections/project/project.intent.md) renders with icons and
// switches tabs. Uses the first project in the list; skips when none exist.
test('wizard shows the three tabs and switches between them', async ({ page }) => {
  await page.goto('/project/projects')
  const firstProject = page.locator('[data-pc-section="content"] a, tbody tr a').first()
  if ((await firstProject.count()) === 0) test.skip(true, 'no projects in dev data')

  await firstProject.click()
  await page.waitForURL(/\/wizard/)

  const tabs = page.getByRole('tab')
  await expect(tabs).toHaveCount(3)
  for (const icon of ['pi-wrench', 'pi-shield', 'pi-chart-line']) {
    await expect(page.locator(`[role="tab"] .${icon}`)).toBeVisible()
  }

  await tabs.nth(1).click()
  await expect(tabs.nth(1)).toHaveAttribute('aria-selected', 'true')
  await tabs.nth(0).click()
  await expect(tabs.nth(0)).toHaveAttribute('aria-selected', 'true')
})
