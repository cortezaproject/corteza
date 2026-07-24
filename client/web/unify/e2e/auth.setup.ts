import { expect, test as setup } from '@playwright/test'

// Logs in once through the real auth server form (fields: email/password —
// server/auth/assets/templates/login.html.tpl) and saves the session for every
// spec via storageState. Credentials come from .env.e2e, never the repo.
setup('authenticate', async ({ page }) => {
  const user = process.env.E2E_USER
  const pass = process.env.E2E_PASS
  if (!user || !pass)
    throw new Error('E2E_USER / E2E_PASS not set — copy .env.e2e.example to .env.e2e')

  await page.goto('/')
  // Unauthenticated visit redirects to the auth server's login page
  await page.waitForURL(/\/auth\//)
  await page.fill('input[name="email"]', user)
  await page.fill('input[name="password"]', pass)
  await page.click('button[type="submit"]')

  // Back in the app: the topbar renders for authenticated users
  await page.waitForURL(url => !url.pathname.includes('/auth/'))
  await expect(page.locator('header')).toBeVisible()

  await page.context().storageState({ path: 'e2e/.auth/state.json' })
})
