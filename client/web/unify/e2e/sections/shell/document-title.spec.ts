import { expect, test } from '@playwright/test'

// The browser tab names where you are: the shell mirrors whatever heading the
// active view teleports into #topbar-title into document.title
// (contract: App.intent.md; the helper is src/utils/documentTitle.js).
//
// The assertion reads the RENDERED topbar text rather than the extraction
// helper — comparing the helper against itself would pass with the mirroring
// switched off.

const BASE_TITLE = 'Human'

type Page = import('@playwright/test').Page

// Every test here navigates two or three times, and a fresh context pays for an
// oauth roundtrip on the first one; the default per-test budget runs out mid-run.
test.describe.configure({ timeout: 90_000 })

// The shell landmark, not `header`: vite compiles a route the first time it is
// asked for, so the default expect timeout calls a cold route a broken one.
const openShell = async (page: Page, path: string) => {
  await page.goto(path)
  await expect(page.getByTestId('app-topbar')).toBeVisible({ timeout: 30_000 })
}

const heading = (page: Page) =>
  page.evaluate(() => document.getElementById('topbar-title')?.innerText.trim() ?? '')

// A heading arrives with the data behind it (a project name, a page title), so
// it is polled for rather than read once.
const headingArrives = (page: Page) =>
  expect.poll(() => heading(page), { timeout: 20_000 }).not.toBe('')

// Routes whose heading is plain text with nothing beside it.
const MIRRORED = [
  '/admin/system/users',
  '/admin/system/roles',
  '/taq',
  '/compose/namespaces',
  '/workflow/list',
  '/agentic',
  '/chatbot',
  '/project/projects',
]

for (const path of MIRRORED) {
  test(`tab title is the topbar heading at ${path}`, async ({ page }) => {
    await openShell(page, path)
    await headingArrives(page)

    expect(await page.title()).toBe(await heading(page))
  })
}

test('a route that states no heading keeps the app name', async ({ page }) => {
  await openShell(page, '/')

  expect(await heading(page)).toBe('')
  await expect(page).toHaveTitle(BASE_TITLE)
})

test('the title follows navigation, leaving nothing stale behind', async ({ page }) => {
  await openShell(page, '/admin/system/users')
  await headingArrives(page)
  const first = await page.title()

  await openShell(page, '/taq')
  await headingArrives(page)
  const second = await page.title()

  expect(second).not.toBe(first)
  expect(second).toBe(await heading(page))
})

test('the project dashboard titles the project, not the revision switcher', async ({ page }) => {
  // Tolerates a stack with no projects (e2e never seeds), but only after the
  // list has had a chance to load — an empty table is also what a still-
  // fetching one looks like.
  await openShell(page, '/project/projects')
  const row = page.locator('table tbody tr').first()
  await row.waitFor({ state: 'visible' }).catch(() => {})
  test.skip(!(await row.isVisible()), 'no project on this stack')

  await row.click()
  await page.waitForURL(/\/project\/projects\/\d+/)
  await headingArrives(page)

  const [topbar, title] = [await heading(page), await page.title()]

  // The switcher sits inside #topbar-title and carries data-title-exclude, so
  // the topbar says more than the tab does.
  expect(title).not.toBe(BASE_TITLE)
  expect(topbar.startsWith(title)).toBe(true)
  expect(topbar.length).toBeGreaterThan(title.length)
})
