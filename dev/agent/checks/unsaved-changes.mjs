// Unsaved-changes guards: does leaving an editor with staged edits warn first?
//
// Run: node dev/agent/drive.mjs dev/agent/checks/unsaved-changes.mjs
//      node dev/agent/drive.mjs dev/agent/checks/unsaved-changes.mjs --only "builder"
//
// The page builder and the record create form each shipped without a guard, so
// a block added to a layout or a filled-in record was lost to one stray click.
// Both halves matter: an editor that never warns loses work, and one that warns
// on load cries wolf until nobody reads it.
import { drive, check, ids } from '../drive.mjs'

const NS = 'catalogue'
const id = ids(NS)

const BUILDER = `/compose/namespace/${NS}/admin/pages/${id.page.catalogue_content}/builder`
const RECORD_CREATE = `/compose/namespace/${NS}/admin/modules/${id.module.catalogue_field}/records/create`

async function openBuilder(page) {
  await page.open(BUILDER)
  await page
    .locator('.grid-stack-item')
    .first()
    .waitFor({ state: 'visible', timeout: 20000 })
    .catch(() => {})
  return page
}

// The record form waits on the module store, and rebuilds its record once the
// module lands — which blanks anything typed before that.
async function openRecordForm(page) {
  await page.open(RECORD_CREATE)
  await page.locator('form').first().waitFor({ state: 'visible', timeout: 25000 })
  await page.raw.waitForTimeout(2500)
  return page
}

drive('the builder warns after a block is added', async page => {
  await openBuilder(page)
  await page.click('button', { hasText: 'Add block' })
  await page.click('.p-dialog .grid > div.cursor-pointer', { hasText: 'Content' })
  // The new block is staged in its configurator; saving that puts it on the grid.
  await page.click('.p-dialog-footer button', { hasText: 'Save' })

  page.clearDialog()
  await page.back()
  const d = page.dialog()
  check('leaving prompts a confirm', !!d, d || 'no dialog')
})

drive('the builder warns after a block is edited', async page => {
  await openBuilder(page)
  await page.click('.block-toolbox button', { nth: 0 })
  await page.click('.p-dialog [role="tab"]', { hasText: 'General' })
  await page.fill('.p-dialog input[placeholder="Block title"]', 'Retitled by a check')
  await page.click('.p-dialog-footer button', { hasText: 'Save' })

  page.clearDialog()
  await page.back()
  const d = page.dialog()
  check('leaving prompts a confirm', !!d, d || 'no dialog')
})

drive('an untouched builder does not warn', async page => {
  await openBuilder(page)
  page.clearDialog()
  await page.back()
  const d = page.dialog()
  check('no confirm on a clean leave', !d, d || '')
})

drive('the record create form warns once a field is filled', async page => {
  await openRecordForm(page)
  // Typed, not filled: a fill() sets the value with no pointer or key event, so
  // it exercises a path no user can take.
  //
  // The module has 26 text fields and any one of them dirties the form; the
  // first is as good as another, so long as the choice is stated.
  await page.click('form input[type="text"]', { nth: 0 })
  await page.raw.keyboard.type('Draft that must not vanish')
  await page.raw.keyboard.press('Tab')
  await page.raw.waitForTimeout(1000)

  page.clearDialog()
  // No Back button here — a user leaves by a topbar link, which routes in-app.
  // (Browser back unloads the SPA and never reaches the route guard.)
  await page.click('#topbar-tools a', { hasText: 'Edit module' })
  await page.raw.waitForTimeout(1500)
  const d = page.dialog()
  check('leaving prompts a confirm', !!d, d || 'no dialog')
})

drive('an untouched record create form does not warn', async page => {
  await openRecordForm(page)
  page.clearDialog()
  await page.click('#topbar-tools a', { hasText: 'Edit module' })
  await page.raw.waitForTimeout(1500)
  const d = page.dialog()
  check('no confirm on a clean leave', !d, d || '')
})
