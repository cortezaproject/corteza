import { expect, test } from '@playwright/test'

// Smoke: the home landing loads at / with the app topbar
// (contract: sections/home/home.intent.md).
test('home landing renders', async ({ page }) => {
  await page.goto('/')
  await expect(page.locator('header')).toBeVisible()
  await expect(page).toHaveURL(/\/$/)
})
