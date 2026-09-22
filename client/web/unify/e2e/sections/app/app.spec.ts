import { expect, test } from '@playwright/test'

// A custom application renders inside the two-frame sandbox at /app/:id
// (contract: sections/app/app.intent.md). The stack is not seeded here: the
// spec uses whichever custom application the app menu already offers and
// skips when there is none.

async function firstCustomAppHref(page) {
  await page.goto('/')
  const tile = page.locator('a[href^="/app/"]').first()
  // The app list arrives after the shell; give it the time a slow stack needs.
  const offered = await tile
    .waitFor({ timeout: 10000 })
    .then(() => true)
    .catch(() => false)
  return offered ? tile.getAttribute('href') : null
}

test('a custom application renders inside the sandbox', async ({ page }) => {
  const href = await firstCustomAppHref(page)
  test.skip(!href, 'no custom application is offered to this user')

  await page.goto(href!)

  const outer = page.frameLocator('iframe[srcdoc]').first()
  const inner = outer.frameLocator('iframe[sandbox="allow-scripts"]').first()
  await expect(inner.locator('body')).toBeVisible()

  // The outer document is what pins the inner one in place.
  const outerSrcdoc = await page.locator('iframe[srcdoc]').first().getAttribute('srcdoc')
  expect(outerSrcdoc).toContain(`content="frame-src 'none'"`)
  expect(outerSrcdoc).toContain(`connect-src 'none'`)
})

test('the app cannot reach the network', async ({ page }) => {
  const href = await firstCustomAppHref(page)
  test.skip(!href, 'no custom application is offered to this user')

  const refusals: string[] = []
  page.on('console', m => {
    if (/Content Security Policy/.test(m.text())) refusals.push(m.text())
  })

  await page.goto(href!)
  const inner = page
    .frameLocator('iframe[srcdoc]')
    .first()
    .frameLocator('iframe[sandbox="allow-scripts"]')
    .first()
  await expect(inner.locator('body')).toBeVisible()

  await inner.locator('body').evaluate(() => fetch('http://127.0.0.1:9/leak').catch(() => {}))
  await expect.poll(() => refusals.length).toBeGreaterThan(0)
  expect(refusals.join('\n')).toContain(`connect-src 'none'`)
})

test('an application that is not a custom one is refused', async ({ page }) => {
  await page.goto('/app/999')
  await expect(page).toHaveURL(/\/\?denied=app$/)
})
