import { expect, test, type Page } from '@playwright/test'

// Fixture: pages "E2E Comments Newest" (sortDirection asc) and "E2E Comments
// Oldest" (sortDirection desc) over module `e2e_comments`.
//
// Before the fix the block re-sorted comments on the client, so "oldest
// first" showed the newest comment at the top.

async function createComments (page: Page) {
  await page.goto('/ns/e2e_ns/pages')
  await page.getByRole('link', { name: 'E2E Comments Newest', exact: true }).first().click()
  await expect(page.getByRole('heading', { name: 'E2E Comments Newest' })).toBeVisible()

  await page.evaluate(async () => {
    // start from the page heading awaited above; comments may not be rendered yet
    let el: any = document.querySelector('h2.title') || document.body
    let vm: any = el.__vue__
    while (!vm && el.parentElement) { el = el.parentElement; vm = el.__vue__ }
    const api = vm.$ComposeAPI
    const { namespaceID } = vm.$store.getters['namespace/getByHandle'] ? vm.$store.getters['namespace/getByHandle']('e2e_ns') || {} : {}
    const ns = namespaceID || (vm.$store.getters['namespace/set'] || []).find((n: any) => n.slug === 'e2e_ns').namespaceID
    const mod = (vm.$store.getters['module/set'] || []).find((m: any) => m.handle === 'e2e_comments')
    // remove comments of earlier runs, they would push these off the first
    // page of the "oldest first" list
    await api.recordBulkDelete({ namespaceID: ns, moduleID: mod.moduleID, query: "title LIKE 'c%'" })

    const stamp = Date.now()
    for (const n of [1, 2, 3]) {
      await api.recordCreate({ namespaceID: ns, moduleID: mod.moduleID, values: [{ name: 'title', value: `c${n}-${stamp}` }, { name: 'content', value: `comment ${n} ${stamp}` }] })
      await new Promise(r => setTimeout(r, 1100))
    }
    ;(window as any).__stamp = stamp
  })

  return page.evaluate(() => (window as any).__stamp as number)
}

async function visibleOrder (page: Page, title: string, stamp: number) {
  await page.goto('/ns/e2e_ns/pages')
  await page.getByRole('link', { name: title, exact: true }).first().click()
  await expect(page.getByRole('heading', { name: title })).toBeVisible()
  await expect(page.locator('.comment-item').first()).toBeVisible()
  const texts = await page.locator('.comment-item').allInnerTexts()
  return texts.map(t => (t.match(new RegExp(`comment (\\d) ${stamp}`)) || [])[1]).filter(Boolean)
}

// Both modes render comments chronologically (a conversation); "oldest first"
// used to re-sort on the client and put the newest comment at the top.
test('comments are rendered chronologically in both modes', async ({ page }) => {
  const stamp = await createComments(page)

  const oldest = await visibleOrder(page, 'E2E Comments Oldest', stamp)
  expect(oldest).toEqual(['1', '2', '3'])

  // newest first loads the latest page first and scrolls to the end
  const newest = await visibleOrder(page, 'E2E Comments Newest', stamp)
  expect(newest).toEqual(['1', '2', '3'])
})
