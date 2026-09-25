import { expect, test } from '@playwright/test'

// The user's own home application: picked from a tile's menu, badged on its
// tile, opened on entry at `/` and by the topbar's home button without a
// reload (contract: sections/home/home.intent.md).
test('own home application opens on entry', async ({ page }) => {
  // Four full page loads, each through the auth flow.
  test.setTimeout(120000)

  const appsColumn = page.locator('.column-panel').first()
  const tile = appsColumn.locator('[data-drag-item]', { hasText: 'Automation (TAQ)' })

  const pick = async (label: RegExp) => {
    await tile.hover()
    await tile.getByRole('button', { name: 'Application options' }).click()
    const saved = page.waitForResponse(
      r => r.request().method() === 'PUT' && /\/system\/users\/\d+$/.test(r.url()),
    )
    await page.getByRole('menuitem', { name: label }).click()
    expect((await saved).ok()).toBe(true)
    await expect(page.locator('.p-toast-message')).toBeVisible()
  }

  // A query keeps the entry redirect out of the way, whatever an earlier run left.
  await page.goto('/?stay=1')
  await expect(tile).toBeVisible({ timeout: 20000 })
  await pick(/Set as Home/)
  await expect(tile.locator('[data-test-id="app-tile-home"]')).toBeVisible()

  try {
    await page.goto('/')
    await expect(page).toHaveURL(/\/taq/, { timeout: 20000 })

    // A marker a page load would wipe proves the home button routes in place.
    await page.goto('/agentic/')
    await page.evaluate(() => ((window as any).__noReload = true))
    await page.locator('header a:has(.pi-home)').click()
    await expect(page).toHaveURL(/\/taq/)
    expect(await page.evaluate(() => (window as any).__noReload)).toBe(true)
  } finally {
    await page.goto('/?stay=1')
    await expect(tile).toBeVisible({ timeout: 20000 })
    await pick(/Remove as Home/)
  }

  await page.goto('/')
  await expect(tile).toBeVisible({ timeout: 20000 })
  await expect(page).toHaveURL(/\/$/)
})
