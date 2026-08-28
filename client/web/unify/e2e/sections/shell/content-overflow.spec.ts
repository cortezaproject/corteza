import { expect, test } from '@playwright/test'

// The shell fits the viewport, whatever a section puts inside it.
//
// The content column sits beside the sidebar on a margin, and it is a flex
// item: with no floor on its width its automatic minimum size is its content's
// min-content width, so one wide block (a record list with many columns) makes
// the column wider than the space left beside the sidebar, and the margin
// carries its right edge — the topbar with it — off the screen.
//
// Asserted on scrollWidth rather than on a role or a label because the contract
// here is a measurement: nothing names "the page is 300px too wide".
//
// The wide element is injected rather than built from a record list. What the
// column reacts to is a min-content width larger than the space beside the
// sidebar; a 15-column table is one way to produce one, an explicit width is
// the same measurement without a fixture the suite is not allowed to seed.

type Page = import('@playwright/test').Page

const VIEWPORT = { width: 1280, height: 900 }

// Admin draws its sidebar on every route and needs no namespace to exist.
const ROUTE = '/admin/system/users'
const SECTION = 'admin'

test.describe.configure({ timeout: 90_000 })

const openShell = async (page: Page, path: string) => {
  await page.goto(path)
  await expect(page.getByTestId('app-topbar')).toBeVisible({ timeout: 30_000 })
}

const docWidth = (page: Page) => page.evaluate(() => document.documentElement.scrollWidth)

// The drawer slides in and the content column's margin follows it, so both
// edges are moving numbers for a few hundred ms after the route renders; a rect
// read straight away measures the animation rather than the layout.
const settled = async <T>(page: Page, read: () => Promise<T>): Promise<T> => {
  let last = await read()
  for (;;) {
    await page.waitForTimeout(120)
    const now = await read()
    if (JSON.stringify(now) === JSON.stringify(last)) return now
    last = now
  }
}

const edges = (page: Page) =>
  settled(page, () =>
    page.evaluate(() => {
      const sidebar = document.querySelector('[data-testid="app-sidebar"]')
      const column = document.querySelector('main')!.getBoundingClientRect()
      return {
        sidebarWidth: Math.round(sidebar?.getBoundingClientRect().width ?? 0),
        columnLeft: Math.round(column.left),
        columnRight: Math.round(column.right),
      }
    }),
  )

const widen = (page: Page) =>
  page.evaluate(() => {
    const probe = document.createElement('div')
    probe.dataset.testid = 'overflow-probe'
    probe.style.width = '5000px'
    probe.style.height = '1px'
    document.querySelector('main')?.appendChild(probe)
  })

test('wide content cannot push the shell past the viewport', async ({ page }) => {
  await page.setViewportSize(VIEWPORT)

  // Expansion is a per-section localStorage preference, so the origin has to
  // exist before it can be written; the second visit is the one under test.
  await openShell(page, ROUTE)
  await page.evaluate(
    section => localStorage.setItem('ui.sidebar.expanded', JSON.stringify({ [section]: true })),
    SECTION,
  )
  await openShell(page, ROUTE)
  await expect(page.getByTestId('app-sidebar')).toBeVisible({ timeout: 30_000 })
  await edges(page)

  expect(await docWidth(page)).toBeLessThanOrEqual(VIEWPORT.width)

  await widen(page)

  await expect.poll(() => docWidth(page), { timeout: 10_000 }).toBeLessThanOrEqual(VIEWPORT.width)

  // The symptom the measurement stands for: the topbar leaving the screen.
  const topbarRight = await page
    .getByTestId('app-topbar')
    .evaluate(el => Math.round(el.getBoundingClientRect().right))
  expect(topbarRight).toBeLessThanOrEqual(VIEWPORT.width)
})

test('the content column still fills the space beside the sidebar', async ({ page }) => {
  // The floor that stops the overflow must not become a column that under-fills:
  // sidebar and content together account for the whole viewport.
  await page.setViewportSize(VIEWPORT)

  await openShell(page, ROUTE)
  await page.evaluate(
    section => localStorage.setItem('ui.sidebar.expanded', JSON.stringify({ [section]: true })),
    SECTION,
  )
  await openShell(page, ROUTE)
  await expect(page.getByTestId('app-sidebar')).toBeVisible({ timeout: 30_000 })

  await widen(page)

  // <main> is the column's own box — the topbar carries padding, so its rect
  // answers a different question by a few pixels.
  const { sidebarWidth, columnLeft, columnRight } = await edges(page)

  expect(sidebarWidth).toBeGreaterThan(100)
  expect(columnLeft).toBe(sidebarWidth)
  expect(columnRight).toBe(VIEWPORT.width)
})
