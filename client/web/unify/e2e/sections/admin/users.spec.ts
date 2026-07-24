import { expect, test } from '@playwright/test'

// Smoke: system users admin list renders rows and opens the editor
// (contract: sections/admin/views/system/User/List.intent.md).
test('user list renders and opens a user editor', async ({ page }) => {
  await page.goto('/admin/system/users')
  const rows = page.locator('tbody tr')
  await expect(rows.first()).toBeVisible()

  await rows.first().click()
  await page.waitForURL(/\/admin\/system\/users\//)
  // The editor's required identity field (Editor.intent.md: email is required)
  await expect(page.locator('input[name="email"]')).toBeVisible()
})
