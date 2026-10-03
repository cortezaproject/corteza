import { expect, test, type Page } from '@playwright/test'

// Page builder: pasting a block is a window "paste" event with the block's
// JSON (copy adds the source pageID). Before the fixes, pasting the same
// block twice froze the builder (two grid entries with one ID) and a Tabs
// block pasted from another page grabbed the blocks with the same numbers on
// the target page.

async function openDashboardBuilder (page: Page) {
  await page.goto('/ns/e2e_ns/pages')
  await page.getByRole('link', { name: 'E2E Dashboard', exact: true }).first().click()
  await expect(page.getByRole('heading', { name: 'E2E Dashboard' })).toBeVisible()
  await page.locator('a[href*="/admin/pages/"][href*="/builder"]').first().click()
  await expect(page.locator('#page-builder')).toBeVisible()
}

function builderState (page: Page) {
  return page.evaluate(() => {
    let el: any = document.querySelector('#page-builder')
    let vm: any = el.__vue__
    while (vm && typeof vm.pasteBlock === 'undefined') vm = vm.$parent
    return {
      pageID: vm.page.pageID as string,
      blocks: vm.blocks.map((b: any) => ({ kind: b.kind, title: b.title, tabs: b.kind === 'Tabs' ? (b.options.tabs || []).length : undefined })),
    }
  })
}

async function paste (page: Page, block: object) {
  await page.locator('#page-builder').focus()
  await page.evaluate((json) => {
    const dt = new DataTransfer()
    dt.setData('text', json)
    window.dispatchEvent(new ClipboardEvent('paste', { clipboardData: dt, bubbles: true, cancelable: true }))
  }, JSON.stringify(block))
}

const tabsBlock = (pageID: string) => ({
  pageID,
  kind: 'Tabs',
  title: 'Pasted tabs',
  xywh: [0, 20, 6, 4],
  options: { tabs: [{ blockID: '1', title: 'First' }, { blockID: '2', title: 'Second' }] },
})

test('pasting the same block twice adds two blocks', async ({ page }) => {
  await openDashboardBuilder(page)
  const before = await builderState(page)

  await paste(page, tabsBlock(before.pageID))
  await paste(page, tabsBlock(before.pageID))

  await expect.poll(async () => (await builderState(page)).blocks.length).toBe(before.blocks.length + 2)
  const after = await builderState(page)
  expect(after.blocks.filter(b => b.title === 'Pasted tabs')).toHaveLength(2)
  // tabs from the same page keep their block references
  expect(after.blocks.filter(b => b.title === 'Pasted tabs').map(b => b.tabs)).toEqual([2, 2])
})

test('tabs pasted from another page do not take over same-numbered blocks', async ({ page }) => {
  await openDashboardBuilder(page)
  const before = await builderState(page)

  await paste(page, tabsBlock('999999999'))

  await expect.poll(async () => (await builderState(page)).blocks.length).toBe(before.blocks.length + 1)
  const after = await builderState(page)
  expect(after.blocks.find(b => b.title === 'Pasted tabs')!.tabs).toBe(0)
})
